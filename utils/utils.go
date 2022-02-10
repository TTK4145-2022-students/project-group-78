package utils

import log "github.com/sirupsen/logrus"

func PanicIf(err error) {
	if err != nil {
		log.Panic()
	}
}
