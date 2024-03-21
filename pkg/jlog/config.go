package jlog

import (
	"fmt"
	"time"

	"github.com/wolfelee/gocomm/pkg/conf"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var CallerSkip = 1

func SetCallerSkip(caller int) {
	CallerSkip = caller
}

// Config ...
type Config struct {
	// Dir 日志输出目录
	Dir string
	// Name 日志文件名称
	Name string
	// Level 日志初始等级
	Level string
	// 日志初始化字段
	Fields []zap.Field
	// 是否添加调用者信息
	AddCaller bool
	// 日志前缀
	Prefix string
	// 日志输出文件最大长度，超过改值则截断
	MaxSize   int
	MaxAge    int
	MaxBackup int
	// 日志磁盘刷盘间隔
	Interval      time.Duration
	CallerSkip    int
	Async         bool
	Queue         bool
	QueueSleep    time.Duration
	Core          zapcore.Core
	Debug         bool
	EncoderConfig *zapcore.EncoderConfig
	configKey     string
}

// Filename ...
func (config *Config) Filename() string {
	return fmt.Sprintf("%s/%s", config.Dir, config.Name)
}

func StdConfig() *Config {
	var logCfg = DefaultConfig()
	logLevel := conf.GetString("logLevel")
	if len(logLevel) > 0 {
		logCfg.Level = logLevel
	}
	return logCfg
}

// DefaultConfig ...
func DefaultConfig() *Config {
	return &Config{
		Name:          "default.log",
		Dir:           ".",
		Level:         "info",
		MaxSize:       500, // 500M
		MaxAge:        1,   // 1 day
		MaxBackup:     10,  // 10 backup
		Interval:      24 * time.Hour,
		CallerSkip:    CallerSkip,
		AddCaller:     true,
		Async:         true,
		Queue:         false,
		QueueSleep:    100 * time.Millisecond,
		EncoderConfig: DefaultZapConfig(),
	}
}

// Build ...
func (config Config) Build() *Logger {
	if config.EncoderConfig == nil {
		config.EncoderConfig = DefaultZapConfig()
	}
	//if config.Debug {
	//	config.EncoderConfig.EncodeLevel = DebugEncodeLevel
	//}
	JLogger = newLogger(&config)
	//if config.configKey != "" {
	//	logger.AutoLevel(config.configKey + ".level")
	//}
	return JLogger
}
