package config

import (
    "io/ioutil"
    "log"
    "fmt"

    "gopkg.in/yaml.v3"
    "github.com/AceStructor/healthcheck-backend/db"
    "github.com/AceStructor/healthcheck-backend/helper"
)

type RawConfig struct {
	HTTP     []RawHTTPElement `yaml:"http"`
	TLS      []RawTLSElement `yaml:"tls"`
	DNS      []RawDNSElement `yaml:"dnss"`
}

type RawHTTPElement struct {
    Name     string `yaml:"name"`
    Target   string `yaml:"url"`
    Interval int    `yaml:"interval"` // in seconds
    Timeout	 int    `yaml:"timeout"`
    Method   *string `yaml:"method,omitempty"`
    Headers  *string `yaml:"headers,omitempty"`
    ExpectStatus *int `yaml:"expect_status,omitempty"`
}

type RawTLSElement struct {
    Name     string `yaml:"name"`
    Target   string `yaml:"host"`
    Port     string `yaml:"port"`
    Interval int    `yaml:"interval"` // in seconds
    Timeout	 int    `yaml:"timeout"`
}

type RawDNSElement struct {
    Name     string `yaml:"name"`
    Target   string `yaml:"hostname"`
    Interval int    `yaml:"interval"` // in seconds
    RecordType *string `yaml:"record_type,omitempty"`
    ExpectIP *string `yaml:"expect_ip,omitempty"`
    DNSServer *string `yaml:"dns_server,omitempty"`
}

func NewHTTPConfig(httpc RawHTTPElement) db.Config {
	return db.Config{
		Type:            "http"
		Name:            httpc.Name,
		Target:          httpc.Target,
		IntervalSeconds: httpc.Interval,
		Timeout:		 httpc.Timeout,
		Method:          helper.stringOrDefault(httpc.Method, "GET"),
		Headers:         helper.stringOrDefault(httpc.Headers, ""),
		ExpectStatus:    helper.intOrDefault(httpc.ExpectStatus, 200),
	}
}

func NewTLSConfig(tlsc RawTLSElement) db.Config {
	return db.Config{
		Type:            "tls"
		Name:            tlsc.Name,
		Target:          tlsc.Target,
		Port:            tlsc.Port,
		IntervalSeconds: tlsc.Interval,
		Timeout:		 tlsc.Timeout
	}
}

func NewDNSConfig(dnsc RawDNSElement) db.Config {
	return db.Config{
		Type:            "dns"
		Name:            dnsc.Name,
		Target:          dnsc.Target,
		IntervalSeconds: dnsc.Interval,
		RecordType:		 helper.stringOrDefault(dnsc.RecordType, "A"),
		ExpectIP:	     helper.stringOrDefault(dnsc.ExpectIP, "")
        DNSServer:       helper.stringOrDefault(dnsc.DNSServer, "1.1.1.1")
	}
}

func TranslateConfig(path String, WarningLog *log.Logger, InfoLog *log.Logger) ([]db.Config, error) {
    InfoLog.Println("Starting Config Translation...")
    var cfgs []db.Config
    file, err := os.OpenFile(path)
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
    for _, tlsc := range rawConfig.TLS {
        cfgs = append(cfgs, NewTLSConfig(tlsc))
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
        return fmt.Errorf("Error while translating config: %w", err
    }
    
    if err := AddConfig(cfgs); err != nil {
        return fmt.Errorf("Error while adding configs to database: %w", err)
    }
    
    InfoLog.Println("Configuration Initialized!")
    return nil
}

func AddConfig(cfgs []db.Config) error {
    for cfg := range cfgs {
        if err := db.WriteConfig(cfg); err != nil {
            return fmt.Errorf("Error while writing config %v: %w", cfg.ID, err)
        }
    }
    
    return nil
}
