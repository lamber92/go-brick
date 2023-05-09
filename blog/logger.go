package blog

var (
	Access = newAccessLogger()
	Biz    = newBizLogger()
	Infra  = newInfraLogger()
)

func ReplaceLogger(lgr Logger) {
	replaceAccessLogger(lgr)
	replaceBizLogger(lgr)
	replaceInfraLogger(lgr)
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
