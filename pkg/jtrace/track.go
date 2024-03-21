package jtrace

import (
	"context"
	"google.golang.org/grpc/metadata"
	"strings"
)

const TraceIDKey = "traceId"

// 生成新的traceId到ctx中
func ToContext(ctx context.Context, traceId string) context.Context {
	md, ok := metadata.FromOutgoingContext(ctx)
	mdtrace := metadata.Pairs(TraceIDKey, traceId)
	if ok {
		md = metadata.Join(md, mdtrace)
	} else {
		md = mdtrace
	}
	return metadata.NewOutgoingContext(ctx, md)
}

//增加一个判断有没有的，有traceId则继续，没有才重新生成
// Extract
func TraceId(ctx context.Context) string {
	if ctx == nil {
		return ""
	}
	if md, ok := metadata.FromIncomingContext(ctx); ok {
		if t, ok := md[TraceIDKey]; ok {
			return t[0]
		}
	}
	return ""
}

func ExtractAppName(ctx context.Context) string {
	if md, ok := metadata.FromIncomingContext(ctx); ok {
		return strings.Join(md.Get("appName"), ",")
	}
	return ""
}

func ExtractTraceId(ctx context.Context) []string {
	if md, ok := metadata.FromIncomingContext(ctx); ok {
		return md.Get(TraceIDKey)
	}
	return []string{}
}
