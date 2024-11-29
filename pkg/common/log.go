package common

import (
	"ctyun-code.srdcloud.cn/aiplat/cwai-watcher/pkg/config"
	"fmt"
	"github.com/sirupsen/logrus"
	"io"
	"log"
	"os"
	"path"
)

func init() {
	// 设置日志格式为json格式
	logrus.SetFormatter(&logrus.JSONFormatter{
		TimestampFormat: "2006-01-02 15:04:05",
	})
	logrus.SetReportCaller(true)
}
func Debug(fields logrus.Fields, args ...interface{}) {
	setOutPutFile(logrus.DebugLevel)
	logrus.WithFields(fields).Debug(args)
}
func Info(fields logrus.Fields, args ...interface{}) {
	setOutPutFile(logrus.InfoLevel)
	logrus.WithFields(fields).Info(args)
}
func Warn(fields logrus.Fields, args ...interface{}) {
	setOutPutFile(logrus.WarnLevel)
	logrus.WithFields(fields).Warn(args)
}
func Fatal(fields logrus.Fields, args ...interface{}) {
	setOutPutFile(logrus.FatalLevel)
	logrus.WithFields(fields).Fatal(args)
}
func Error(fields logrus.Fields, args ...interface{}) {
	setOutPutFile(logrus.ErrorLevel)
	logrus.WithFields(fields).Error(args)
}
func Panic(fields logrus.Fields, args ...interface{}) {
	setOutPutFile(logrus.PanicLevel)
	logrus.WithFields(fields).Panic(args)
}
func Trace(fields logrus.Fields, args ...interface{}) {
	setOutPutFile(logrus.TraceLevel)
	logrus.WithFields(fields).Trace(args)
}
func setOutPutFile(level logrus.Level) {
	if _, err := os.Stat(config.Config.Log.Dir); os.IsNotExist(err) {
		err = os.MkdirAll(config.Config.Log.Dir, 0777)
		if err != nil {
			panic(fmt.Errorf("create log dir '%s' error: %s", config.Config.Log.Dir, err))
		}
	}
	name := ""
	switch level {
	case logrus.DebugLevel:
		name = "debug"
	case logrus.InfoLevel:
		name = "info"
	case logrus.WarnLevel:
		name = "warn"
	case logrus.FatalLevel:
		name = "fatal"
	case logrus.ErrorLevel:
		name = "error"
	case logrus.PanicLevel:
		name = "panic"
	case logrus.TraceLevel:
		name = "trace"
	default:
		panic(fmt.Errorf("invaild log level error %d", logrus.ErrorLevel))
	}
	fileName := path.Join(config.Config.Log.Dir, name+".log")
	var err error
	os.Stderr, err = os.OpenFile(fileName, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		fmt.Println("open log file err", err)
	}
	log.SetOutput(io.MultiWriter(os.Stderr, os.Stdout)) // io.MultiWriter 返回一个 io.Writer 对象
	logrus.SetLevel(level)
	return
}
