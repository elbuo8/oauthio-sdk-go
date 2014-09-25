package oauthio

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

type OAuthRequestObject struct {
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
	Client       *http.Client
}

type OAuthRequestObjectPayload struct {
	data		OAuthRequestObject
	status		string
}

func (a *OAuthRequestObject) Get(endpoint string) ([]byte, error) {
	return a.makeRequest("GET", endpoint, nil)
}

func (a *OAuthRequestObject) Post(endpoint string, body interface{}) ([]byte, error) {
	payload, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}
	return a.makeRequest("POST", endpoint, bytes.NewReader(payload))
}

func (a *OAuthRequestObject) Put(endpoint string, body interface{}) ([]byte, error) {
	payload, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}
	return a.makeRequest("PUT", endpoint, bytes.NewReader(payload))
}

func (a *OAuthRequestObject) Del(endpoint string) ([]byte, error) {
	return a.makeRequest("DELETE", endpoint, nil)
}

func (a *OAuthRequestObject) Patch(endpoint string, body interface{}) ([]byte, error) {
	payload, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}
	return a.makeRequest("PATCH", endpoint, bytes.NewReader(payload))
}

func (a *OAuthRequestObject) Me(filters []string) ([]byte, error) {
	qs := url.Values{}
	qs.Set("filters", strings.Join(filters, ","))
	return a.makeMeRequest(qs)
}

func (a *OAuthRequestObject) buildHeaders(r *http.Request) {
	headers := url.Values{}
	headers.Set("k", a.appKey)
	if a.AccessToken == "" {
		headers.Set("access_token", a.AccessToken)
	} else {
		headers.Set("oauth_token", a.OAuthToken)
		headers.Set("oauth_token_secret", a.OAuthSecret)
		headers.Set("oauthv", "1")
	}
	r.Header = http.Header{
		"oauthio": []string{headers.Encode()},
	}
}

func (a *OAuthRequestObject) makeMeRequest(filters url.Values) ([]byte, error) {
	req, _ := http.NewRequest("GET", a.OAuthdURL+"/auth/"+a.Provider+"/me?"+filters.Encode(), nil)
	a.buildHeaders(req)
	return excuteRequest(a.Client, req)
}

func (a *OAuthRequestObject) makeRequest(method, endpoint string, body io.Reader) ([]byte, error) {
	req, _ := http.NewRequest(method, a.OAuthdURL+"/request/"+a.Provider+endpoint, body)
	a.buildHeaders(req)
	return excuteRequest(a.Client, req)
}

func excuteRequest(c *http.Client, r *http.Request) ([]byte, error) {
	response, err := c.Do(r)
	if err != nil {
		return nil, errors.New("oauth_request.go: Couldn't reach Oauthd")
	}
	respBody, err := ioutil.ReadAll(response.Body)
	response.Body.Close()
	if err != nil {
		return nil, errors.New("oauth_request.go: Couldn't read Oauthd response")
	}
	return respBody, nil
}
