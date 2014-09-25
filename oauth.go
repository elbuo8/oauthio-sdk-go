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
	"fmt"
	//"runtime/debug"
)

const (
	OAuthdURL = "https://oauth.io"
	OAuthdBase = "/auth"
	Version   = "0.0.1"
)

type OAuth struct {
	OAuthdURL  string
	OAuthdBase string
	appKey     string
	appSecret  string
	Version    string
	Client     *http.Client
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
		OAuthdBase: OAuthdBase,
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

func (o *OAuth) SetOAuthdBase(url string) {
	o.OAuthdBase = url
}

func (o *OAuth) GenerateStateToken() (string, error) {
	id, err := uuid.NewV4()
	if err != nil {
		return "", err
	}
	return id.String(), nil
}


func (o *OAuth) Redirect(provider string, redirectTo string) func(http.ResponseWriter, *http.Request) {
	fmt.Println("haahaha")
	return func(res http.ResponseWriter, req *http.Request) {
		fmt.Println("uhuhuhuh")
		csrf_token, _ := o.GenerateStateToken()
		Url, _ := url.Parse(o.OAuthdURL + o.OAuthdBase + "/" + provider)
		parameters := url.Values{}
		parameters.Add("k", o.appKey)
		parameters.Add("opts", "{\"state\":\"" + csrf_token + "\"}")
		parameters.Add("redirect_type", "server")
		parameters.Add("redirect_uri", redirectTo)
		Url.RawQuery = parameters.Encode()
		fmt.Println(Url.String(), http.StatusFound, csrf_token)
		http.Redirect(res, req, Url.String(), http.StatusFound)
	}
}

/*func Redirect(provider string, redirectTo string) func(http.ResponseWriter, *http.Request) {
	return func(res http.ResponseWriter, req *http.Request) {
		Url, _ := url.Parse("https://oauth.io/auth/" + provider)
		parameters := url.Values{}
		parameters.Add("k", "fsjhfkdsjhfkdjs")
		parameters.Add("opts", "{\"state\":\"fdsafdsa\"}")
		parameters.Add("redirect_type", "server")
		parameters.Add("redirect_uri", redirectTo)
		Url.RawQuery = parameters.Encode()
		fmt.Println(Url.String(), http.StatusFound)
		http.Redirect(res, req, Url.String(), http.StatusFound)
	}
}*/


type OAuthResponseData struct {
	code 		string
}

type OAuthResponse struct {
	status 		string
	data 		OAuthResponseData
	state 		string
	provider 	string
}

func (o *OAuth) Callback(cb func(*OAuthRequestObject, error, http.ResponseWriter, *http.Request)) func(http.ResponseWriter, *http.Request) {
	return func(rw http.ResponseWriter, req *http.Request) {
		oauthioParam := req.URL.Query()["oauthio"][0]
		var r OAuthResponse
		json.Unmarshal([]byte(oauthioParam), r)
		fmt.Println(oauthioParam)
		oauthResp, error := o.Auth(r.data.code)
		cb(oauthResp, error, rw, req)
		//cb(nil, nil, nil, nil)
	}
}

func (o *OAuth) Auth(code string) (*OAuthRequestObject, error) {
	data := url.Values{}
	data.Set("code", code)
	data.Set("key", o.appKey)
	data.Set("secret", o.appSecret)
	fmt.Println(code, o.appKey, o.appSecret)
	response, err := http.PostForm(o.OAuthdURL+"/auth/access_token", data)
	if err != nil {
		return nil, errors.New("oauth.go: Couldn't communicate with Oauthd")
	}
	body, err := ioutil.ReadAll(response.Body)
	response.Body.Close()
	if err != nil {
		return nil, errors.New("oauth.go: Couldn't read Oauthd response")
	}
	fmt.Println(string(body))

	var oauthRespPayload OAuthRequestObjectPayload
	err = json.Unmarshal(body, oauthRespPayload)
	fmt.Println(string(body))
	oauthResp := &oauthRespPayload.data
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
	return nil, nil
}

func (o *OAuth) RefreshCredentials(creds *OAuthRequestObject, force bool) error {
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
