// Copyright 2015 The Xorm Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package db

import (
	"github.com/wolfelee/gocomm/pkg/jlog"
	xormLog "xorm.io/xorm/log"
)

const (
	// !nashtsai! following level also match syslog.Priority value
	LOG_DEBUG xormLog.LogLevel = iota
	LOG_INFO
	LOG_WARNING
	LOG_ERR
	LOG_OFF
	LOG_UNKNOWN
)

type SimpleJdLogger struct {
	level   xormLog.LogLevel
	showSQL bool
}

var _ xormLog.Logger = &SimpleJdLogger{}

func NewSimpleJdLogger() *SimpleJdLogger {
	return &SimpleJdLogger{}
}

// Error implement ILogger
func (s *SimpleJdLogger) Error(v ...interface{}) {
	jlog.Error("", jlog.Any("sql", v))
}

// Errorf implement ILogger
func (s *SimpleJdLogger) Errorf(format string, v ...interface{}) {
	jlog.Error("", jlog.Any("sql", v))
}

// Debug implement ILogger
func (s *SimpleJdLogger) Debug(v ...interface{}) {
	jlog.Debug("", jlog.Any("sql", v))
}

// Debugf implement ILogger
func (s *SimpleJdLogger) Debugf(format string, v ...interface{}) {
	jlog.Debug("", jlog.Any("sql", v))
}

// Info implement ILogger
func (s *SimpleJdLogger) Info(v ...interface{}) {
	jlog.Info("", jlog.Any("sql", v))
}

// Infof implement ILogger
func (s *SimpleJdLogger) Infof(format string, v ...interface{}) {
	if len(v) >= 4 {
		jlog.Info("", jlog.Any("sql", v), jlog.Any("sqlcost", v[3]))
	} else {
		jlog.Info("", jlog.Any("sql", v))
	}
}

// Warn implement ILogger
func (s *SimpleJdLogger) Warn(v ...interface{}) {
	jlog.Warn("", jlog.Any("sql", v))
}

// Warnf implement ILogger
func (s *SimpleJdLogger) Warnf(format string, v ...interface{}) {
	jlog.Warn("", jlog.Any("sql", v))
}

// Level implement ILogger
func (s *SimpleJdLogger) Level() xormLog.LogLevel {
	return s.level
}

// SetLevel implement ILogger
func (s *SimpleJdLogger) SetLevel(l xormLog.LogLevel) {
	s.level = l
}

// ShowSQL implement ILogger
func (s *SimpleJdLogger) ShowSQL(show ...bool) {
	if len(show) == 0 {
		s.showSQL = true
		return
	}
	s.showSQL = show[0]
}

// IsShowSQL implement ILogger
func (s *SimpleJdLogger) IsShowSQL() bool {
	return s.showSQL
}
