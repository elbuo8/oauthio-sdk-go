package oauthio

import (
	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
)

const (
	OAuthdURL = "https://oauth.io"
)

type AuthRes struct {
	AccessToken  string `json:"access_token"`
	OAuthToken   string `json:"oauth_token"`
	OAuthSecret  string `json:"oauth_token_secret"`
	State        string `json:"state"`
	Provider     string `json:"provider"`
	ExpiresIn    int64  `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
	ExpireDate   int64
	Refreshed    bool
	OAuthdURL    string
	appKey       string
	Client       *http.DefaultClient
}

func (a *AuthRes) Get(endpoint string) ([]byte, error) {
	return a.makeRequest("GET", endpoint, nil)
}

func (a *AuthRes) Post(endpoint string, body interface{}) ([]byte, error) {
	payload := json.Marshal(body)
	return a.makeRequest("POST", endpoint, bytes.NewReader(payload))
}

func (a *AuthRes) Put(endpoint string, body interface{}) ([]byte, error) {
	payload := json.Marshal(body)
	return a.makeRequest("PUT", endpoint, bytes.NewReader(payload))
}

func (a *AuthRes) Del(endpoint string, body interface{}) ([]byte, error) {
	payload := json.Marshal(body)
	return a.makeRequest("DELETE", endpoint, bytes.NewReader(payload))
}

func (a *AuthRes) Patch(endpoint string, body interface{}) ([]byte, error) {
	payload := json.Marshal(body)
	return a.makeRequest("PATCH", endpoint, bytes.NewReader(payload))
}

func (a *AuthRes) Me(endpoint string, body interface{}) ([]byte, error) {

}

func (a *AuthRes) makeRequest(method, endpoint string, body io.Reader) ([]byte, error) {
	req, _ := http.NewRequest(method, a.OAuthdURL+"/request/"+a.Provider+endpoint, body)
	headers := url.Values{}
	headers.Set("k", a.appKey)
	if a.AccessToken == "" {
		headers.Set("access_token", a.AccessToken)
	} else {
		headers.Set("oauth_token", o.OAuthToken)
		headers.Set("oauth_token_secret", o.OAuthSecret)
		headers.Set("oauthv1", "1")
	}
	req.Header = http.Header{
		"oauthio": headers.Encode(),
	}
	response, err := a.Client.Do(req)
	if err != nil {
		return nil, "oauth_request.go: Couldn't reach Oauthd"
	}
	respBody, err := ioutil.ReadAll(response.Body)
	response.Body.Close()
	if err != nil {
		return nil, "oauth_request.go: Couldn't read Oauthd response"
	}
	return respBody, nil
}
