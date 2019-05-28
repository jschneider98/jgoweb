package services

import (
	"net/http"
	"encoding/json"
	"io/ioutil"
	"crypto/rand"
	"encoding/base64"
	"golang.org/x/net/context"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"golang.org/x/oauth2/jwt"
)

// Google Service Account
type GoogleServiceAccount struct {
    Subject string             `json:"subext,omitempty"`
    Type string                `json:"type"`
    ProjectId string           `json:"project_id"`
    PrivateKeyId string        `json:"private_key_id"`
    PrivateKey string          `json:"private_key"`
    ClientEmail string         `json:"client_email"`
    ClientId string            `json:"client_id"`
    AuthUri string             `json:"authUri"`
    TokenUri string            `json:"token_uri"`
    AuthProviderCertUrl string `json:"auth_provider_x509_cert_url"`
    ClientCertUrl string       `json:"client_x509_cert_url"`
    Scopes []string            `json:"scopes,omitempty"`
}

var googleOauth2Config *oauth2.Config

// Credentials which stores google ids.
type Credentials struct {
	ClientID     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
}

// Init google sign in creds and oauth2 conf
// @TODO: Dynamically load creds, redirect url, etc
var GetGoogleOauth2Config = func(redirectUrl string) *oauth2.Config {
	
	if googleOauth2Config != nil {
		return googleOauth2Config
	}

	var cred Credentials

	file, err := ioutil.ReadFile("./creds.json")
	
	if err != nil {
		panic(err)
	}

	json.Unmarshal(file, &cred)

	googleOauth2Config = &oauth2.Config{
		ClientID:     cred.ClientID,
		ClientSecret: cred.ClientSecret,
		RedirectURL:  redirectUrl,
		Scopes: []string{
			"https://www.googleapis.com/auth/userinfo.email", 
		},
		Endpoint: google.Endpoint,
	}

	return googleOauth2Config
}

// Google Service Account Info
var GetGoogleServiceAccount = func(clientName string) (*GoogleServiceAccount, error) {
	// @TEMP: Pull info down from data store
	// Debugf("Getting service account for: %s\n", clientName)
	json_str, err := ioutil.ReadFile("./service_account.json")
	
	//fmt.Println("json_str")
	//fmt.Printf("%s", json_str)
	
	if err != nil {
		return nil, err
	}

	var serviceAccount GoogleServiceAccount

	err = json.Unmarshal(json_str, &serviceAccount)

	if err != nil {
		return nil, err
	}

	return &serviceAccount, nil
}

// Return a JWT configured http client
var GetGoogleClient = func(serviceAccount *GoogleServiceAccount) *http.Client {
	ctx := context.Background()

	config := &jwt.Config{
		Email: serviceAccount.ClientEmail,
		PrivateKey: []byte(serviceAccount.PrivateKey),
		Scopes: []string{
			"https://www.googleapis.com/auth/gmail.readonly",
			"https://www.googleapis.com/auth/admin.directory.user.readonly",
		},
		TokenURL: google.JWTTokenURL,
		Subject: serviceAccount.Subject,
	}

	return config.Client(ctx)
}

// Random token for Google Sign in
var GoogleRandomToken = func() string {
	b := make([]byte, 32)
	rand.Read(b)
	
	return base64.StdEncoding.EncodeToString(b)
}

// Google sign in url
var GetGoogleLoginURL = func(redirectUrl string, state string) string {
	config := GetGoogleOauth2Config(redirectUrl)

	return config.AuthCodeURL(state)
}
