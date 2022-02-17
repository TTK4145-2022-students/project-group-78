package utils

import (
	"runtime"

	"github.com/elliotchance/pie/pie"
	"github.com/sirupsen/logrus"
)

func NewLogger(pkg string, fieldOrder ...string) *logrus.Logger {
	l := logrus.New()
	l.SetReportCaller(true)
	l.SetFormatter(&logrus.TextFormatter{
		ForceColors:     true,
		FullTimestamp:   true,
		TimestampFormat: "15:04:05",
		SortingFunc: func(ks []string) {
			// Place ordered fields first, and then sort the rest alphabetically
			fieldOrder := pie.Strings(fieldOrder)
			keys := pie.Strings(ks).Sort()
			i := 0
			for _, key := range fieldOrder {
				if keys.Contains(key) {
					ks[i] = key
					i++
				}
			}
			for _, key := range keys {
				if !fieldOrder.Contains(key) {
					ks[i] = key
					i++
				}
			}
		},
		CallerPrettyfier: func(_ *runtime.Frame) (function string, file string) {
			return "", pkg
		},
	})
	return l
}
