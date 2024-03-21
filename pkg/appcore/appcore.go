package appcore

import (
	"context"
	"flag"
	"fmt"
	"github.com/wolfelee/gocomm/pkg/jlog"
	"path/filepath"
	"sync"

	"github.com/wolfelee/gocomm/pkg/conf"
	"github.com/wolfelee/gocomm/pkg/server"
	"github.com/wolfelee/gocomm/pkg/signals"
	"github.com/wolfelee/gocomm/pkg/util/xcycle"
	"github.com/wolfelee/gocomm/pkg/util/xdefer"
	"github.com/wolfelee/gocomm/pkg/util/xgo"
	"github.com/wolfelee/gocomm/pkg/worker"
	job "github.com/wolfelee/gocomm/pkg/worker/xjob"
	"golang.org/x/sync/errgroup"
)

const (
	//StageAfterStop after app stop
	StageAfterStop uint32 = iota + 1
	//StageBeforeStop before app stop
	StageBeforeStop
)

type Application struct {
	cycle       *xcycle.Cycle
	smu         *sync.RWMutex
	initOnce    sync.Once
	startupOnce sync.Once
	stopOnce    sync.Once
	servers     []server.Server
	workers     []worker.Worker
	jobs        map[string]job.Runner
	hooks       map[uint32]*xdefer.DeferStack

	ConfPath string //配置所在路径
}

// New new a Application
func New(fns ...func() error) (*Application, error) {
	app := &Application{}
	if err := app.Startup(fns...); err != nil {
		return nil, err
	}
	return app, nil
}

func DefaultApp() *Application {
	app := &Application{}
	app.initialize()
	return app
}

// init hooks
func (app *Application) initHooks(hookKeys ...uint32) {
	app.hooks = make(map[uint32]*xdefer.DeferStack, len(hookKeys))
	for _, k := range hookKeys {
		app.hooks[k] = xdefer.NewStack()
	}
}

// run hooks
func (app *Application) runHooks(k uint32) {
	hooks, ok := app.hooks[k]
	if ok {
		hooks.Clean()
	}
}

func (app *Application) RegisterHooks(k uint32, fns ...func() error) error {
	hooks, ok := app.hooks[k]
	if ok {
		hooks.Push(fns...)
		return nil
	}
	return fmt.Errorf("hook stage not found")
}

// initialize application
func (app *Application) initialize() {
	app.initOnce.Do(func() {
		//assign
		app.cycle = xcycle.NewCycle()
		app.smu = &sync.RWMutex{}
		app.servers = make([]server.Server, 0)
		app.workers = make([]worker.Worker, 0)
		app.jobs = make(map[string]job.Runner)
		app.initHooks(StageBeforeStop, StageAfterStop)
	})
}

func (app *Application) startup() (err error) {
	app.startupOnce.Do(func() {
		err = xgo.SerialUntilError(
			app.parseFlags,
			app.loadConfig,
			app.initLogger,
		)()
	})
	return
}

func (app *Application) Startup(fns ...func() error) error {
	app.initialize()
	if err := app.startup(); err != nil {
		return err
	}
	return xgo.SerialUntilError(fns...)()
}

func (app *Application) Serve(s ...server.Server) error {
	app.smu.Lock()
	defer app.smu.Unlock()
	app.servers = append(app.servers, s...)
	return nil
}

func (app *Application) Schedule(w worker.Worker) error {
	app.workers = append(app.workers, w)
	return nil
}

func (app *Application) Run(servers ...server.Server) error {
	app.smu.Lock()
	app.servers = append(app.servers, servers...)
	app.smu.Unlock()

	app.waitSignals() //start signal listen task in goroutine
	defer app.clean()

	// todo jobs not graceful
	app.startJobs()

	// start servers and govern server
	app.cycle.Run(app.startServers)
	// start workers
	app.cycle.Run(app.startWorkers)

	//blocking and wait quit
	if err := <-app.cycle.Wait(); err != nil {
		jlog.Error("shutdown with error", jlog.Any("err", err))
		return err
	}
	jlog.Info("shutdown, bye!")
	return nil
}

func (app *Application) clean() {
}

// Stop application immediately after necessary cleanup
func (app *Application) Stop() (err error) {
	app.stopOnce.Do(func() {
		app.runHooks(StageBeforeStop)

		//stop servers
		app.smu.RLock()
		for _, s := range app.servers {
			func(s server.Server) {
				app.cycle.Run(s.Stop)
			}(s)
		}
		app.smu.RUnlock()

		//stop workers
		for _, w := range app.workers {
			func(w worker.Worker) {
				app.cycle.Run(w.Stop)
			}(w)
		}
		<-app.cycle.Done()
		app.runHooks(StageAfterStop)
		app.cycle.Close()
	})
	return
}

// GracefulStop application after necessary cleanup
func (app *Application) GracefulStop(ctx context.Context) (err error) {
	app.stopOnce.Do(func() {
		app.runHooks(StageBeforeStop)

		//stop servers
		app.smu.RLock()
		for _, s := range app.servers {
			func(s server.Server) {
				app.cycle.Run(func() error {
					return s.GracefulStop(ctx)
				})
			}(s)
		}
		app.smu.RUnlock()

		//stop workers
		for _, w := range app.workers {
			func(w worker.Worker) {
				app.cycle.Run(w.Stop)
			}(w)
		}
		<-app.cycle.Done()
		app.runHooks(StageAfterStop)
		app.cycle.Close()
	})
	return err
}

// waitSignals wait signal
func (app *Application) waitSignals() {
	signals.Shutdown(func(grace bool) { //when get shutdown signal
		//todo: support timeout
		if grace {
			app.GracefulStop(context.TODO())
		} else {
			app.Stop()
		}
	})
}

func (app *Application) startServers() error {
	var eg errgroup.Group
	// start multi servers
	for _, s := range app.servers {
		s := s
		eg.Go(func() (err error) {
			defer jlog.Info("exit server:" + s.Info().Name)
			err = s.Serve()
			return
		})
	}
	return eg.Wait()
}

func (app *Application) startWorkers() error {
	var eg errgroup.Group
	// start multi workers
	for _, w := range app.workers {
		w := w
		eg.Go(func() error {
			return w.Run()
		})
	}
	return eg.Wait()
}

// todo handle error
func (app *Application) startJobs() error {
	if len(app.jobs) == 0 {
		return nil
	}
	var jobs = make([]func(), 0)
	//warp jobs
	for name, runner := range app.jobs {
		jobs = append(jobs, func() {
			jlog.Info("job run begin:" + name)
			defer jlog.Info("job run end:" + name)
			// runner.Run panic 错误在更上层抛出
			runner.Run()
		})
	}
	xgo.Parallel(jobs...)()
	return nil
}

func (app *Application) parseFlags() error {
	flag.StringVar(&app.ConfPath, "c", "conf", "config path")
	flag.Parse()
	return nil
}

func (app *Application) initLogger() error {
	jlog.StdConfig().Build()
	return nil
}

func (app *Application) loadConfig() error {
	conf.Init(filepath.Join(app.ConfPath, "app.yml"))
	return nil
}
