package config

import "time"

// Cfg is the config struct
type Cfg struct {
	HTTPLog      bool         `json:"httpLog"`
	ServerConfig ServerConfig `json:"serverConfig"`
	MongoConfig  MongoConfig  `json:"mongoConfig"`
	Pricing      Pricing      `json:"pricing"`
}

// Pricing is the pricing models
type Pricing struct {
	Default     PriceValue            `json:"default"`
	VehicleType map[string]PriceValue `json:"vehicleType"`
}

// PriceValue is priceValue model having oneWay & twoWay field value
type PriceValue struct {
	OneWay float64 `json:"oneWay"`
	TwoWay float64 `json:"twoWay"`
}

// ServerConfig is the configuration to start the web service
type ServerConfig struct {
	Host         string        `json:"host"`
	Port         string        `json:"port"`
	ReadTimeout  time.Duration `json:"readTimeout"`
	WriteTimeout time.Duration `json:"writeTimeout"`
}

// MongoConfig is the configuration to connect to mongodb
type MongoConfig struct {
	Host string `json:"host"`
}
