package utils

import (
	"github.com/sirupsen/logrus"
)

var Log *logrus.Logger

func InitLogger() {
	Log = logrus.New()

	Log.SetFormatter(&logrus.TextFormatter{
		ForceColors:     true,
		FullTimestamp:   true,
		TimestampFormat: "2006-01-02 15:04:05",
		// PrettyPrint:     true,
	})

	Log.SetLevel(logrus.DebugLevel)
}
