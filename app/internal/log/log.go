package log

import (
	"fmt"
	"os"
	"strings"

	"github.com/fatih/color"
	"github.com/sirupsen/logrus"
)

var (
	LEVEL_COLOURS = map[logrus.Level]*color.Color{
		logrus.DebugLevel: color.New(color.FgBlack, color.Bold),
		logrus.InfoLevel:  color.New(color.FgBlue, color.Bold),
		logrus.WarnLevel:  color.New(color.FgYellow, color.Bold),
		logrus.ErrorLevel: color.New(color.FgRed),
		logrus.FatalLevel: color.New(color.FgRed, color.BgWhite),
		logrus.PanicLevel: color.New(color.FgRed, color.BgWhite),
	}
	log *logrus.Logger
)

type ColourFormatter struct{}

func (f *ColourFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	time := entry.Time.Format("2006-01-02 15:04:05")
	level := strings.ToUpper(entry.Level.String())
	levelColor := LEVEL_COLOURS[entry.Level]

	formattedMessage := fmt.Sprintf("%s [ %s ]%s %s%s\n",
		color.New(color.FgCyan).Sprint(time),
		levelColor.Sprint(level),
		color.New(color.Reset).Sprint(),
		color.New(color.FgWhite).Sprint(entry.Message),
		color.New(color.Reset).Sprint(),
	)

	return []byte(formattedMessage), nil
}

func init() {
	log = logrus.New()

	level := logrus.InfoLevel
	if os.Getenv("LOG_LEVEL") == "DEBUG" {
		level = logrus.DebugLevel
	}

	log.SetFormatter(&ColourFormatter{})
	log.SetLevel(level)
}

func Info(msg string) {
	log.Info(msg)
}

func Infof(msg string, args ...interface{}) {
	log.Infof(msg, args...)
}

func Warn(msg string) {
	log.Warn(msg)
}

func Warnf(msg string, args ...interface{}) {
	log.Warnf(msg, args...)
}

func Error(msg string) {
	log.Error(msg)
}

func Errorf(msg string, args ...interface{}) {
	log.Errorf(msg, args...)
}

func Debug(msg string) {
	log.Debug(msg)
}

func Debugf(msg string, args ...interface{}) {
	log.Debugf(msg, args...)
}

func Fatal(msg string) {
	log.Fatal(msg)
}

func Fatalf(msg string, args ...interface{}) {
	log.Fatalf(msg, args...)
}

func Panic(msg string) {
	log.Panic(msg)
}

func Panicf(msg string, args ...interface{}) {
	log.Panicf(msg, args...)
}
