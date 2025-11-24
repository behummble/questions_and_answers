package config

import (

)

type Config struct {
	Server ServerConfig
	Log LogConfig
	Storage StorageConfig
}

type ServerConfig struct {
	Host string
	Port int
	ReadTimeout int
	
}

type LogConfig struct {
	Path string
	Level int
}

type StorageConfig struct {
	Host string
	Port int
	DBName string
	Username string
	Password string
	TimeZone string
}

func NewConfig() *Config {
	return &Config{}
}