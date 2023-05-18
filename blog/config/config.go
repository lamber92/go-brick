package config

import (
	"io"
	"log"
	"os"
	"path"
	"strings"
	"time"

	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type ZapConfig struct {
	Level      string   `json:"Level" yaml:"Level"`
	Debug      bool     `json:"Debug" yaml:"Debug"`
	Stacktrace string   `json:"Stacktrace" yaml:"Stacktrace"`
	Encoding   string   `json:"Encoding" yaml:"Encoding"`
	Output     []string `json:"Output" yaml:"Output"`
	// file attr
	FileAttr *FileAttr `json:"FileAttr" yaml:"FileAttr"`
}

func NewDefault() *ZapConfig {
	return &ZapConfig{
		Level:      "debug",
		Debug:      true,
		Stacktrace: "warn",
		Encoding:   "json",
		Output:     []string{"stdout"},
	}
}

// TODO: 加入到可读取配置文件获取日志配置

func LoadConfig() *ZapConfig {
	return &ZapConfig{}
}

func (c *ZapConfig) GetWriterSyncer() zapcore.WriteSyncer {
	var ws []zapcore.WriteSyncer
	for _, v := range c.Output {
		switch v {
		case "file":
			if c.FileAttr == nil {
				log.Panicln("defaultConfig.FileAttr is nil")
			}
			ws = append(ws, zapcore.AddSync(c.FileAttr.getWriter()))
		case "stdout":
			fallthrough
		default:
			ws = append(ws, zapcore.AddSync(os.Stdout))
		}
	}
	return zapcore.NewMultiWriteSyncer(ws...)
}

func (c *ZapConfig) GetLogLevel(lvl string) zapcore.LevelEnabler {
	atom := zap.DebugLevel
	switch lvl {
	case "debug":
		atom = zap.DebugLevel
	case "info":
		atom = zap.InfoLevel
	case "warn":
		atom = zap.WarnLevel
	case "error":
		atom = zap.ErrorLevel
	case "panic":
		atom = zap.PanicLevel
	default:
		atom = zap.DebugLevel
	}
	return zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl >= atom
	})
}

func (c *ZapConfig) GetEncoder() zapcore.Encoder {
	switch c.Encoding {
	case "console":
		return zapcore.NewConsoleEncoder(*c.buildEncodingConfig())
	case "json":
		fallthrough
	default:
		return zapcore.NewJSONEncoder(*c.buildEncodingConfig())
	}
}

func (c *ZapConfig) buildEncodingConfig() *zapcore.EncoderConfig {
	return &zapcore.EncoderConfig{
		TimeKey:     "time",
		LevelKey:    "level",
		NameKey:     "type",
		FunctionKey: "func",
		MessageKey:  "msg",
		// Do not use the built-in stack trace, because in the access log scenario,
		// the location where the log is output is different from the location where the error is generated,
		// which will cause the output to nest too many intermediate layers.
		// Use the berror.Stack() instead.
		// StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.CapitalLevelEncoder,
		EncodeTime:     zapcore.RFC3339TimeEncoder,
		EncodeDuration: zapcore.MillisDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}
}

type FileAttr struct {
	Directory string `json:"Directory" yaml:"Directory"`
	Name      string `json:"Name" yaml:"Name"`
	MaxAge    string `json:"MaxAge" yaml:"MaxAge"`
}

func (attr *FileAttr) getWriter() io.Writer {
	if err := os.MkdirAll(attr.Directory, os.ModeDir|0755); err != nil {
		log.Panicln(err)
	}
	maxAge, err := attr.getFileMaxAge()
	if err != nil {
		log.Panicln(err)
	}
	hook, err := rotatelogs.New(
		path.Join(attr.Directory, strings.Replace(attr.Name, ".log", "", -1)+"-%Y%m%d.log"),
		rotatelogs.WithMaxAge(maxAge),
		rotatelogs.WithRotationTime(time.Hour*24),
	)
	if err != nil {
		log.Panicln(err)
	}
	return hook
}

func (attr *FileAttr) getFileMaxAge() (time.Duration, error) {
	if len(attr.MaxAge) == 0 {
		return time.Hour * 24 * 7, nil
	} else {
		return time.ParseDuration(attr.MaxAge)
	}
}
