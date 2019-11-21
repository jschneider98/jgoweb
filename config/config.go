package config

import (
	"fmt"
	"os"
	"errors"
	"io/ioutil"
	"encoding/json"
	"golang.org/x/crypto/acme/autocert"
	"github.com/jschneider98/jgocache/autocert/cache"
)

// Config file definition
type Config struct {
	Server ServerOptions `json:"server"`
	DbConns []DbConnOptions `json:"dbConns"`
	GoogleOauth2Creds GoogleOauth2Credentials `json:"googleOauth2Credentials"`
	Autocert AutocertOptions `json:"autocert"`
	AutocertCache autocert.Cache `json:"-"`
}

// Server configuratoin
type ServerOptions struct {
	SessionName string `json:sessionName`
	SessionKey string `json:sessionKey`
	EnableSsl bool `json:enableSsl`
	HttpsHost string `json:httpsHost`
	HttpHost string `json:httpHost`
	HealthHost string `json:healthHost`
}

// DB Connection Strings
type DbConnOptions struct {
	ShardName string `json:"shardName"`
	Dsn string `json:"dsn"`
}

// Google Oauth2 Credentials
type GoogleOauth2Credentials struct {
	ClientID     string `json:"clientId"`
	ClientSecret string `json:"clientSecret"`
}

// Autocert configuration
type AutocertOptions struct {
	AllowedHost string `json:"allowedHost"`
	Email string `json:"email"`
	DirectoryURL string `json:"directoryURL"`
	CacheOptions map[string]string `json:"cacheOptions"`
}

// Reads json configuration file and returns Config
func New(path string) (*Config, error) {
	config, _ := NewFromEnv()

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

	_, err = config.GetAutocertCache()

	if err != nil {
		return nil, err
	}

	return config, nil
}

//
func NewFromEnv() (*Config, error) {
	var err error

	conf := os.Getenv("JGO_CONFIG")

	if conf == "" {
		err = errors.New("Missing JGO_CONFIG environment varriable.")
		fmt.Println("Missing JGO_CONFIG environment varriable.")

		return nil, err
	}

	config, err := parse([]byte(conf))

	if err != nil {
		return nil, err
	}

	config.EnsureBasicOptions()

	_, err = config.GetAutocertCache()

	if err != nil {
		return nil, err
	}
fmt.Println("Loaded config from env")
	return config, nil
}

// Conditionally load default values 
func (c *Config) EnsureBasicOptions() {

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
