package utils

import "github.com/sirupsen/logrus"

func NewLogger(pkg string) *logrus.Entry {
	l := logrus.New()
	l.SetFormatter(&logrus.TextFormatter{
		ForceColors:     true,
		FullTimestamp:   true,
		TimestampFormat: "15:04:05",
		//PadLevelText:    true,
	})
	return l.WithField("pkg", pkg)
}
