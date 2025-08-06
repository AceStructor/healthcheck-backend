package request

import (
	"fmt"
	"log"
	"time"

	"github.com/AceStructor/healthcheck-backend/db"
)

func RunChecks(WarningLog *log.Logger, InfoLog *log.Logger) error {
	cfgs, err := db.ReadConfig(WarningLog, InfoLog)
	if err != nil {
		return fmt.Errorf("Error reading configs from database: %w", err)
	}

	for _, cfg := range cfgs {
		InfoLog.Printf("Executing check for config %v \n", cfg.ID)
		tStart := time.Now()
		skipInterval, _ := time.ParseDuration(fmt.Sprintf("%ds", cfg.IntervalSeconds))
		if tStart.Before(cfg.LastChecked.Add(skipInterval)) {
			continue
		}

		var (
			res *db.Result
			err error
		)

		switch cfg.Type {
		case "http":
			res, err = HTTPCheck(cfg, WarningLog, InfoLog)
		case "tcp":
			res, err = TCPCheck(cfg, WarningLog, InfoLog)
		case "dns":
			res, err = DNSCheck(cfg, WarningLog, InfoLog)
		default:
			WarningLog.Printf("unsupported check type: %v \n", cfg.Type)
		}
		if err != nil {
			return fmt.Errorf("Error during Check: %w", err)
		}
		if res == nil {
			return fmt.Errorf("Error during Check: No result was created!")
		}
		duration := time.Since(tStart)

		cfg.LastChecked = tStart
		res.ResponseTime = duration

		if err := db.UpdateConfig(cfg, WarningLog, InfoLog); err != nil {
			return fmt.Errorf("Error writing execution time to database: %w", err)
		}

		if err := db.WriteResult(res, cfg.ID, WarningLog, InfoLog); err != nil {
			return fmt.Errorf("Error writing new result: %w", err)
		}
	}
	return nil
}
