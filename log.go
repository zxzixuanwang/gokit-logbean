package logbean

import (
	"io"
	"os"

	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
)

const (
	Std Logtype = iota
	File
	JsonFile
)

type (
	Logtype int

	logInfo struct {
		filePosition string
		level        string
		logtype      Logtype
		l            log.Logger
		service      *string
		caller       int
		timeFormat   log.Valuer
		writer       io.Writer
		otherPrefix  map[string]interface{}
	}

	Options        func(*logInfo)
	LogBeanService interface {
		GetType() Logtype
		GetLevel() string
		GetFilePostion() string
		Warn(log ...interface{})
		Info(log ...interface{})
		Debug(log ...interface{})
		Error(log ...interface{})
	}
)

var _ LogBeanService = (*logInfo)(nil)

func WithOtherPrefix(input map[string]interface{}) Options {
	return func(li *logInfo) {
		li.otherPrefix = input
	}
}

func WithWriter(writer io.Writer) Options {
	return func(li *logInfo) {
		li.writer = writer
	}
}

func WithService(service string) Options {
	return func(li *logInfo) {
		li.service = &service
	}
}

func WithCall(caller int) Options {
	return func(li *logInfo) {
		li.caller = caller
	}
}

func WithFilePostion(position string) Options {
	return func(li *logInfo) {
		li.filePosition = position
	}
}

func WithOutput(output Logtype) Options {
	return func(li *logInfo) {
		li.logtype = output
	}
}

func WithLevel(lev string) Options {
	return func(li *logInfo) {
		li.level = lev
	}
}

func WithTime(format log.Valuer) Options {
	return func(li *logInfo) {
		li.timeFormat = format
	}
}

func WithType(logtype Logtype) Options {
	return func(li *logInfo) {
		li.logtype = logtype
	}
}

func (li *logInfo) GetType() Logtype {
	return li.logtype
}

func (li *logInfo) GetLevel() string {
	return li.level
}

func (li *logInfo) GetFilePostion() string {
	return li.filePosition
}

func (li *logInfo) Info(log ...interface{}) {
	_ = level.Info(li.l).Log(log...)
}

func (li *logInfo) Warn(log ...interface{}) {
	_ = level.Warn(li.l).Log(log...)
}

func (li *logInfo) Debug(log ...interface{}) {
	_ = level.Debug(li.l).Log(log...)
}

func (li *logInfo) Error(log ...interface{}) {
	_ = level.Error(li.l).Log(log...)
}

func openFile(position string) *os.File {
	f, err := os.OpenFile(position, os.O_CREATE|os.O_RDWR|os.O_APPEND, 0644)
	if err != nil {
		panic(err)
	}
	return f
}

func InitLogBean(opt ...Options) LogBeanService {
	li := defaultInfo()
	for _, op := range opt {
		op(li)
	}
	if li.logtype > 0 && li.writer == nil {
		li.writer = openFile(li.filePosition)
	}

	var logger log.Logger
	switch li.logtype {
	case JsonFile:
		logger = log.NewJSONLogger(li.writer)
	case File:
		logger = log.NewLogfmtLogger(li.writer)
	default:
		logger = log.NewJSONLogger(os.Stdout)
	}

	logger = log.With(logger, "ts", li.timeFormat, "caller", log.Caller(li.caller))

	if li.service != nil {
		logger = log.With(logger, "service", li.service)
	}
	if len(li.otherPrefix) > 0 {
		for k, v := range li.otherPrefix {
			logger = log.With(logger, k, v)
		}
	}

	if li.level == "all" {
		logger = level.NewFilter(logger, level.AllowAll())
	} else {
		logger = level.NewFilter(logger, level.Allow(logLevelFilter(li.level)))
	}

	li.l = logger

	return li
}

func defaultInfo() *logInfo {
	return &logInfo{
		level:        "info",
		logtype:      Std,
		caller:       6,
		timeFormat:   log.DefaultTimestamp,
		filePosition: "./app.log",
	}
}

func logLevelFilter(lev string) level.Value {
	switch lev {
	case "debug":
		return level.DebugValue()
	case "error":
		return level.ErrorValue()
	case "warn":
		return level.WarnValue()
	default:
		return level.InfoValue()
	}
}
