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

	LogInfo struct {
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

	Options func(*LogInfo)
)

func WithOtherPrefix(input map[string]interface{}) Options {
	return func(li *LogInfo) {
		li.otherPrefix = input
	}
}

func WithWriter(writer io.Writer) Options {
	return func(li *LogInfo) {
		li.writer = writer
	}
}

func WithService(service string) Options {
	return func(li *LogInfo) {
		li.service = &service
	}
}

func WithCall(caller int) Options {
	return func(li *LogInfo) {
		li.caller = caller
	}
}

func WithFilePostion(position string) Options {
	return func(li *LogInfo) {
		li.filePosition = position
	}
}

func WithOutput(output Logtype) Options {
	return func(li *LogInfo) {
		li.logtype = output
	}
}

func WithLevel(lev string) Options {
	return func(li *LogInfo) {
		li.level = lev
	}
}

func WithTime(format log.Valuer) Options {
	return func(li *LogInfo) {
		li.timeFormat = format
	}
}

func WithType(logtype Logtype) Options {
	return func(li *LogInfo) {
		li.logtype = logtype
	}
}

func (li *LogInfo) GetType() Logtype {
	return li.logtype
}

func (li *LogInfo) GetLevel() string {
	return li.level
}

func (li *LogInfo) GetFilePostion() string {
	return li.filePosition
}

func (li *LogInfo) Info(log ...interface{}) {
	level.Info(li.l).Log(log...)
}

func (li *LogInfo) Warn(log ...interface{}) {
	level.Warn(li.l).Log(log...)
}

func (li *LogInfo) Debug(log ...interface{}) {
	level.Debug(li.l).Log(log...)
}

func (li *LogInfo) Error(log ...interface{}) {
	level.Error(li.l).Log(log)
}

func openFile(position string) *os.File {
	f, err := os.OpenFile(position, os.O_CREATE|os.O_RDWR|os.O_APPEND, 0644)
	if err != nil {
		panic(err)
	}
	return f
}

func InitLogBean(opt ...Options) *LogInfo {
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

func defaultInfo() *LogInfo {
	return &LogInfo{
		level:        "info",
		logtype:      Std,
		caller:       5,
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
