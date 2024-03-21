package db

import (
	"context"
	"errors"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/wolfelee/gocomm/pkg/jlog"
	"github.com/wolfelee/gocomm/pkg/jtrace"
	"google.golang.org/grpc/metadata"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"time"
	"xorm.io/xorm"
	xormLog "xorm.io/xorm/log"
)

var (
	dbCfg     *dbConfig
	DataGroup map[string]*xorm.EngineGroup
)

type dbGroupConfig struct {
	OpenConns       int `yaml:"openConns"`
	IdleConns       int `yaml:"idleConns"`
	ConnMaxLifetime int `yaml:"maxLifetime"`

	Master *dbDetail   `yaml:"master"`
	Slaves []*dbDetail `yaml:"slaves"`
}

type dbDetail struct {
	Name       string `yaml:"name"`
	User       string `yaml:"user"`
	Password   string `yaml:"password"`
	Host       string `yaml:"host"`
	Port       int    `yaml:"port"`
	Charset    string `yaml:"charset"`
	TzLocation string `yaml:"tz_location"`
	DatabaseTz string `yaml:"database_tz"`
}

func (dd *dbDetail) parseConnStr() string {
	//"xxx:xxx@tcp(192.168.1.111:3306)/gotest?loc=Local&charset=utf8mb4"
	addr := fmt.Sprintf("tcp(%s:%d)", dd.Host, dd.Port)
	connStr := fmt.Sprintf("%s:%s@%s/%s?charset=%s",
		dd.User,
		dd.Password,
		addr,
		dd.Name,
		dd.Charset,
	)
	return connStr
}

type dbConfig struct {
	Adapter string                    `yaml:"adapter"`
	ShowSQL bool                      `yaml:"showsql"`
	Db      map[string]*dbGroupConfig `yaml:"db"`
}

func initDataGroup() (map[string]*xorm.EngineGroup, error) {
	var groups = make(map[string]*xorm.EngineGroup)
	if dbCfg == nil {
		return nil, errors.New("db config setting error")
	}
	for g, e := range dbCfg.Db {
		dataSourceSlice := make([]string, 0)
		dataSourceSlice = append(dataSourceSlice, e.Master.parseConnStr())
		for _, sn := range dbCfg.Db[g].Slaves {
			dataSourceSlice = append(dataSourceSlice, sn.parseConnStr())
		}
		if len(dataSourceSlice) > 0 {
			group, err := xorm.NewEngineGroup(dbCfg.Adapter, dataSourceSlice)
			if err != nil {
				return nil, errors.New("创建数据库组链接错误：" + err.Error())
			}

			logger := NewSimpleJdLogger()
			group.SetLogger(xormLog.NewLoggerAdapter(logger))

			group.SetMaxOpenConns(dbCfg.Db[g].OpenConns)
			group.SetMaxIdleConns(dbCfg.Db[g].IdleConns)
			group.SetConnMaxLifetime(time.Duration(dbCfg.Db[g].ConnMaxLifetime) * time.Second)
			group.ShowSQL(dbCfg.ShowSQL)
			group.EnableSessionID(true)
			err = group.Ping()
			if err != nil {
				return nil, errors.New("DB ping error：" + err.Error())
			}

			if e.Master.DatabaseTz != "" {
				group.DatabaseTZ, _ = time.LoadLocation(e.Master.DatabaseTz)
			}
			if e.Master.TzLocation != "" {
				group.TZLocation, _ = time.LoadLocation(e.Master.TzLocation)
			}

			groups[g] = group

			jlog.Info(fmt.Sprintf("%s EngineGroup Opened", g))
		}
	}
	return groups, nil
}

// 兼用旧的使用方式
func DB(dbName string) *xorm.EngineGroup {
	return Use(dbName)
}

func UseWithTraceId(dbName string, traceId string) *xorm.Session {
	if traceId != "" {
		ctx1 := context.WithValue(context.Background(), "__xorm_session_id", traceId)
		return Use(dbName).NewSession().Context(ctx1)
	} else {
		return Use(dbName).NewSession()
	}
}

func extractTraceId(ctx context.Context) string {
	if md, ok := metadata.FromOutgoingContext(ctx); ok {
		return md.Get(jtrace.TraceIDKey)[0]
	}
	return ""
}

func UseWithCtx(ctx context.Context, dbName string) *xorm.Session {
	//从ctx中提取出ctx来
	ctx1 := context.WithValue(context.Background(), "__xorm_session_id", extractTraceId(ctx))
	return Use(dbName).NewSession().Context(ctx1)
}

func Use(dbName string) *xorm.EngineGroup {
	if DataGroup == nil {
		var err error
		DataGroup, err = initDataGroup()
		if err != nil {
			jlog.Error(err.Error())
			return nil
		}
	}
	if g, ok := DataGroup[dbName]; ok {
		return g
	} else {
		jlog.Error(dbName + " - database does not exist.")
		return nil
	}
}

func Init(dbCfgFile string) error {
	buf, err := ioutil.ReadFile(dbCfgFile)
	if err != nil {
		jlog.Error(dbCfgFile + " file read error")
		return err
	}
	err = yaml.Unmarshal(buf, &dbCfg)
	if err != nil {
		jlog.Error(dbCfgFile + "file unmarshal error")
		return err
	}
	DataGroup, err = initDataGroup()
	if err != nil {
		jlog.Error(err.Error())
		return err
	}
	return nil
}

func Close() {
	for n, db := range DataGroup {
		db.Close()
		jlog.Info(fmt.Sprintf("%s EngineGroup Closed", n))
	}
}
