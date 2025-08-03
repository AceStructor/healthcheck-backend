package config

import (
    "io/ioutil"
    "log"
    "fmt"

    "gopkg.in/yaml.v3"
    "github.com/AceStructor/healthcheck-backend/db"
)

type RawConfig struct {
    Services []RawConfigElement `yaml:"services"`
}

type RawConfigElement struct {
    Name     string `yaml:"name"`
    Type     string `yaml:"type"`
    Address  string `yaml:"address"`
    Interval int    `yaml:"interval"` // in seconds
    Timeout	 int    `yaml:"timeout"`
}

func TranslateConfig(path String, WarningLog *log.Logger, InfoLog *log.Logger) ([]db.Config, error) {
    InfoLog.Println("Starting Config Translation...")
    var cfgs []db.Config
    content, err := ioutil.ReadFile(path)
    if err != nil {
        return nil, fmt.Errorf("Error while reading config at %v: %w", path, err)
    }
    
    var rawConfig RawConfig
    if err := yaml.Unmarshal(content, &rawConfig); err != nil {
        return nil, fmt.Errorf("Error while unmarshaling yaml: %w", err)
    }
    
    for _, svc := range rawConfig.Services {
        cfgs = append(cfgs, db.Config{
            Name:            svc.Name,
            Type:            svc.Type,
            Address:         svc.Address,
            IntervalSeconds: svc.Interval,
            Timeout:		 svc.Timeout
        })
    }
    
    InfoLog.Println("Config Translation Successful!")
    return cfgs, nil
}

func InitConfig(WarningLog *log.Logger, InfoLog *log.Logger) error {
    InfoLog.Println("Initializing Configuration...")
    var cfgs []db.Config
    
    cfgs, err := TranslateConfig("config/exampleConf.yaml")
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
