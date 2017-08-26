package firebase_test

import (
	"context"
	"testing"

	"github.com/acoshift/go-firebase-admin"
	"github.com/stretchr/testify/assert"
)

func TestCreateCustomToken(t *testing.T) {
	app := initApp(t)
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

func TestUser(t *testing.T) {
	ctx := context.Background()
	app := initApp(t)
	auth := app.Auth()

	createUser := &firebase.User{
		DisplayName:   "Tester",
		Email:         "test@test.com",
		EmailVerified: true,
		Password:      "test123",
	}

	var userID string

	t.Run("Create", func(t *testing.T) {
		uID, err := auth.CreateUser(ctx, createUser)
		assert.NoError(t, err)
		assert.NotNil(t, uID)
		userID = uID
	})

	t.Run("Get", func(t *testing.T) {
		u, err := auth.GetUser(ctx, userID)
		assert.NoError(t, err)
		assert.NotNil(t, u)
		assert.Equal(t, createUser.DisplayName, u.DisplayName)
		assert.Equal(t, createUser.Email, u.Email)
		assert.Equal(t, createUser.EmailVerified, u.EmailVerified)
		assert.NotEmpty(t, u.Metadata.CreatedAt)
	})

	t.Run("Delete", func(t *testing.T) {
		err := auth.DeleteUser(ctx, userID)
		assert.NoError(t, err)

		u, err := auth.GetUser(ctx, userID)
		assert.Error(t, err)
		assert.Nil(t, u)
	})
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
