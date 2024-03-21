package xgrpc

import (
	"context"
	"errors"
	"fmt"
	"github.com/wolfelee/gocomm/pkg/jlog"
	"github.com/wolfelee/gocomm/pkg/jtrace"
	"go.uber.org/zap"
	"net"
	"runtime"
	"strings"
	"time"

	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/peer"

	"google.golang.org/grpc"
)

type contextedServerStream struct {
	grpc.ServerStream
	ctx context.Context
}

// Context ...
func (css contextedServerStream) Context() context.Context {
	return css.ctx
}

func extractAppName(ctx context.Context) string {
	if md, ok := metadata.FromIncomingContext(ctx); ok {
		return strings.Join(md.Get("appName"), ",")
	}
	return ""
}

func extractHiddenLog(ctx context.Context) string {
	if md, ok := metadata.FromIncomingContext(ctx); ok {
		return strings.Join(md.Get("hiddenLog"), ",")
	}
	return ""
}

func extractTraceId(ctx context.Context) []string {
	if md, ok := metadata.FromIncomingContext(ctx); ok {
		return md.Get(jtrace.TraceIDKey)
	}
	return []string{}
}

func defaultStreamServerInterceptor() grpc.StreamServerInterceptor {
	return func(srv interface{}, stream grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) (err error) {
		var beg = time.Now()
		var fields = make([]jlog.Field, 0, 8)
		var event = "normal"
		defer func() {
			if rec := recover(); rec != nil {
				switch rec := rec.(type) {
				case error:
					err = rec
				default:
					err = fmt.Errorf("%v", rec)
				}
				stack := make([]byte, 4096)
				stack = stack[:runtime.Stack(stack, true)]
				fields = append(fields, jlog.FieldStack(stack))
				event = "recover"
			}

			fields = append(fields,
				jlog.FieldMethod(info.FullMethod),
				jlog.FieldCost(time.Since(beg)),
				jlog.FieldEvent(event),
			)

			for key, val := range getPeer(stream.Context()) {
				fields = append(fields, jlog.Any(key, val))
			}

			if err != nil {
				fields = append(fields, zap.String("err", err.Error()))
				jlog.Error("grpcStreamS", fields...)
				return
			}
			jlog.Info("grpcStreamS", fields...)
		}()

		return handler(srv, stream)
	}
}

func defaultUnaryServerInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {

		appName := extractAppName(ctx)
		if len(appName) == 0 {
			return nil, errors.New("appName is nil")
		}
		var beg = time.Now()
		var fields = make([]jlog.Field, 0, 8)
		var event = "normal"
		defer func() {
			if rec := recover(); rec != nil {
				switch rec := rec.(type) {
				case error:
					err = rec
				default:
					err = fmt.Errorf("%v", rec)
				}

				stack := make([]byte, 4096)
				stack = stack[:runtime.Stack(stack, true)]
				fields = append(fields, jlog.FieldStack(stack))
				event = "recover"
			}

			fields = append(fields,
				jlog.FieldMethod(info.FullMethod),
				jlog.FieldCost(time.Since(beg)),
				jlog.FieldEvent(event),
			)

			for key, val := range getPeer(ctx) {
				fields = append(fields, jlog.Any(key, val))
			}

			if err != nil {
				fields = append(fields, zap.String("err", err.Error()))
				//服务端不处理报警，由调用端进行错误的处理
				jlog.Warn("access", fields...)
				return
			}

			fields = append(fields,
				jlog.Any("req", req),
				jlog.Any("traceId", extractTraceId(ctx)),
				jlog.Any("appName", appName),
			)
			if appName != "health" && extractHiddenLog(ctx) == "" {
				//心跳接口不打印, 如果有屏蔽日志的请求也不打印
				jlog.Info("grpcS", fields...)
			}
		}()

		logger := jlog.With(jlog.Any(jtrace.TraceIDKey, extractTraceId(ctx)))
		newCtx := jlog.ToContext(ctx, logger)

		return handler(newCtx, req)
	}
}

func getClientIP(ctx context.Context) (string, error) {
	pr, ok := peer.FromContext(ctx)
	if !ok {
		return "", fmt.Errorf("[getClinetIP] invoke FromContext() failed")
	}
	if pr.Addr == net.Addr(nil) {
		return "", fmt.Errorf("[getClientIP] peer.Addr is nil")
	}
	addSlice := strings.Split(pr.Addr.String(), ":")
	return addSlice[0], nil
}

func getPeer(ctx context.Context) map[string]string {
	var peerMeta = make(map[string]string)
	if md, ok := metadata.FromIncomingContext(ctx); ok {
		if val, ok := md["aid"]; ok {
			peerMeta["aid"] = strings.Join(val, ";")
		}
		var clientIP string
		if val, ok := md["client-ip"]; ok {
			clientIP = strings.Join(val, ";")
		} else {
			ip, err := getClientIP(ctx)
			if err == nil {
				clientIP = ip
			}
		}
		peerMeta["clientIP"] = clientIP
		if val, ok := md["client-host"]; ok {
			peerMeta["host"] = strings.Join(val, ";")
		}
	}
	return peerMeta

}
