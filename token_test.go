package firebase

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestValidToken(t *testing.T) {
	t.Run("Valid", func(t *testing.T) {
		token := &Token{}

		err := token.Valid()
		assert.NoError(t, err)
	})

	t.Run("verifyIssuedAt", func(t *testing.T) {
		now := time.Now()
		token := &Token{
			IssuedAt: now.Unix(),
		}

		res := token.verifyIssuedAt(now.Unix())
		assert.True(t, res)

		res = token.verifyIssuedAt(now.Unix() + 60 /* leeway */)
		assert.True(t, res)
	})

	t.Run("verifyExpiresAt", func(t *testing.T) {
		now := time.Now()
		token := &Token{
			IssuedAt:  now.Unix(),
			ExpiresAt: 1500651130,
		}

		res := token.verifyExpiresAt(now.Unix())
		assert.False(t, res)
	})
}
