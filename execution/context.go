package execution

import (
	"context"
	"fmt"
)

type ctxKey string

const (
	executorIDCtxKey    ctxKey = "executorID"
	executorPathCtxKey  ctxKey = "executorPath"
	executorDepthCtxKey ctxKey = "executorDepth"
)

func WithExecutorID(ctx context.Context, executorID string) context.Context {
	currentPath := ExecutorPath(ctx)
	newPath := fmt.Sprintf("%s/%s", currentPath, executorID)
	currentDepth := ExecutorDepth(ctx)
	newDepth := currentDepth + 1
	ctx = context.WithValue(ctx, executorIDCtxKey, executorID)
	ctx = context.WithValue(ctx, executorPathCtxKey, newPath)
	ctx = context.WithValue(ctx, executorDepthCtxKey, newDepth)
	return ctx
}

func ExecutorID(ctx context.Context) string {
	v, ok := ctx.Value(executorIDCtxKey).(string)
	if !ok {
		return ""
	}
	return v
}

func ExecutorPath(ctx context.Context) string {
	v, ok := ctx.Value(executorPathCtxKey).(string)
	if !ok {
		return ""
	}
	return v
}

func ExecutorDepth(ctx context.Context) int {
	v, ok := ctx.Value(executorDepthCtxKey).(int)
	if !ok {
		return 0
	}
	return v
}
