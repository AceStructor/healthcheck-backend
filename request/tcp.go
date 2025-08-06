package request

import (
	"fmt"
	"log"
	"net"
	"time"

	"github.com/AceStructor/healthcheck-backend/db"
)

func TCPCheck(cfg *db.Config, WarningLog *log.Logger, InfoLog *log.Logger) (*db.Result, error) {
	InfoLog.Printf("Running tcp check for %v \n", cfg.Name)
	var res *db.Result

	timeoutInterval, _ := time.ParseDuration(fmt.Sprintf("%ds", cfg.Timeout))
	conn, err := net.DialTimeout("tcp", fmt.Sprintf("%s:%d", cfg.Target, cfg.Port), timeoutInterval)
	if err != nil {
		res.Text = "Connection failed: " + err.Error()
		res.Status = false
		WarningLog.Printf("An error during http connection: %v \n", err)
		return res, nil
	}
	defer func() {
		if cerr := conn.Close(); cerr != nil {
			WarningLog.Printf("TCP Connection could not be closed: %v \n", cerr)
			res.Text += fmt.Sprintf(" (note: connection close failed: %v)", cerr)
		}
	}()

	res.Status = true
	res.Text = "TCP connection successful"
	return res, nil
}
