package utils

import (
	"context"
	"time"

	"github.com/surfe/logger/v2"
)

func LogFunctionDuration(functionName string, startTime time.Time) {
	logger.Log(context.Background()).Infof("func %v started at %v", functionName, startTime)
	logger.Log(context.Background()).Infof("func %v completed in %v seconds", functionName, time.Since(startTime).Seconds())
}

func LogFunctionDurationWithContext(ctx context.Context, functionName string, startTime time.Time) {
	logger.Log(ctx).Infof("Func %v started at %v", functionName, startTime)
	logger.Log(ctx).Infof("Func %v completed in %v seconds", functionName, time.Since(startTime).Seconds())
}
