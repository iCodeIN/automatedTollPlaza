package config

import "time"

// Cfg ..
type Cfg struct {
	HTTPLog      bool         `json:"httpLog"`
	ServerConfig ServerConfig `json:"serverConfig"`
	MongoConfig  MongoConfig  `json:"mongoConfig"`
}

// ServerConfig ..
type ServerConfig struct {
	Host         string        `json:"host"`
	Port         string        `json:"port"`
	ReadTimeout  time.Duration `json:"readTimeout"`
	WriteTimeout time.Duration `json:"writeTimeout"`
}

// MongoConfig ..
type MongoConfig struct {
	Host string `json:"host"`
}
