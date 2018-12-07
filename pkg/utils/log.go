package utils

import (
	"fmt"
	"os"
	"time"

	"github.com/sirupsen/logrus"
)

type log struct {
	logrus.Logger
	OSFile *os.File
	err    error
}

// Log handle app logs
var Log = new(log)

func init() {
	formatter := &logrus.TextFormatter{}
	formatter.TimestampFormat = "2006-01-02 15:04:05"
	Log.Formatter = formatter
	Log.Hooks = make(logrus.LevelHooks)
	Log.Level = logrus.InfoLevel

	path := "../var/logs/"
	_, err := os.Stat(path)
	if err != nil {
		// create dir
		err = os.MkdirAll(path, os.ModePerm)
		if err != nil {
			fmt.Printf("mkdir failed![%v][path:%s]\n", err, path)
		}
	}
	logOut := Config.GetString("log.out")
	filename := time.Now().Format("2006-01-02")
	Log.OSFile, Log.err = os.OpenFile(fmt.Sprintf(path + "%s.log", filename), os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if Log.err == nil {
		if logOut == "console" {
			Log.Out = os.Stdout
		} else {
			Log.Out = Log.OSFile
		}
	} else {
		fmt.Printf("Failed to log to file - %s.log, using default stderr\n", filename)
	}
}
