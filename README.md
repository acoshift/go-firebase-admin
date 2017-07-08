go-firebase-admin
==============================

[![Go Report Card](https://goreportcard.com/badge/github.com/acoshift/go-firebase-admin)](https://goreportcard.com/report/github.com/acoshift/go-firebase-admin)
[![GoDoc](https://godoc.org/github.com/acoshift/go-firebase-admin?status.svg)](https://godoc.org/github.com/acoshift/go-firebase-admin)

Firebase Admin SDK for Golang

On Wednesday, May 17, 2017 [Google announced at Google IO][1] : Open sourcing the Firebase SDKs.
But for now, there is no official Admin SDK for Golang, only Java, Node and Python official SDKs.

If you decide to use this still in development unofficial SDK, it will be a lot of breaking change, please use any package manager to fix version. 

This go-firebase-admin SDK supports the following functions :

- Authentication
  * CreateCustomToken : [Generate JSON Web Tokens (JWTs) on your server][3], pass them back to a client device, and then use them to authenticate via the signInWithCustomToken() method.
  * VerifyIDToken : [verify the integrity and authenticity of the ID token][4] and retrieve the uid from it.
- Manage Users
  * GetUser : fetching the profile information of users by their uid
  * GetUsers : fetching list of profile information of users by their uid
  * GetUserByEmail : fetching the profile information of users by their email
  * GetUsersByEmail : fetching list of profile information of users by their email
  * ListUsers : fetching the profile information of users
  * CreateUser : create a new Firebase Authentication user
  * UpdateUser : modifying an existing Firebase user's data.
  * DeleteUser : deleting existing Firebase user by uid
  * SendPasswordResetEmail : send password reset for the given user
  * VerifyPassword : verifies given email and password

- Realtime Database
  * 
  
- Cloud Messaging (FCM)
  * TODO

Installation
------------

Install the package with go:

    go get github.com/acoshift/go-firebase-admin

To-Do List
----------

- [ ] update documentation
- [ ] add examples
- [ ] add FCM

Documentation
-------------

You can find documentation for more details on [godoc.org][2].


Initialize Firebase Admin SDK
-----------------------------

```go
package main

import (
  "io/ioutil"

  "github.com/acoshift/go-firebase-admin"
)

func main() {
  // Init App
  serviceAccount, _ := ioutil.ReadFile("service_account.json")
  firApp, err := admin.InitializeApp(context.Background(), admin.AppOptions{
    ServiceAccount: serviceAccount,
    ProjectID: "YOUR_PROJECT_ID",
    DatabaseURL: "YOUR_DATABASE_URL",
  })
  if err != nil {
    panic(err)
  }
  firAuth := firApp.Auth()
  firDatabase := firApp.Database()
  // ...
}
```

[1]: https://opensource.googleblog.com/2017/05/open-sourcing-firebase-sdks.html
[2]: https://godoc.org/github.com/acoshift/go-firebase-admin
[3]: https://firebase.google.com/docs/auth/admin/create-custom-tokens
[4]: https://firebase.google.com/docs/auth/admin/verify-id-tokens
[5]: https://firebase.google.com/docs/auth/admin/manage-users

