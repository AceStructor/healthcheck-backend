package db

import (
    "github.com/go-openapi/errors"
)

func WriteConfig(cfg Config) error {
    if err := DB.Create(&cfg).Error; err != nil {
        return errors.New(500, "failed to create configuration entry: %v", err)
    }
    
    return nil
}

func ReadConfig() ([]Config, error) {
    var cfgs []Config
    
    if err := DB.Where("disabled = 0").Find(&cfgs).Error; err != nil {
        return nil, errors.New(500, "failed to read configs from database: %v", err)
    }
    
    return cfgs, nil
}

func DisableConfig(cfg Config) error {
	var targetConfig Config
	
	if err := DB.First(&targetConfig, cfg.ID).Error; err != nil {
        return errors.New(500, "failed to find config in database: %v", err)
    }
    
    if err := DB.Model(&targetConfig).Updates(Config{Disabled: true}).Error; err != nil {
		return errors.New(500, "failed to update config: %v", err)
	}

}

func UpdateConfig(cfg Config) error {
	var targetConfig Config
	
	if err := DB.First(&targetConfig, cfg.ID).Error; err != nil {
        return errors.New(500, "failed to find config in database: %v", err)
    }
    
    if err := DB.Model(&targetConfig).Updates(Config{Name: cfg.Name, Type: cfg.Type, Address: cfg.Address, IntervalSeconds: cfg.IntervalSeconds, LastChecked: cfg.LastChecked, Timeout: cfg.Timeout}).Error; err != nil {
		return errors.New(500, "failed to update config: %v", err)
	}
}
