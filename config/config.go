package config

import (
    "io/ioutil"
    "log"

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
}

func TranslateConfig(path String) ([]db.Config, error) {
    var configs []db.Config
    content, err := ioutil.ReadFile(path)
    if err != nil {
        return nil, err
    }
    
    var rawConfig RawConfig
    if err := yaml.Unmarshal(content, &rawConfig); err != nil {
        return nil, err
    }
    
    for _, svc := range rawConfig.Services {
        configs = append(configs, db.Config{
            Name:            svc.Name,
            Type:            svc.Type,
            Address:         svc.Address,
            IntervalSeconds: svc.Interval
        })
    }
    
    return configs, nil
}

func InitConfig() error {
    var configs []db.Config
    
    configs, err := TranslateConfig("config/exampleConf.yaml")
    if err != nil {
        return err
    }
    
    if err := AddConfig(configs); err != nil {
        return err
    }
    
    return nil
}

func AddConfig(configs []db.Config) error {
    for config := range configs {
        if err := db.WriteConfig(config); err != nil {
            return err
        }
    }
    
    return nil
}