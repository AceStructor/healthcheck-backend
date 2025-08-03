package request 

import (
    "fmt"
    "log"
	"time"
	
    "github.com/AceStructor/healthcheck-backend/db"
)

func RunChecks(WarningLog *log.Logger, InfoLog *log.Logger) error {
    if cfgs, err := db.ReadConfig(); err != nil {
		return fmt.Errorf("Error reading configs from database: %w", err)
	}
	
	for cfg := range cfgs {
		InfoLog.Printfln("Executing check for config %v", cfg.ID)
		tStart := time.Now()
        if tStart < cfg.LastChecked.Add(time.Second * cfg.IntervalSeconds) {
            continue
        }
        
        var (
            res db.Result
            err error
        )
        
		switch cfg.Type {
        case "http":
            res, err = HTTPCheck(cfg)
        case "tls":
            res, err = TLSCheck(cfg)
        case "dns":
            res, err = DNSCheck(cfg)
        default:
            WarningLog.Printfln("unsupported check type: %v", cfg.Type)
        }
        if err != nil {
            return fmt.Errorf("Error during Check: %w", err)
        }
        duration := time.Since(tStart)
        
        cfg.LastChecked = tStart
        res.ResponseTime = duration
        
        if err := db.UpdateConfig(cfg); err != nil {
            return fmt.Errorf("Error writing execution time to database: %w", err)
        }
        
        if err := db.WriteResult(res, cfg.ID); err != nil {
            return fmt.Errorf("Error writing new result: %w", err)
        }
	}
	return nil
}
