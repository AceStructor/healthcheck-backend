package db

import (
	"fmt"
	"log"
)

func WriteConfig(cfg *Config, WarningLog *log.Logger, InfoLog *log.Logger) error {
	InfoLog.Printf("Creating Config %v (%v) \n", cfg.Name, cfg.Target)
	if err := DB.Create(&cfg).Error; err != nil {
		return fmt.Errorf("failed to create configuration entry: %w", err)
	}

	return nil
}

func ReadConfig(WarningLog *log.Logger, InfoLog *log.Logger) ([]*Config, error) {
	InfoLog.Println("Reading Configs")
	var cfgs []*Config

	if err := DB.Where("disabled = 0").Find(&cfgs).Error; err != nil {
		return nil, fmt.Errorf("failed to read configs from database: %w", err)
	}

	return cfgs, nil
}

func DisableConfig(cfg *Config, WarningLog *log.Logger, InfoLog *log.Logger) error {
	InfoLog.Printf("Disabling Config %v \n", cfg.ID)
	var targetConfig *Config

	if err := DB.First(&targetConfig, cfg.ID).Error; err != nil {
		return fmt.Errorf("failed to find config in database: %w", err)
	}

	if err := DB.Model(&targetConfig).Updates(Config{Disabled: true}).Error; err != nil {
		return fmt.Errorf("failed to update config: %w", err)
	}

	return nil
}

func UpdateConfig(cfg *Config, WarningLog *log.Logger, InfoLog *log.Logger) error {
	InfoLog.Printf("Updating Config %v \n", cfg.ID)
	var targetConfig *Config

	if err := DB.First(&targetConfig, cfg.ID).Error; err != nil {
		return fmt.Errorf("failed to find config in database: %w", err)
	}

	if err := DB.Model(&targetConfig).Updates(Config{
		Name:            cfg.Name,
		Type:            cfg.Type,
		Target:          cfg.Target,
		IntervalSeconds: cfg.IntervalSeconds,
		LastChecked:     cfg.LastChecked,
		Timeout:         cfg.Timeout,
		Method:          cfg.Method,
		Headers:         cfg.Headers,
		ExpectStatus:    cfg.ExpectStatus,
		RecordType:      cfg.RecordType,
		ExpectIP:        cfg.ExpectIP,
		DNSServer:       cfg.DNSServer,
	}).Error; err != nil {
		return fmt.Errorf("failed to update config: %w", err)
	}

	return nil
}
