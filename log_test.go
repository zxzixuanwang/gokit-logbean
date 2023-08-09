package logbean

import (
	"testing"
)

func TestStdOutput(t *testing.T) {
	l := InitLogBean(WithService("test"))
	l.Info("i am", "testing")
}

func TestJsonFile(t *testing.T) {
	l := InitLogBean(WithService("test_json_file"), WithFilePostion("./app.log"), WithType(JsonFile))
	l.Info("i am", "info")
	l.Debug("i am", "debug")
}

func TestFmtFile(t *testing.T) {
	l := InitLogBean(WithService("test_fmt_file"), WithFilePostion("./app.log"), WithType(File))
	l.Info("i am", "info")
	l.Debug("i am", "debug")
}

func TestCustomWriter(t *testing.T) {
	l := InitLogBean(WithService("test_writer"), WithWriter(openFile("./custom.log")), WithType(File))
	l.Info("i am", "info")
	l.Debug("i am", "debug")
}
