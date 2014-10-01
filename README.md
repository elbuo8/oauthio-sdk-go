# OAuth.io SDK for Go

This package simplify the interaction with [OAuth.io](https://oauth.io) or [oauthd](https://github.com/oauth-io/oauthd) using Golang.

Special thanks to [Yamil Asusta](https://github.com/elbuo8) from Sendgrid who is the original creator of this SDK.

### Installation

```bash
go get github.com/oauth-io/sdk-go
```

### Examples

This library was built to be unobtrusive. Meaning that handling codes is up to the user. Examples will be provided on how to use the library on common frameworks.

#### Martini

```go
package main

import (
  "fmt"
  "github.com/oauth-io/sdk-go"
  "github.com/go-martini/martini"
)

func main() {
  m := martini.Classic()
  oauth := oauthio.New("PUBLIC_KEY", "SECRET_KEY")

  //redirect the user to the Facebook authorization page then redirect him to /oauth/redirect
  m.Get("/signin", oauth.Redirect("facebook", "http://localhost:3000/oauth/redirect"))

  //Once redirected, handle the callback and get back a oauthio.OAuthRequestObject object
  m.Get("/oauth/redirect", oauth.Callback(func(res *oauthio.OAuthRequestObject, err error, rw http.ResponseWriter, req *http.Request) {
    if err != nil {
      fmt.Println(err)
      return
    }

    r, _ := res.Me([]string{})

    //fmt.Println(res.AccessToken, err)
    fmt.Println(string(r))
  }))

  m.Run()
}
```

#### Gorilla

```go
package main

import (
  "fmt"
  "github.com/oauth-io/sdk-go"
  "github.com/gorilla/mux"
  "net/http"
)

func main() {
  m := mux.NewRouter()
  oauth := oauthio.New("PUBLIC_KEY", "SECRET_KEY")

  m.HandleFunc("/signin", oauth.Redirect("facebook", "http://localhost:3000/oauth/redirect"))

  m.HandleFunc("/oauth/redirect", oauth.Callback(func(res *oauthio.OAuthRequestObject, err error, rw http.ResponseWriter, req *http.Request) {
    if err != nil {
      fmt.Println(err)
      return
    }

    r, _ := res.Me([]string{})

    //fmt.Println(res.AccessToken, err)
    fmt.Println(string(r))
  }))

  http.Handle("/", m)
}
```

It is completely up to the user to decide where and how to store the tokens.

## MIT

Enjoy, feel free to submit pull requests!