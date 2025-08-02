package request 

import (
	"time"
	
    "github.com/AceStructor/healthcheck-backend/db"
    "github.com/go-openapi/errors"
)

func RunChecks() error {
    if configs, err := db.ReadConfig(); err != nil {
		return err
	}
	
	for _, config := range configs {
		
		if config.Type == "http" {
			err := httpCheck(config)
		} else if config.Type == "tls" {
			err := tlsCheck(config)
		} else if config.Type == "dns" {
			err := dnsCheck(config)
		} else {
			return errors.New(500, "Type not found")
		}
	}
	return nil
}
