package config

import (
    "io/ioutil"
    "log"
    "fmt"

    "gopkg.in/yaml.v3"
    "github.com/AceStructor/healthcheck-backend/db"
)

type RawConfig struct {
	HTTP     []RawHTTPElement `yaml:"http"`
	TLS      []RawTLSElement `yaml:"tls"`
	DNS      []RawDNSElement `yaml:"dnss"`
}

type RawHTTPElement struct {
    Name     string `yaml:"name"`
    URL      string `yaml:"url"`
    Interval int    `yaml:"interval"` // in seconds
    Timeout	 int    `yaml:"timeout"`
    Method   string `yaml:"method"`
    Headers  string `yaml:"headers"`
    ExpectStatus string `yaml:"expect_status"`
}

type RawTLSElement struct {
    Name     string `yaml:"name"`
    Host     string `yaml:"host"`
    Port     string `yaml:"port"`
    Interval int    `yaml:"interval"` // in seconds
    Timeout	 int    `yaml:"timeout"`
}

type RawDNSElement struct {
    Name     string `yaml:"name"`
    Hostname string `yaml:"hostname"`
    Interval int    `yaml:"interval"` // in seconds
    RecordType string `yaml:"record_type"`
    ExpectIP string `yaml:"expect_IP"`
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
        cfgs = append(cfgs, db.Config{
            Name:            httpc.Name,
            Address:         httpc.Address,
            IntervalSeconds: httpc.Interval,
            Timeout:		 httpc.Timeout
        })
    }
    for _, tlsc := range rawConfig.TLS {
        cfgs = append(cfgs, db.Config{
            Name:            tlsc.Name,
            Address:         tlsc.Address,
            IntervalSeconds: tlsc.Interval,
            Timeout:		 tlsc.Timeout
        })
    }
    for _, dnsc := range rawConfig.DNS {
        cfgs = append(cfgs, db.Config{
            Name:            dnsc.Name,
            Address:         dnsc.Address,
            IntervalSeconds: dnsc.Interval,
            Timeout:		 dnsc.Timeout
        })
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
