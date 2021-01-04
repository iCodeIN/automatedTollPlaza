package config

import "time"

// Cfg ..
type Cfg struct {
	HTTPLog      bool         `json:"httpLog"`
	ServerConfig ServerConfig `json:"serverConfig"`
	MongoConfig  MongoConfig  `json:"mongoConfig"`
	Pricing      Pricing      `json:"pricing"`
}

// Pricing ..
type Pricing struct {
	Default     PriceValue            `json:"default"`
	VehicleType map[string]PriceValue `json:"vehicleType"`
}

// PriceValue ..
type PriceValue struct {
	OneWay float64 `json:"oneWay"`
	TwoWay float64 `json:"twoWay"`
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
