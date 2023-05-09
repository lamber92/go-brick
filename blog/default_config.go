package blog

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

type defaultConfig struct {
	Level      string   `json:"Level" yaml:"Level"`
	Debug      bool     `json:"Debug" yaml:"Debug"`
	Stacktrace string   `json:"Stacktrace" yaml:"Stacktrace"`
	Encoding   string   `json:"Encoding" yaml:"Encoding"`
	Output     []string `json:"Output" yaml:"Output"`
	// file attr
	FileAttr *configFileAttr `json:"FileAttr" yaml:"FileAttr"`
}

func (c *defaultConfig) getWriterSyncer() zapcore.WriteSyncer {
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

func (c *defaultConfig) getZapLogLevel(lvl string) zapcore.LevelEnabler {
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

func (c *defaultConfig) getEncoder() zapcore.Encoder {
	switch c.Encoding {
	case "console":
		return zapcore.NewConsoleEncoder(*c.buildZapEncodingConfig())
	case "json":
		fallthrough
	default:
		return zapcore.NewJSONEncoder(*c.buildZapEncodingConfig())
	}
}

func (c *defaultConfig) buildZapEncodingConfig() *zapcore.EncoderConfig {
	return &zapcore.EncoderConfig{
		TimeKey:        "time",
		LevelKey:       "level",
		NameKey:        "default",
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder,
		EncodeTime:     zapcore.RFC3339TimeEncoder,
		EncodeDuration: zapcore.MillisDurationEncoder,
	}
}

type configFileAttr struct {
	Directory string `json:"Directory" yaml:"Directory"`
	Name      string `json:"Name" yaml:"Name"`
	MaxAge    string `json:"MaxAge" yaml:"MaxAge"`
}

func (attr *configFileAttr) getWriter() io.Writer {
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

func (attr *configFileAttr) getFileMaxAge() (time.Duration, error) {
	if len(attr.MaxAge) == 0 {
		return time.Hour * 24 * 7, nil
	} else {
		return time.ParseDuration(attr.MaxAge)
	}
}
