package config

import (
    "os"
	"fmt"
	"io"
	"log"

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
	Name         string  `yaml:"name"`
	Target       string  `yaml:"url"`
	Interval     int     `yaml:"interval"` // in seconds
	Timeout      int     `yaml:"timeout"`
	Method       *string `yaml:"method,omitempty"`
	Headers      *string `yaml:"headers,omitempty"`
	ExpectStatus *int    `yaml:"expect_status,omitempty"`
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

func NewHTTPConfig(httpc RawHTTPElement) db.Config {
	return db.Config{
		Type:            "http",
		Name:            httpc.Name,
		Target:          httpc.Target,
		IntervalSeconds: httpc.Interval,
		Timeout:         httpc.Timeout,
		Method:          helper.StringOrDefault(httpc.Method, "GET"),
		Headers:         helper.StringOrDefault(httpc.Headers, ""),
		ExpectStatus:    helper.IntOrDefault(httpc.ExpectStatus, 200),
	}
}

func NewTCPConfig(tcpc RawTCPElement) db.Config {
	return db.Config{
		Type:            "tcp",
		Name:            tcpc.Name,
		Target:          tcpc.Target,
		Port:            tcpc.Port,
		IntervalSeconds: tcpc.Interval,
		Timeout:         tcpc.Timeout,
	}
}

func NewDNSConfig(dnsc RawDNSElement) db.Config {
	return db.Config{
		Type:            "dns",
		Name:            dnsc.Name,
		Target:          dnsc.Target,
		IntervalSeconds: dnsc.Interval,
		RecordType:      helper.StringOrDefault(dnsc.RecordType, "A"),
		ExpectIP:        helper.StringOrDefault(dnsc.ExpectIP, ""),
		DNSServer:       helper.StringOrDefault(dnsc.DNSServer, "1.1.1.1"),
	}
}

func TranslateConfig(path string, WarningLog *log.Logger, InfoLog *log.Logger) ([]db.Config, error) {
	InfoLog.Println("Starting Config Translation...")
	var cfgs []db.Config
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

	for _, httpc := range rawConfig.HTTP {
		cfgs = append(cfgs, NewHTTPConfig(httpc))
	}
	for _, tcpc := range rawConfig.TCP {
		cfgs = append(cfgs, NewTCPConfig(tcpc))
	}
	for _, dnsc := range rawConfig.DNS {
		cfgs = append(cfgs, NewDNSConfig(dnsc))
	}

	InfoLog.Println("Config Translation Successful!")
	return cfgs, nil
}

func InitConfig(WarningLog *log.Logger, InfoLog *log.Logger) error {
	InfoLog.Println("Initializing Configuration...")
	var cfgs []db.Config

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

func AddConfig(cfgs []db.Config, WarningLog *log.Logger, InfoLog *log.Logger) error {
	for _, cfg := range cfgs {
		if err := db.WriteConfig(cfg, WarningLog, InfoLog); err != nil {
			return fmt.Errorf("Error while writing config %v: %w", cfg.ID, err)
		}
	}

	return nil
}
