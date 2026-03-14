package common

type Config struct {
	Server ServerConfig `yaml:"server"`
}

type ServerConfig struct {
	Listen   string `yaml:"listen"`
	Upstream string `yaml:"upstream"`
}
