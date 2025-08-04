package request 

import (
    "fmt"
	"time"
    "net"
	
    "github.com/AceStructor/healthcheck-backend/db"
)

func TCPCheck(cfg db.Config, WarningLog *log.Logger, InfoLog *log.Logger) (db.Result, error) {
    InfoLog.Printfln("Running tcp check for %v", cfg.Name)
    var res db.Result
    
    host, _, err := net.SplitHostPort(cfg.Address)
    if err != nil {
        res.Text = "Address is not of form host:port: " + err.Error()
        res.Status = false
        WarningLog.Printfln("Address is not of form host:port: %v", err)
        return res, nil
    }
    
    _, err := net.LookupHost(host)
    if err != nil {
        res.Text = "Error in Name resolution: " + err.Error()
        res.Status = false
        WarningLog.Printfln("Error in Name resolution: %v", err)
        return res, nil
    }
    
    conn, err := net.DialTimeout("tcp", cfg.Address, cfg.Timeout*time.Second)
    if err != nil {
        res.Text = "Connection failed: " + err.Error()
        res.Status = false
        WarningLog.Printfln("An error during http connection: %v", err)
        return res, nil
    }
    defer conn.Close()

    res.Status = true
    res.Text = "TCP connection successful"
    return res, nil
}