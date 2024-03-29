package apollo

import (
	"github.com/apolloconfig/agollo/v4/component/log"
	"github.com/lamber92/go-brick/blog"
	"github.com/lamber92/go-brick/blog/logger"
)

const (
	moduleKey  = "module"
	moduleName = "apollo"
)

type defaultLogger struct {
	logger logger.Logger
	debug  bool
}

func newDefaultLogger(debug bool) log.LoggerInterface {
	return &defaultLogger{
		logger: logger.Infra.WithOptions(logger.AddCallerSkip(2)),
		debug:  debug,
	}
}

func (d *defaultLogger) Debugf(format string, params ...interface{}) {
	if d.debug {
		d.logger.With(blog.String(moduleKey, moduleName)).Debugf(format, params...)
	}
}

func (d *defaultLogger) Infof(format string, params ...interface{}) {
	d.logger.With(blog.String(moduleKey, moduleName)).Infof(format, params...)
}

func (d *defaultLogger) Warnf(format string, params ...interface{}) {
	d.logger.With(blog.String(moduleKey, moduleName)).Warnf(format, params...)
}

func (d *defaultLogger) Errorf(format string, params ...interface{}) {
	d.logger.With(blog.String(moduleKey, moduleName)).Errorf(format, params...)
}

func (d *defaultLogger) Debug(msg string) {
	if d.debug {
		d.logger.With(blog.String(moduleKey, moduleName)).Debug(msg)
	}
}

func (d *defaultLogger) Info(msg string) {
	d.logger.With(blog.String(moduleKey, moduleName)).Info(msg)
}

func (d *defaultLogger) Warn(msg string) {
	d.logger.With(blog.String(moduleKey, moduleName)).Warn(msg)
}

func (d *defaultLogger) Error(msg string) {
	d.logger.With(blog.String(moduleKey, moduleName)).Error(msg)
}
