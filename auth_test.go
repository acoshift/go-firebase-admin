package admin_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateCustomToken(t *testing.T) {

	app := initApp()
	firAuth := app.Auth()

	assert.NotNil(t, app)
	assert.NotNil(t, firAuth)

	// my claims
	myClaims := make(map[string]string)
	myClaims["name"] = "go-firebase-admin"

	token, err := firAuth.CreateCustomToken("go-firebase-admin", myClaims)

	assert.Nil(t, err)
	assert.NotNil(t, token)
}

// func TestVerifyIDToken(t *testing.T) {
// 	app := initApp()
// 	firAuth := app.Auth()

// 	assert.NotNil(t, app)
// 	assert.NotNil(t, firAuth)

// 	// my claims
// 	myClaims := make(map[string]string)
// 	myClaims["name"] = "go-firebase-admin"
// 	myClaims["kid"] = "polo"

// 	token, err := firAuth.CreateCustomToken("uid", myClaims)

// 	fmt.Printf("%s", token)

// 	assert.Nil(t, err)
// 	assert.NotNil(t, token)

// 	claims, err := firAuth.VerifyIDToken(token)

// 	assert.Nil(t, err)
// 	assert.NotNil(t, claims)
// }
