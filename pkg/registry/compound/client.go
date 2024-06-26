package compound

import (
	"context"

	registry2 "github.com/wolfelee/gocomm/pkg/registry"
	"github.com/wolfelee/gocomm/pkg/server"
	"golang.org/x/sync/errgroup"
)

type compoundRegistry struct {
	registries []registry2.Registry
}

// ListServices ...
func (c compoundRegistry) ListServices(ctx context.Context, name string, scheme string) ([]*server.ServiceInfo, error) {
	var eg errgroup.Group
	var services = make([]*server.ServiceInfo, 0)
	for _, registry := range c.registries {
		registry := registry
		eg.Go(func() error {
			infos, err := registry.ListServices(ctx, name, scheme)
			if err != nil {
				return err
			}
			services = append(services, infos...)
			return nil
		})
	}
	err := eg.Wait()
	return services, err
}

// WatchServices ...
func (c compoundRegistry) WatchServices(ctx context.Context, s string, s2 string) (chan registry2.Endpoints, error) {
	panic("compound registry doesn't support watch services")
}

// RegisterService ...
func (c compoundRegistry) RegisterService(ctx context.Context, bean *server.ServiceInfo) error {
	var eg errgroup.Group
	for _, registry := range c.registries {
		registry := registry
		eg.Go(func() error {
			return registry.RegisterService(ctx, bean)
		})
	}
	return eg.Wait()
}

// UnregisterService ...
func (c compoundRegistry) UnregisterService(ctx context.Context, bean *server.ServiceInfo) error {
	var eg errgroup.Group
	for _, registry := range c.registries {
		registry := registry
		eg.Go(func() error {
			return registry.UnregisterService(ctx, bean)
		})
	}
	return eg.Wait()
}

// Close ...
func (c compoundRegistry) Close() error {
	var eg errgroup.Group
	for _, registry := range c.registries {
		registry := registry
		eg.Go(func() error {
			return registry.Close()
		})
	}
	return eg.Wait()
}

// New ...
func New(registries ...registry2.Registry) registry2.Registry {
	return compoundRegistry{
		registries: registries,
	}
}
