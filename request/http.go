package request

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/AceStructor/healthcheck-backend/db"
)

func HTTPCheck(cfg *db.Config, WarningLog *log.Logger, InfoLog *log.Logger) (*db.Result, error) {
	InfoLog.Printf("Running http check for %v \n", cfg.Name)
	var res db.Result

	timeoutInterval, _ := time.ParseDuration(fmt.Sprintf("%ds", cfg.Timeout))
	client := &http.Client{Timeout: timeoutInterval}
	res.CheckedAt = time.Now()

	var body io.Reader
	req, err := http.NewRequest(cfg.Method, cfg.Target, body)
	if err != nil {
		res.Text = "Error in request parameters: " + err.Error()
		res.Status = false
		WarningLog.Printf("Error in request parameters: %v \n", err)
		return &res, nil
	}
	if cfg.Headers != nil {
		for key, value := range *cfg.Headers {
			req.Header.Set(key, value)
		}
	}

	resp, err := client.Do(req)
	if err != nil {
		res.Text = "Connection failed: " + err.Error()
		res.Status = false
		WarningLog.Printf("An error during http connection: %v \n", err)
		return &res, nil
	}
	if resp.StatusCode == cfg.ExpectStatus {
		res.Status = true
	} else {
		res.Status = false
	}

	res.Text = resp.Status

	return &res, nil
}
