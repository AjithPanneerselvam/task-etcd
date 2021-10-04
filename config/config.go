package config

import (
	"github.com/kelseyhightower/envconfig"
)

// Config represents the environment config values
type Config struct {
	LogLevel string `envconfig:"LOG_LEVEL" required:"true"`

	HostName   string `envconfig:"HOST_NAME" required:"true"`
	ListenPort string `envconfig:"LISTEN_PORT" required:"true"`

	EtcdURLS []string `envconfig:"ETCD_URLS" required:"true"`

	GithubClientID     string `envconfig:"GITHUB_CLIENT_ID" required:"true"`
	GithubClientSecret string `envconfig:"GITHUB_CLIENT_SECRET" required:"true"`
	GithubOAuthURL     string `envconfig:"GITHUB_OAUTH_URL" required:"true"`
	GithubTimeoutInSec int32  `envconfig:"GITHUB_TIMEOUT_IN_SEC" required:"true"`
	GithubAPIURL       string `envconfig:"GITHUB_API_URL" required:"true"`
}

// Load loads the config
func Load() (*Config, error) {
	var config Config
	err := envconfig.Process("", &config)
	return &config, err
}
