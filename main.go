package main

import (
	"github.com/TTK4145-2022-students/project-group-78/distributor"
	log "github.com/sirupsen/logrus"
)

const PORT = 41875

func main() {
	log.SetLevel(log.DebugLevel)
	distributor.New(1)
}
