package config

import (
	"io/ioutil"
	"encoding/json"
	"golang.org/x/crypto/acme/autocert"
	"github.com/jschneider98/jgocache/autocert/cache"
)

// Config file definition
type Config struct {
	Server ServerOptions `json:"server"`
	DbConns []DbConnOptions `json:"dbConns"`
	AutocertCache autocert.Cache `json:"-"`
	Autocert AutocertOptions `json:"autocert"`
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

// Autocert configuration
type AutocertOptions struct {
	AllowedHost string `json:"allowedHost"`
	Email string `json:"email"`
	DirectoryURL string `json:"directoryURL"`
	CacheOptions map[string]string `json:"cacheOptions"`
}

// Reads json configuration file and returns Config
func New(path string) (*Config, error) {
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