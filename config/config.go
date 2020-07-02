package config

import (
	"github.com/knadh/koanf"
	"github.com/knadh/koanf/providers/confmap"
	"github.com/knadh/koanf/providers/env"
)

var k = koanf.New(".")

// Config contains configuration parameters for this program
type Config struct {
	HubAPIURL                      string
	ConnectTimeout, RequestTimeout int
}

// NewConfig reads configuration from environment variables
func NewConfig() *Config {

	k.Load(confmap.Provider(map[string]interface{}{
		"HUB_API_URL":         "http://localhost:8001/rpc/api",
		"HUB_CONNECT_TIMEOUT": 10,
		"HUB_REQUEST_TIMEOUT": 10,
	}, "."), nil)

	k.Load(env.Provider("HUB_", ".", nil), nil)

	return &Config{
		HubAPIURL:      k.String("HUB_API_URL"),
		ConnectTimeout: k.Int("HUB_CONNECT_TIMEOUT"),
		RequestTimeout: k.Int("HUB_REQUEST_TIMEOUT"),
	}
}
