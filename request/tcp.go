package request 

import (
    "fmt"
	"time"
    "net"
	
    "github.com/AceStructor/healthcheck-backend/db"
)

func TCPCheck(cfg db.Config, WarningLog *log.Logger, InfoLog *log.Logger) (db.Result, error) {
    InfoLog.Printf("Running tcp check for %v \n", cfg.Name)
    var res db.Result
    
    host, _, err := net.SplitHostPort(cfg.Address)
    if err != nil {
        res.Text = "Address is not of form host:port: " + err.Error()
        res.Status = false
        WarningLog.Printf("Address is not of form host:port: %v \n", err)
        return res, nil
    }
    
    _, err := net.LookupHost(host)
    if err != nil {
        res.Text = "Error in Name resolution: " + err.Error()
        res.Status = false
        WarningLog.Printf("Error in Name resolution: %v \n", err)
        return res, nil
    }
    
    conn, err := net.DialTimeout("tcp", cfg.Address, cfg.Timeout*time.Second)
    if err != nil {
        res.Text = "Connection failed: " + err.Error()
        res.Status = false
        WarningLog.Printf("An error during http connection: %v \n", err)
        return res, nil
    }
    defer func() {
        if cerr = conn.Close(); cerr != nil {
            WarningLog.Printf("TCP Connection could not be closed: %v \n", cerr)
            res.Text += fmt.Sprintf(" (note: connection close failed: %v)", cerr)
        }
    }()

    res.Status = true
    res.Text = "TCP connection successful"
    return res, nil
}