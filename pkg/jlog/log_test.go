package jlog_test

import (
	"github.com/wolfelee/gocomm/pkg/jlog"
	"testing"
)

func Test_Info(t *testing.T) {
	jlog.Info("hello", jlog.Any("a", "b"))
}
