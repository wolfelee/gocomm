package grpc

import (
	"context"
	. "github.com/smartystreets/goconvey/convey"
	"github.com/wolfelee/gocomm/pkg/util/xtest/proto/testproto"
	"github.com/wolfelee/gocomm/pkg/util/xtest/server/yell"
	"testing"
	"time"
)

// TestBase test direct dial with New()
func TestDirectGrpc(t *testing.T) {
	Convey("test direct grpc", t, func() {
		ctx := context.Background()
		ctx, cancel := context.WithTimeout(ctx, time.Second)
		defer cancel()
		res, err := directClient.SayHello(ctx, &testproto.HelloRequest{
			Name: "hello",
		})
		So(err, ShouldBeNil)
		So(res.Message, ShouldEqual, yell.RespFantasy.Message)
	})
}

func TestConfigBlockTrue(t *testing.T) {
	Convey("test no address block, and panic", t, func() {
		flag := false
		defer func() {
			if r := recover(); r != nil {
				flag = true
			}
			So(flag, ShouldEqual, true)
		}()
		cfg := DefaultConfig()
		cfg.OnDialError = "panic"
		newGRPCClient(cfg)
	})
}

func TestConfigBlockFalse(t *testing.T) {
	Convey("test no address and no block", t, func() {
		cfg := DefaultConfig()
		cfg.OnDialError = "panic"
		cfg.Block = false
		conn, _ := newGRPCClient(cfg)
		So(conn.GetState().String(), ShouldEqual, "IDLE")
	})
}
