# go-firebase-admin

[![GoDoc](https://godoc.org/github.com/acoshift/go-firebase-admin?status.svg)](https://godoc.org/github.com/acoshift/go-firebase-admin)

> Work in Progress

Firebase Admin SDK for Golang

## Usage

### Init

```go
package main

import (
  "io/ioutil"
  "log"

  "github.com/acoshift/go-firebase-admin"
)

func main() {
  // Init App
  serviceAccount, _ := ioutil.ReadFile("service_account.json")
  firApp, err := admin.InitializeApp(admin.ProjectID("YOUR_PROJECT_ID"), admin.ServiceAccount(serviceAccount))
  if err != nil {
    panic(err)
  }
  firAuth := firApp.Auth()

  // ...
}
```

### CreateCustomToken

```go
userID := "12345678"
claims := map[string]interface{}{"isAdmin": true}
token, err := firAuth.CreateCustomToken(userID, claims)
```

### VerifyIDToken

```go
idToken := "ID_TOKEN_FROM_CLIENT"
claims, err := firAuth.VerifyIDToken(idToken)
if err != nil {
  panic(err)
}

userID := claims.UserID
log.Println(userID)
```

### GetAccountInfoByUID

```go
user, err := firAuth.GetAccountInfoByUID("123312121")
```

### GetAccountInfoByUIDs

```go
users, err := firAuth.GetAccountInfoByUIDs([]string{"123312121", "2433232", "12121211"})
```

### GetAccountInfoByEmail

```go
user, err := firAuth.GetAccountInfoByEmail("abc@gmail.com")
```

### GetAccountInfoByEmails

```go
users, err := firAuth.GetAccountInfoByEmails([]string{"abc@gmail.com", "qqq@hotmail.com", "aaaqaq@aaa.com"})
```

### DeleteAccount

```go
err := firAuth.DeleteAccount("USER_ID")
```

### CreateAccount

```go
userID, err := firApp.CreateAccount(&admin.Account{
  Email:         "aaa@bbb.com",
  EmailVerified: true,
  Password:      "12345678",
  DisplayName:   "AAA BBB",
})
```

### ListAccount

```go
cursor := firApp.ListAccount(100)
for {
  users, err := cursor.Next()
  if users == nil || err != nil {
    break
  }
  log.Println(len(users))
}
```

### UpdateAccount

```go
err := firApp.UpdateAccount(&UpdateAccount{
  LocalID: "12121212",
  Email: "new_email@email.com",
  Password: "new_password",
  DisplayName: "new name",
})
