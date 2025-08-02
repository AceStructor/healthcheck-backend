package db

import (
    "github.com/go-openapi/errors"
)

func WriteConfig(config Config) error {
    if err := DB.Create(&config).Error; err != nil {
        return errors.New(500, "failed to create configuration entry: %v", err)
    }
    
    return nil
}

func ReadConfig() ([]Config, error) {
    var configs []Config
    
    if err := DB.Where("disabled = 0").Find(&configs).Error; err != nil {
        return nil, errors.New(500, "failed to read configs from database: %v", err)
    }
    
    return configs, nil
}

func DisableConfig(config Config) error {
	var targetConfig Config
	
	if err := DB.First(&targetConfig, config.ID).Error; err != nil {
        return errors.New(500, "failed to find config in database: %v", err)
    }
    
    if err := DB.Model(&targetConfig).Updates(Config{Disabled: true}).Error; err != nil {
		return errors.New(500, "failed to update config: %v", err)
	}

}

func UpdateConfig(config Config) error {
	var targetConfig Config
	
	if err := DB.First(&targetConfig, config.ID).Error; err != nil {
        return errors.New(500, "failed to find config in database: %v", err)
    }
    
    if err := DB.Model(&targetConfig).Updates(Config{Name: config.Name, Type: config.Type, Address: config.Address, IntervalSeconds: config.IntervalSeconds, LastChecked: config.LastChecked}).Error; err != nil {
		return errors.New(500, "failed to update config: %v", err)
	}
}
