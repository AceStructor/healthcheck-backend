package request 

import (
    "fmt"
	"time"
    "net/http"
	
    "github.com/AceStructor/healthcheck-backend/db"
)

func HTTPCheck(cfg db.Config) (db.Result, error) {
    var res db.Result
    
    client := &http.Client{ Timeout: cfg.Timeout * time.Second }
    res.CheckedAt = time.Now()
    
    resp, err := client.Get(cfg.Address)
    if err != nil {
        res.Text = "Connection failed: " + err.Error()
        res.Status = false
        return res, nil
    }
    if resp.StatusCode >= 200 && resp.StatusCode < 300 {
        res.Status = true
    } else {
        res.Status = false
    }
    
    res.Text = resp.Status
    
    return res, nil
}