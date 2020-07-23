package main

import (
	"github.com/grindlabs/sorteiagram/cmd"
	log "github.com/sirupsen/logrus"
)

func init() {
	log.SetLevel(log.DebugLevel)
}

func main() {
	cmd.Execute()
}
