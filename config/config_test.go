package config

import (
    "log"
    "fmt"
    "testing"
    
    "github.com/google/go-cmp/cmp"
    "github.com/AceStructor/healthcheck-backend/db"
    "github.com/AceStructor/healthcheck-backend/helper"
)

var minCfg []*db.Config = []*db.Config{
            {
                Type: "http",
                Name: "Frontend",
                Target: "https://example.com",
                IntervalSeconds: 60,
                Timeout: 5,
                Method: "GET",
                ExpectStatus: 200,
            },
            {
                Type: "tcp",
                Name: "API",
                Target: "127.0.0.1",
                Port: 3000,
                IntervalSeconds: 30,
                Timeout: 10,
            },
            {
                Type: "dns",
                Name: "Server",
                Target: "example.com",
                IntervalSeconds: 300,
                RecordType: "A",
                DNSServer: "1.1.1.1",
            },
        }
        
var fullCfg []*db.Config = []*db.Config{
            {
                Type: "http",
                Name: "Frontend",
                Target: "https://example.com",
                IntervalSeconds: 60,
                Timeout: 5,
                Method: "GET",
                Headers: &map[string]string{
                    "Authorization": "Bearer xyz123",
                    "Content-Type": "application/json",
                },
                ExpectStatus: 200,
            },
            {
                Type: "tcp",
                Name: "API",
                Target: "127.0.0.1",
                Port: 3000,
                IntervalSeconds: 30,
                Timeout: 10,
            },
            {
                Type: "dns",
                Name: "Server",
                Target: "example.com",
                IntervalSeconds: 300,
                RecordType: "A",
                ExpectIP: "192.168.1.20",
                DNSServer: "8.8.8.8",
            },
        }
        
var wrongConfigErr error = fmt.Errorf("validation failed:\nHTTP config 1: missing required fields in config 'Frontend': [timeout]\nTCP config 1: missing required fields in config 'API': [timeout]\nDNS config 1: missing required fields in config 'Server': [interval]")

func TestTranslateConfig(t *testing.T) {
    tests := []struct {
        name     string
        path     string
        expectedErr error
        expectedRes []*db.Config
    }{
        {
            name:     "full Configuration",
            path:     "./exampleConf.yaml",
            expectedErr: nil,
            expectedRes: fullCfg,
        },
        {
            name:     "minimal Configuration",
            path:     "./minimalExampleConf.yaml",
            expectedErr: nil,
            expectedRes: minCfg,
        },
        {
            name:     "invalid Configuration",
            path:     "./invalidConf.yaml",
            expectedErr: wrongConfigErr,
            expectedRes: nil,
        },
    }
    
    for _, tt := range tests {
        tt := tt 
        t.Run(tt.name, func(t *testing.T) {
            t.Parallel() 

            result, err := TranslateConfig(tt.path, log.Default(), log.Default())
            if err != nil {
                if err.Error() != tt.expectedErr.Error() {
                    t.Errorf("Expected %q, got %q", tt.expectedErr.Error(), err.Error())
                }
            }
            if diff := cmp.Diff(helper.DerefConfigs(tt.expectedRes), helper.DerefConfigs(result)); diff != "" {
                t.Errorf("Mismatch (-expected +got):\n%s", diff)
            }
        })
    }
}