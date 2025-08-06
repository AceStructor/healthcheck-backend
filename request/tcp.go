package request

import (
	"fmt"
	"net"
    "log"
	"time"

	"github.com/AceStructor/healthcheck-backend/db"
)

func TCPCheck(cfg db.Config, WarningLog *log.Logger, InfoLog *log.Logger) (db.Result, error) {
	InfoLog.Printf("Running tcp check for %v \n", cfg.Name)
	var res db.Result

	conn, err := net.DialTimeout("tcp", cfg.Target+":"+cfg.Port, cfg.Timeout*time.Second)
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
