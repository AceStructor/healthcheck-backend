package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/AceStructor/healthcheck-backend/config"
	"github.com/AceStructor/healthcheck-backend/db"
	"github.com/AceStructor/healthcheck-backend/request"
)

var (
	WarningLog *log.Logger
	InfoLog    *log.Logger
	ErrorLog   *log.Logger
)

func initLogging() {
	file, err := os.OpenFile("log.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		if err = file.Close(); err != nil {
			log.Fatal(err)
		}
	}()

	InfoLog = log.New(file, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile)
	WarningLog = log.New(file, "WARNING: ", log.Ldate|log.Ltime|log.Lshortfile)
	ErrorLog = log.New(file, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile)
}

func main() {
	initLogging()
	InfoLog.Println("Starting healthcheck backend")

	if err := db.InitDB(WarningLog, InfoLog); err != nil {
		ErrorLog.Printf("Database init failed: %v \n", err)
		os.Exit(1)
	}

	if err := config.InitConfig(WarningLog, InfoLog); err != nil {
		ErrorLog.Printf("Config init failed: %v \n", err)
		os.Exit(1)
	}

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			if err := request.RunChecks(WarningLog, InfoLog); err != nil {
				ErrorLog.Printf("Check run failed: %v \n", err)
			}
		case <-stop:
			InfoLog.Println("Received termination signal, exiting...")
			return
		}
	}
}
