package main

import (
	"io"
	"os"
	"strconv"

	"github.com/grindlabs/sorteiagram/cmd"
	"github.com/joho/godotenv"
	log "github.com/sirupsen/logrus"
	"gopkg.in/natefinch/lumberjack.v2"
)

func init() {
	if err := godotenv.Load(); err != nil {
		panic("Unable to load the env file")
	}

	debug, err := strconv.Atoi(os.Getenv("DEBUG"))

	if err != nil {
		panic("Unable to convert the DEBUG env var")
	}

	log.SetLevel(log.InfoLevel)

	if debug == 1 {
		log.SetLevel(log.DebugLevel)
	}

	log.SetFormatter(&log.TextFormatter{
		FullTimestamp: true,
		DisableColors: true,
	})

	multiWriter := io.MultiWriter(os.Stdout, &lumberjack.Logger{
		Filename:   "logs/sorteiagram.log",
		MaxSize:    50,
		MaxBackups: 3,
		MaxAge:     28,
		Compress:   true,
	})

	log.SetOutput(multiWriter)
}

func main() {
	cmd.Execute()
}
