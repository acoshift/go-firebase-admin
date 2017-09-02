package firebase_test

import (
	"testing"
	"time"

	"fmt"

	firebase "github.com/acoshift/go-firebase-admin"
	"github.com/stretchr/testify/assert"
)

func TestValidToken(t *testing.T) {
	t.Run("Valid", func(t *testing.T) {
		now := time.Now()
		token := &firebase.Token{
			IssuedAt:  now.Unix(),
			ExpiresAt: now.Unix(),
		}

		err := token.Valid()
		assert.Nil(t, err)
	})

	t.Run("usedbefore", func(t *testing.T) {
		now := time.Now().AddDate(1, 0, 0)
		now2 := time.Now()
		token := &firebase.Token{
			IssuedAt:  now.Unix(),
			ExpiresAt: now2.Unix(),
		}

		err := token.Valid()
		assert.NotNil(t, err)
		assert.Equal(t, fmt.Errorf("token used before issued"), err)
	})

	t.Run("expired", func(t *testing.T) {
		now := time.Now()
		token := &firebase.Token{
			IssuedAt:  now.Unix(),
			ExpiresAt: 1500651130,
		}

		err := token.Valid()
		assert.NotNil(t, err)
		assert.Equal(t, fmt.Errorf("token is expired by %v", time.Unix(now.Unix(), 0).Sub(time.Unix(token.ExpiresAt, 0))), err)
	})
}
