package service

import (
	"github.com/vibeitco/go-utils/config"
	"github.com/vibeitco/go-utils/storage"
)

type Config struct {
	config.Core
	MongoDB MongoDB `json:"mongoDb" yaml:"mongoDb"`
}

type MongoDB struct {
	Addresses []string `json:"addresses" yaml:"addresses"`
	Database  string   `json:"database" yaml:"database"`
	Username  string   `json:"username" yaml:"username"`
	Password  string   `json:"password" yaml:"password" envconfig:"MONGODB_PASSWORD" required:"true"`
}

func (m MongoDB) ToStorageConfig() storage.Config {
	return storage.Config{
		Addresses: m.Addresses,
		Database:  m.Database,
		Username:  m.Username,
		Password:  m.Password,
		Tls:       true,
	}
}
