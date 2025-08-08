package helper

import (
    "fmt"
    
	"github.com/AceStructor/healthcheck-backend/db"
)

func StringOrDefault(ptr *string, def string) string {
	if ptr != nil {
		return *ptr
	}
	return def
}

func IntOrDefault(ptr *int, def int) int {
	if ptr != nil {
		return *ptr
	}
	return def
}

func DerefConfigs(input []*db.Config) []db.Config {
	result := make([]db.Config, len(input))
	for i, ptr := range input {
		copy := *ptr

		if ptr.Headers != nil {
			headersCopy := make(map[string]string, len(*ptr.Headers))
			for k, v := range *ptr.Headers {
				headersCopy[k] = v
			}
			copy.Headers = &headersCopy
		}

		result[i] = copy
	}
	return result
}

func DerefResults(input []*db.Result) []db.Result {
	result := make([]db.Result, len(input))
	for i, ptr := range input {
		result[i] = *ptr
	}
	return result
}

func ValidateConfig(cType string, cfg db.Config) error {
    var missing []string

    if cfg.Name == "" {
        missing = append(missing, "name")
    }
    if cfg.Target == "" {
        switch cType {
        case "http": 
            missing = append(missing, "url")
        case "tcp": 
            missing = append(missing, "host")
        case "dns": 
            missing = append(missing, "hostname")
        default: 
            missing = append(missing, "target")
        }
    }
    if cfg.IntervalSeconds == 0 {
        missing = append(missing, "interval")
    }
    if (cType == "http" || cType == "tcp") && cfg.Timeout == 0 {
        missing = append(missing, "timeout")
    }
    if cType == "tcp" && cfg.Port == 0 {
        missing = append(missing, "port")
    }
    
    if len(missing) > 0 {
        return fmt.Errorf("missing required fields in config '%s': %v", cfg.Name, missing)
    }
    
    return nil
}