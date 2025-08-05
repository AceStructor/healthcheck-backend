package request 

import (
    "fmt"
	"time"
    "net"
	
    "github.com/AceStructor/healthcheck-backend/db"
)

func DNSCheck(cfg db.Config, WarningLog *log.Logger, InfoLog *log.Logger) (db.Result, error) {
    InfoLog.Printf("Running DNS check for %v \n", cfg.Name)
    var res db.Result
    
    resp, err := net.LookupHost(cfg.Target)
    if err != nil {
        res.Text = "Error in Name resolution: " + err.Error()
        res.Status = false
        WarningLog.Printf("Error in Name resolution: %v \n", err)
        return res, nil
    }
    
    res.Status = true
    res.Text = fmt.Printf("DNS resolution was successful: %v", resp)
    return res, nil
}
