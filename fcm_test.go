package admin_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSendToDevices(t *testing.T) {

	app := initApp()
	firFCM := app.FCM()

	assert.NotNil(t, firFCM)
}
