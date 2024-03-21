package grpc

import (
	"context"
	"github.com/wolfelee/gocomm/pkg"
	"github.com/wolfelee/gocomm/pkg/jtrace"
	"github.com/wolfelee/gocomm/pkg/util/juuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

func extractTraceId(ctx context.Context) []string {
	if md, ok := metadata.FromOutgoingContext(ctx); ok {
		return md.Get(jtrace.TraceIDKey)
	}
	return []string{}
}

func loggerUnaryClientInterceptor() grpc.UnaryClientInterceptor {
	return func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		md, ok := metadata.FromOutgoingContext(ctx)

		clientAidMD := metadata.Pairs("appName", pkg.AppName())
		if ok {
			md = metadata.Join(md, clientAidMD)
		} else {
			md = clientAidMD
		}

		//检测有没有trackId，如果没有自动生成一个
		if len(extractTraceId(ctx)) == 0 {
			mdtrace := metadata.Pairs(jtrace.TraceIDKey, juuid.ShortUUID())
			md = metadata.Join(md, mdtrace)
		}

		ctx = metadata.NewOutgoingContext(ctx, md)

		return invoker(ctx, method, req, reply, cc, opts...)
	}
}
