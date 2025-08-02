package db

import (
    "time"
)

type Config struct {
    ID             uint `gorm:"primaryKey"`
    Name           string
    Type           string // "http", "tcp", "dns"
    Address        string
    IntervalSeconds int
    CreatedAt      time.Time
    UpdatedAt      time.Time
    Disabled       bool `gorm:"not null;default:false"`
    
}

type Result struct {
    ID            uint `gorm:"primaryKey"`
    ConfigID      uint
    Status        bool
    ResponseTime  int
    CheckedAt     time.Time
}