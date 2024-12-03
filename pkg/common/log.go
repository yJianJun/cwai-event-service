package common

import (
	"github.com/sirupsen/logrus"
	"log"
	"os"
)

func init() {
	// 设置日志格式为json格式
	logrus.SetFormatter(&logrus.JSONFormatter{
		TimestampFormat: "2006-01-02 15:04:05",
	})
	logrus.SetReportCaller(true)
}
func Debug(fields logrus.Fields, args ...interface{}) {
	setOutput(logrus.DebugLevel)
	logrus.WithFields(fields).Debug(args)
}
func Info(fields logrus.Fields, args ...interface{}) {
	setOutput(logrus.InfoLevel)
	logrus.WithFields(fields).Info(args)
}
func Warn(fields logrus.Fields, args ...interface{}) {
	setOutput(logrus.WarnLevel)
	logrus.WithFields(fields).Warn(args)
}
func Fatal(fields logrus.Fields, args ...interface{}) {
	setOutput(logrus.FatalLevel)
	logrus.WithFields(fields).Fatal(args)
}
func Error(fields logrus.Fields, args ...interface{}) {
	setOutput(logrus.ErrorLevel)
	logrus.WithFields(fields).Error(args)
}
func Panic(fields logrus.Fields, args ...interface{}) {
	setOutput(logrus.PanicLevel)
	logrus.WithFields(fields).Panic(args)
}
func Trace(fields logrus.Fields, args ...interface{}) {
	setOutput(logrus.TraceLevel)
	logrus.WithFields(fields).Trace(args)
}
func setOutput(level logrus.Level) {
	log.SetOutput(os.Stdout)
	logrus.SetLevel(level)
	return
}
