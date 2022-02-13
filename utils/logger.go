package utils

import "github.com/sirupsen/logrus"

func NewLogger() *logrus.Logger {
	l := logrus.New()
	l.SetFormatter(&logrus.TextFormatter{ForceColors: true, FullTimestamp: true, TimestampFormat: "15:04:05"})
	return l
}
