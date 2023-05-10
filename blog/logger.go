package blog

var (
	Access = newAccessLogger()
	Biz    = newBizLogger()
	Infra  = newInfraLogger()
)

// ReplaceLogger replace all built-in logging engines
func ReplaceLogger(lgr Logger) {
	replaceAccessLogger(lgr)
	replaceBizLogger(lgr)
	replaceInfraLogger(lgr)
}

// Close disable all built-in logging engines
func Close() (out []error) {
	out = make([]error, 0, 3)
	handle := func(err error) []error {
		out = append(out, err)
		return out
	}

	handle(Access.Close())
	handle(Biz.Close())
	handle(Infra.Close())

	return out
}

func replaceAccessLogger(lgr Logger) {
	Access = lgr
}

func newAccessLogger() Logger {
	return &accessLogger{defaultLogger: newDefaultLogger(TypeAccess)}
}

type accessLogger struct {
	*defaultLogger
}

func replaceBizLogger(lgr Logger) {
	Biz = lgr
}

func newBizLogger() Logger {
	return &bizLogger{defaultLogger: newDefaultLogger(TypeBiz)}
}

type bizLogger struct {
	*defaultLogger
}

func replaceInfraLogger(lgr Logger) {
	Infra = lgr
}

func newInfraLogger() Logger {
	return &infraLogger{defaultLogger: newDefaultLogger(TypeInfra)}
}

type infraLogger struct {
	*defaultLogger
}
