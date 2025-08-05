package request 

import (
    "fmt"
	"time"
    "net/http"
	
    "github.com/AceStructor/healthcheck-backend/db"
)

func HTTPCheck(cfg db.Config, WarningLog *log.Logger, InfoLog *log.Logger) (db.Result, error) {
    InfoLog.Printf("Running http check for %v \n", cfg.Name)
    var res db.Result
    
    client := &http.Client{ Timeout: cfg.Timeout * time.Second }
    res.CheckedAt = time.Now()
    
	req, err := http.NewRequest(cfg.Method, cfg.Target, "")	
	if err != nil {
		res.Text = "Error in request parameters: " + err.Error()
		res.Status = false
		WarningLog.Printf("Error in request parameters: %v \n", err)
        return res, nil
	}
	if cfg.Headers != "" {
		req.Header = cfg.Headers
	}
	
	res, err := client.Do(req)
    if err != nil {
        res.Text = "Connection failed: " + err.Error()
        res.Status = false
        WarningLog.Printf("An error during http connection: %v \n", err)
        return res, nil
    }
    if resp.StatusCode == cfg.ExpectStatus {
        res.Status = true
    } else {
        res.Status = false
    }
    
    res.Text = resp.Status
    
    return res, nil
}
