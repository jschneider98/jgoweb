package jgoweb

import(
	"io/ioutil"
	"encoding/json"
	"golang.org/x/crypto/acme/autocert"
	"github.com/jschneider98/jgocache/autocert/cache"
)

//
var NewAutocertCache = func() (autocert.Cache, error) {
	options := make(map[string]string)

	file, err := ioutil.ReadFile("./autocert.json")
	
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(file, &options)

	if err != nil {
		return nil, err
	}

	return cache.NewCacheFactory(options)
}
