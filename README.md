# go-firebase-admin

> Work in Progress

Firebase Admin SDK for Golang

## Usage

```go
package main

import (
  "log"

  "github.com/acoshift/go-firebase-admin"
)

func main() {
  // Init App
  firApp, err := admin.InitializeApp(admin.ProjectID("YOUR_PROJECT_ID"))
  if err != nil {
    panic(err)
  }
  firAuth := firApp.Auth()

  idToken := "ID_TOKEN_FROM_CLIENT"
  claims, err := firAuth.VerifyIDToken(idToken)
  if err != nil {
    panic(err)
  }

  log.Println(claims)
}
```

## Available functions

- VerifyIDToken


