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
    
    if err := DB.Where("disabled = false").Find(&configs).Error; err != nil {
        return errors.New(500, "failed to read configs from database: %v", err)
    }
}

func UpdateConfig(config Config) error {

}