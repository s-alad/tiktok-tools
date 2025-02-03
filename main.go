package main

import (
	"os"

	_ "github.com/joho/godotenv/autoload"
	"github.com/rifflock/lfshook"
	"github.com/s-alad/tiktok-of-alexandria/cmd/alexandria"
	log "github.com/sirupsen/logrus"
)

func init() {
	/* log.SetReportCaller(true) */

	logfile := "run.log"

	log.SetOutput(os.Stdout)
	log.SetFormatter(&log.TextFormatter{
		ForceColors:   true,
		FullTimestamp: true,
	})
	_, err := os.OpenFile(logfile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("Failed to open log file: %v", err)
	}

	log.AddHook(lfshook.NewHook(lfshook.PathMap{
		log.DebugLevel: logfile,
		log.InfoLevel:  logfile,
		log.WarnLevel:  logfile,
		log.ErrorLevel: logfile,
		log.FatalLevel: logfile,
		log.PanicLevel: logfile,
	}, &log.TextFormatter{
		DisableColors: true,
		FullTimestamp: true,
	}))
}

func main() {
	log.Info("initializing library")
	alexandria.Run(
		"./data/02-01-25.json",
	)
}
