package admin_test

import (
	"testing"
	"time"

	"fmt"

	admin "github.com/acoshift/go-firebase-admin"
	"github.com/stretchr/testify/assert"
)

func TestValidClaims(t *testing.T) {

	t.Run("Valid", func(t *testing.T) {
		now := time.Now()
		claim := &admin.Claims{
			IssuedAt:  now.Unix(),
			ExpiresAt: now.Unix(),
		}

		err := claim.Valid()
		assert.Nil(t, err)
	})

	t.Run("usedbefore", func(t *testing.T) {
		now := time.Now().AddDate(1, 0, 0)
		now2 := time.Now()
		claim := &admin.Claims{
			IssuedAt:  now.Unix(),
			ExpiresAt: now2.Unix(),
		}

		err := claim.Valid()
		assert.NotNil(t, err)
		assert.Equal(t, fmt.Errorf("token used before issued"), err)
	})

	t.Run("expired", func(t *testing.T) {
		now := time.Now()
		claim := &admin.Claims{
			IssuedAt:  now.Unix(),
			ExpiresAt: 1500651130,
		}

		err := claim.Valid()
		assert.NotNil(t, err)
		assert.Equal(t, fmt.Errorf("token is expired by %v", time.Unix(now.Unix(), 0).Sub(time.Unix(claim.ExpiresAt, 0))), err)
	})
}
