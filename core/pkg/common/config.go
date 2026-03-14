package common

type Config struct {
	Server       ServerConfig       `yaml:"server"`
	RateLimiting RateLimitingConfig `yaml:"rate_limiting"`
	Security     SecurityConfig     `yaml:"security"`
}

type SecurityConfig struct {
	AnomalyThreshold int `yaml:"anomaly_threshold"`
}

type ServerConfig struct {
	Listen   string    `yaml:"listen"`
	Upstream string    `yaml:"upstream"`
	TLS      TLSConfig `yaml:"tls"`
}

type TLSConfig struct {
	Enabled   bool   `yaml:"enabled"`
	CertFile  string `yaml:"cert_file"`
	KeyFile   string `yaml:"key_file"`
	ListenTLS string `yaml:"listen_tls"`
}

type RateLimitingConfig struct {
	Enabled bool    `yaml:"enabled"`
	RPS     float64 `yaml:"rps"`
	Burst   int     `yaml:"burst"`
}
