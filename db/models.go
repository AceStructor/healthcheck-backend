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
    Timeout		   int
    CreatedAt      time.Time
    LastChecked    time.Time
    Disabled       bool `gorm:"not null;default:false"`
}

type Result struct {
    ID            uint `gorm:"primaryKey"`
    ConfigID      uint
    Status        bool
    Text	      string
    ResponseTime  int
    CheckedAt     time.Time
}

type JoinedResult struct {
	Name	      string
	Type          string
	Address       string
	Status        bool
	Text	      string
}
