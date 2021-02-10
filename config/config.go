package config

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/jschneider98/jgocache/autocert/cache"
	"golang.org/x/crypto/acme/autocert"
	"io/ioutil"
	"net/url"
	"os"
	"strings"
)

// Config file definition
type Config struct {
	Server            ServerOptions           `json:"server"`
	DbConns           []DbConnOptions         `json:"dbConns"`
	GoogleOauth2Creds GoogleOauth2Credentials `json:"googleOauth2Credentials"`
	Integration       IntegrationOptions      `json:"integration"`
	Autocert          AutocertOptions         `json:"autocert"`
	CustomRaw         []string                `json:"custom"`
	Custom            *url.Values             `json:"-"`
	AutocertCache     autocert.Cache          `json:"-"`
}

// Server configuratoin
type ServerOptions struct {
	Mode        string `json:"mode"`
	SessionName string `json:"sessionName"`
	SessionKey  string `json:"sessionKey"`
	EnableSsl   bool   `json:"enableSsl"`
	HttpsHost   string `json:"httpsHost"`
	HttpHost    string `json:"httpHost"`
	HealthHost  string `json:"healthHost"`
}

// DB Connection Strings
type DbConnOptions struct {
	ShardName        string `json:"shardName"`
	Dsn              string `json:"dsn"`
	MaxOpenConns     int    `json:"maxOpenConns"`
	MaxIdleConns     int    `json:"maxIdleConns"`
	ConnMaxLifetime  int    `json:"connMaxLifetime"`
	StatementTimeout int    `json:"statementTimeout"`
}

// Google Oauth2 Credentials
type GoogleOauth2Credentials struct {
	ClientID     string `json:"clientId"`
	ClientSecret string `json:"clientSecret"`
}

// Autocert configuration
type AutocertOptions struct {
	AllowedHost  string            `json:"allowedHost"`
	Email        string            `json:"email"`
	DirectoryURL string            `json:"directoryURL"`
	CacheOptions map[string]string `json:"cacheOptions"`
}

// Integration test configuration
type IntegrationOptions struct {
	ShardName string `json:"shardName"`
	AccountId string `json:"accountId"`
	UserEmail string `json:"userEmail"`
}

// Reads json configuration file and returns Config
func New(path string, envVar string) (*Config, error) {
	config, _ := NewFromEnv(envVar)

	if config != nil {
		return config, nil
	}

	return NewFromFile(path)
}

//
func NewFromFile(path string) (*Config, error) {
	file, err := read(path)

	if err != nil {
		return nil, err
	}

	config, err := parse(file)

	if err != nil {
		return nil, err
	}

	config.EnsureBasicOptions()

	err = config.LoadCustomOptions()

	if err != nil {
		return nil, err
	}

	_, err = config.GetAutocertCache()

	if err != nil {
		return nil, err
	}

	return config, nil
}

//
func NewFromEnv(envVar string) (*Config, error) {
	var err error

	conf := os.Getenv(envVar)

	if conf == "" {
		err = errors.New(fmt.Sprintf("Missing '%s' environment varriable.", envVar))

		return nil, err
	}

	config, err := parse([]byte(conf))

	if err != nil {
		return nil, err
	}

	config.EnsureBasicOptions()

	err = config.LoadCustomOptions()

	if err != nil {
		return nil, err
	}

	_, err = config.GetAutocertCache()

	if err != nil {
		return nil, err
	}

	return config, nil
}

// Conditionally load default values
func (c *Config) EnsureBasicOptions() {

	if c.Server.Mode == "" {
		c.Server.Mode = "prod"
	}

	if c.Server.SessionName == "" {
		c.Server.SessionName = "web-session"
	}

	if c.Server.SessionKey == "" {
		c.Server.SessionKey = "u46IpCV9y5Vjsi5YvODJEhgOY8m9JVE4"
	}

	// Default acme URL
	if c.Autocert.DirectoryURL == "" {
		c.Autocert.DirectoryURL = "https://acme-v01.api.letsencrypt.org/directory"
	}
}

//
func (c *Config) LoadCustomOptions() error {

	for _, option := range c.CustomRaw {
		optionParts := strings.Split(option, "-:-")

		if len(optionParts) == 2 {
			c.Custom.Set(optionParts[0], optionParts[1])
		}
	}

	return nil
}

//
func (c *Config) GetAutocertCache() (autocert.Cache, error) {
	var err error

	if c.AutocertCache != nil {
		return c.AutocertCache, nil
	}

	if c.Server.EnableSsl == false {
		return nil, nil
	}

	c.AutocertCache, err = cache.NewCacheFactory(c.Autocert.CacheOptions)

	return c.AutocertCache, err
}

func read(path string) ([]byte, error) {
	file, err := ioutil.ReadFile(path)

	if err != nil {
		return nil, err
	}

	return file, nil
}

func parse(jsonConfig []byte) (*Config, error) {
	var config Config

	err := json.Unmarshal(jsonConfig, &config)

	if err != nil {
		return nil, err
	}

	return &config, nil
}
