package oauthio

import (
	"encoding/json"
	"errors"
	"github.com/nu7hatch/gouuid"
	"io/ioutil"
	"net"
	"net/http"
	"net/url"
	"time"
)

const (
	OAuthdURL = "https://oauth.io"
	Version   = "0.0.1"
)

type OAuth struct {
	OAuthdURL string
	appKey    string
	appSecret string
	Version   string
	Client    *http.Client
}

func timeoutHandler(network, address string) (net.Conn, error) {
	return net.DialTimeout(network, address, time.Duration(5*time.Second))
}

func New(appkey, appsecret string) *OAuth {
	transport := http.Transport{
		Dial: timeoutHandler,
	}
	return &OAuth{
		OAuthdURL: OAuthdURL,
		appKey:    appkey,
		appSecret: appsecret,
		Version:   Version,
		Client: &http.Client{
			Transport: &transport,
		},
	}
}

func (o *OAuth) GetVersion() string {
	return o.Version
}

func (o *OAuth) SetOAuthdURL(url string) {
	o.OAuthdURL = url
}

func (o *OAuth) GenerateStateToken() (string, error) {
	id, err := uuid.NewV4()
	if err != nil {
		return "", err
	}
	return id.String(), nil
}

func (o *OAuth) Auth(code string) (*AuthRes, error) {
	data := url.Values{}
	data.Set("code", code)
	data.Set("key", o.appKey)
	data.Set("secret", o.appSecret)
	response, err := http.PostForm(o.OAuthdURL+"/auth/access_token", data)
	if err != nil {
		return nil, errors.New("oauth.go: Couldn't communicate with Oauthd")
	}
	body, err := ioutil.ReadAll(response.Body)
	response.Body.Close()
	if err != nil {
		return nil, errors.New("oauth.go: Couldn't read Oauthd response")
	}
	oauthResp := &AuthRes{}
	err = json.Unmarshal(body, oauthResp)
	if err != nil {
		return nil, errors.New("oauth.go: Couldn't parse response")
	}
	if oauthResp.State == "" {
		return nil, errors.New("oauth.go: State is missing in response")
	}
	oauthResp.ExpireDate = time.Now().Unix() + oauthResp.ExpiresIn
	oauthResp.OAuthdURL = o.OAuthdURL
	oauthResp.Client = o.Client
	oauthResp.appKey = o.appKey
	return oauthResp, nil
}

func (o *OAuth) RefreshCredentials(creds *AuthRes, force bool) error {
	if force || time.Now().Unix() > creds.ExpireDate {
		data := url.Values{}
		data.Set("token", creds.RefreshToken)
		data.Set("key", o.appKey)
		data.Set("secret", o.appSecret)
		response, err := http.PostForm(o.OAuthdURL+"/auth/refresh_token/", data)
		if err != nil {
			return errors.New("oauth.go: Couldn't communicate with Oauthd")
		}
		body, err := ioutil.ReadAll(response.Body)
		response.Body.Close()
		if err != nil {
			return errors.New("oauth.go: Couldn't read Oauthd response")
		}
		err = json.Unmarshal(body, &creds)
		if err != nil {
			return errors.New("oauth.go: Couldn't parse response")
		}
		creds.Refreshed = true
	}
	return nil
}
