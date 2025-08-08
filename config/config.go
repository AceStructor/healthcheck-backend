package config

import (
	"fmt"
	"io"
	"log"
	"os"
    "strings"

	"github.com/AceStructor/healthcheck-backend/db"
	"github.com/AceStructor/healthcheck-backend/helper"
	"gopkg.in/yaml.v3"
)

type RawConfig struct {
	HTTP []RawHTTPElement `yaml:"http"`
	TCP  []RawTCPElement  `yaml:"tcp"`
	DNS  []RawDNSElement  `yaml:"dns"`
}

type RawHTTPElement struct {
	Name         string             `yaml:"name"`
	Target       string             `yaml:"url"`
	Interval     int                `yaml:"interval"` // in seconds
	Timeout      int                `yaml:"timeout"`
	Method       *string            `yaml:"method,omitempty"`
	Headers      *map[string]string `yaml:"headers,omitempty"`
	ExpectStatus *int               `yaml:"expect_status,omitempty"`
}

type RawTCPElement struct {
	Name     string `yaml:"name"`
	Target   string `yaml:"host"`
	Port     int    `yaml:"port"`
	Interval int    `yaml:"interval"` // in seconds
	Timeout  int    `yaml:"timeout"`
}

type RawDNSElement struct {
	Name       string  `yaml:"name"`
	Target     string  `yaml:"hostname"`
	Interval   int     `yaml:"interval"` // in seconds
	RecordType *string `yaml:"record_type,omitempty"`
	ExpectIP   *string `yaml:"expect_ip,omitempty"`
	DNSServer  *string `yaml:"dns_server,omitempty"`
}

func NewHTTPConfig(httpc RawHTTPElement) (*db.Config, error) {
	cfg := db.Config{
		Type:            "http",
		Name:            httpc.Name,
		Target:          httpc.Target,
		IntervalSeconds: httpc.Interval,
		Timeout:         httpc.Timeout,
		Method:          helper.StringOrDefault(httpc.Method, "GET"),
		Headers:         httpc.Headers,
		ExpectStatus:    helper.IntOrDefault(httpc.ExpectStatus, 200),
	}
    if err := helper.ValidateConfig("http", cfg); err != nil {
        return nil, err
    }
	return &cfg, nil
}

func NewTCPConfig(tcpc RawTCPElement) (*db.Config, error) {
	cfg := db.Config{
		Type:            "tcp",
		Name:            tcpc.Name,
		Target:          tcpc.Target,
		Port:            tcpc.Port,
		IntervalSeconds: tcpc.Interval,
		Timeout:         tcpc.Timeout,
	}
    if err := helper.ValidateConfig("tcp", cfg); err != nil {
        return nil, err
    }
	return &cfg, nil
}

func NewDNSConfig(dnsc RawDNSElement) (*db.Config, error) {
	cfg := db.Config{
		Type:            "dns",
		Name:            dnsc.Name,
		Target:          dnsc.Target,
		IntervalSeconds: dnsc.Interval,
		RecordType:      helper.StringOrDefault(dnsc.RecordType, "A"),
		ExpectIP:        helper.StringOrDefault(dnsc.ExpectIP, ""),
		DNSServer:       helper.StringOrDefault(dnsc.DNSServer, "1.1.1.1"),
	}
    if err := helper.ValidateConfig("dns", cfg); err != nil {
        return nil, err
    }
	return &cfg, nil
}

func TranslateConfig(path string, WarningLog *log.Logger, InfoLog *log.Logger) ([]*db.Config, error) {
	InfoLog.Println("Starting Config Translation...")
	var cfgs []*db.Config
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("Error while reading config at %v: %w", path, err)
	}
	defer func() {
		if cerr := file.Close(); cerr != nil {
			WarningLog.Printf("Config file could not be closed: %v", cerr)
		}
	}()
	content, err := io.ReadAll(file)
	if err != nil {
		return nil, fmt.Errorf("Error while reading config at %v: %w", path, err)
	}

	var rawConfig RawConfig
	if err := yaml.Unmarshal(content, &rawConfig); err != nil {
		return nil, fmt.Errorf("Error while unmarshaling yaml: %w", err)
	}

    var allErrs []string
	for i, httpc := range rawConfig.HTTP {
        cfg, err := NewHTTPConfig(httpc)
        if err != nil {
            allErrs = append(allErrs, fmt.Sprintf("HTTP config %d: %v", i+1, err))
        }
		cfgs = append(cfgs, cfg)
	}
	for i, tcpc := range rawConfig.TCP {
        cfg, err := NewTCPConfig(tcpc)
        if err != nil {
            allErrs = append(allErrs, fmt.Sprintf("TCP config %d: %v", i+1, err))
        }
		cfgs = append(cfgs, cfg)
	}
	for i, dnsc := range rawConfig.DNS {
        cfg, err := NewDNSConfig(dnsc)
        if err != nil {
            allErrs = append(allErrs, fmt.Sprintf("DNS config %d: %v", i+1, err))
        }
		cfgs = append(cfgs, cfg)
	}
    
    if len(allErrs) > 0 {
        return nil, fmt.Errorf("validation failed:\n%s", strings.Join(allErrs, "\n"))
    }

	InfoLog.Println("Config Translation Successful!")
	return cfgs, nil
}

func InitConfig(WarningLog *log.Logger, InfoLog *log.Logger) error {
	InfoLog.Println("Initializing Configuration...")
	var cfgs []*db.Config

	cfgs, err := TranslateConfig("config/exampleConf.yaml", WarningLog, InfoLog)
	if err != nil {
		return fmt.Errorf("Error while translating config: %w", err)
	}

	if err := AddConfig(cfgs, WarningLog, InfoLog); err != nil {
		return fmt.Errorf("Error while adding configs to database: %w", err)
	}

	InfoLog.Println("Configuration Initialized!")
	return nil
}

func AddConfig(cfgs []*db.Config, WarningLog *log.Logger, InfoLog *log.Logger) error {
	for _, cfg := range cfgs {
		if err := db.WriteConfig(cfg, WarningLog, InfoLog); err != nil {
			return fmt.Errorf("Error while writing config %v: %w", cfg.ID, err)
		}
	}

	return nil
}
