package blog

import (
	"go-brick/blog/logger"
)

// AddCallerSkip increases the number of callers skipped by caller annotation
// (as enabled by the AddCaller option).
func AddCallerSkip(skip int) logger.Option {
	return logger.AddCallerSkip(skip)
}
