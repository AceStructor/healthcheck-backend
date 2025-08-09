package request

import (
    "log"
    "testing"
    
    "github.com/google/go-cmp/cmp"
    "github.com/AceStructor/healthcheck-backend/db"
)

var aConfig *db.Config = &db.Config{
                Type: "dns",
                Name: "Example",
                Target: "example.com",
                IntervalSeconds: 300,
                RecordType: "A",
                ExpectIP: "23.192.228.80",
                DNSServer: "8.8.8.8:53",
            }

var mxConfig *db.Config = &db.Config{
                Type: "dns",
                Name: "Wikipedia",
                Target: "wikipedia.org",
                IntervalSeconds: 300,
                RecordType: "MX",
                ExpectIP: "mx-in1001.wikimedia.org.",
                DNSServer: "8.8.8.8:53",
            }
            
var unexpectedIPConfig *db.Config = &db.Config{
                Type: "dns",
                Name: "Example",
                Target: "example.com",
                IntervalSeconds: 300,
                RecordType: "A",
                ExpectIP: "1.1.1.1",
                DNSServer: "8.8.8.8:53",
            }

var invalidConfig *db.Config = &db.Config{
                Type: "dns",
                Name: "Invalid",
                Target: "invalid.local",
                IntervalSeconds: 300,
                RecordType: "AAAA",
                ExpectIP: "",
                DNSServer: "8.8.8.8:53",
            }

var successResult *db.Result = &db.Result{Text: "DNS resolution was successful", Status: true}

var failureResult *db.Result = &db.Result{Text: "DNS query failed with code: 3"}

var ipNotFoundResult *db.Result = &db.Result{Text: "Error in Name resolution: Name was not resolved to the expected IP 1.1.1.1"}

func TestDNSCheck(t *testing.T) {
    tests := []struct {
        name     string
        input       *db.Config
        expectedErr error
        expectedRes *db.Result
    }{
        {
            name:     "Valid A Hostname",
            input:    aConfig,
            expectedErr: nil,
            expectedRes: successResult,
        },
        {
            name:     "Valid MX Hostname",
            input:    mxConfig,
            expectedErr: nil,
            expectedRes: successResult,
        },
        {
            name:     "Unexpected IP",
            input:    unexpectedIPConfig,
            expectedErr: nil,
            expectedRes: ipNotFoundResult,
        },
        {
            name:     "Invalid Hostname",
            input:    invalidConfig,
            expectedErr: nil,
            expectedRes: failureResult,
        },
    }
    
    for _, tt := range tests {
        tt := tt 
        t.Run(tt.name, func(t *testing.T) {
            t.Parallel() 

            result, err := DNSCheck(tt.input, log.Default(), log.Default())
            if (err == nil) != (tt.expectedErr == nil) {
                t.Errorf("Expected error: %v, got: %v", tt.expectedErr, err)
            } else if err != nil && err.Error() != tt.expectedErr.Error() {
                t.Errorf("Expected %q, got %q", tt.expectedErr.Error(), err.Error())
            }
            if diff := cmp.Diff(*tt.expectedRes, *result); diff != "" {
                t.Errorf("Mismatch (-expected +got):\n%s", diff)
            }
        })
    }
}