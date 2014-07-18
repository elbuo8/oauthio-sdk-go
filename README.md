# OAuth.io SDK for Go

This is a simple package to simplify the interaction with [OAuth.io](https://oauth.io) or [OAuthD](https://github.com/oauth-io/oauthd) using Golang.

### Installation

```bash
go get github.com/elbuo8/oauthio-sdk-go
```

### Examples

This library was built to be unobtrusive. Meaning that handling codes is up to the user. Examples will be provided on how to use the library on common frameworks.

#### Martini

```go
package main

import (
  "fmt"
  "github.com/elbuo8/oauthio-sdk-go"
  "github.com/go-martini/martini"
)

func main() {
  m := martini.Classic()
  oauth := oauthio.New("PUBLIC_KEY", "SECRET_KEY")

  // Generate State Token
  m.Get("/oauth/state_token", func() string {
    token, _ := o.GenerateStateToken()
    return token
  })

  // Use the Code generated in the Front End to Authenticate
  m.Get("/oauth/:token", func(p martini.Params) string {
    reqObject, _ := o.Auth(p["token"])
    r, _ := reqObject.Me([]string{})
    fmt.Println(string(r))
    return "OK"
  })
  m.Run()
}
```

#### Gorilla

```go
package main

import (
  "fmt"
  "github.com/elbuo8/oauthio-sdk-go"
  "github.com/gorilla/mux"
  "net/http"
)

func main() {
  m := mux.NewRouter()
  oauth := oauthio.New("PUBLIC_KEY", "SECRET_KEY")

  // Generate State Token
  m.HandleFunc("/oauth/state_token", func(w http.RequestWriter, r *http.Request) {
    token, _ := o.GenerateStateToken()
    fmt.Fprintf(w, token)
  }))

  m.HandleFunc("/oauth/:token", func(w http.RequestWriter, r *http.Request) {
    vars := mux.Vars(r)
    token := vars["token"]
    reqObject, _ := o.Auth(p["token"])
    r, _ := reqObject.Me([]string{})
    fmt.Println(string(r))
    fmt.Fprintf(w, "OK")
  })
  http.Handle("/", m)
}
```

It is completely up to the user to decide where and how to store the tokens.

## MIT

Enjoy, feel free to submit pull requests!
