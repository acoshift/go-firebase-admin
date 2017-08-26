package firebase_test

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/acoshift/go-firebase-admin"
	"github.com/stretchr/testify/assert"
)

func TestSendToDevices(t *testing.T) {

	t.Run("send=success", func(t *testing.T) {
		srv := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
			rw.WriteHeader(http.StatusOK)
			rw.Header().Set("Content-Type", "application/json")
			fmt.Fprint(rw, `{"multicast_id": 5438046884136077786,"success": 1,"failure": 0,"results": [{"message_id":"118218","error": ""}]}`)
		}))
		defer srv.Close()

		app := initApp(t)
		firFCM := app.FCM()
		firFCM.NewFcmSendEndpoint(srv.URL)

		assert.NotNil(t, app)
		assert.NotNil(t, firFCM)

		response, err := firFCM.SendToDevice(context.Background(), "mydevicetoken",
			firebase.Message{Notification: firebase.Notification{
				Title: "Hello go firebase admin",
				Body:  "My little Big Notification",
				Color: "#ffcc33"},
			})

		assert.Nil(t, err)
		assert.NotNil(t, response)
		assert.Equal(t, 1, response.Success)
		assert.Equal(t, 0, response.Failure)
		assert.Equal(t, 0, response.CanonicalIDs)
		assert.Equal(t, int64(5438046884136077786), response.MulticastID)
		assert.Equal(t, []firebase.Result{firebase.Result{MessageID: "118218"}}, response.Results)
	})

	t.Run("send=BadToken", func(t *testing.T) {
		srv := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
			rw.WriteHeader(http.StatusOK)
			rw.Header().Set("Content-Type", "application/json")
			fmt.Fprint(rw, `{"multicast_id": 5438046884136077786,"success": 0,"failure": 1,"results": [{"error": "InvalidRegistration"}]}`)
		}))
		defer srv.Close()

		app := initApp(t)
		firFCM := app.FCM()
		firFCM.NewFcmSendEndpoint(srv.URL)

		assert.NotNil(t, app)
		assert.NotNil(t, firFCM)

		response, err := firFCM.SendToDevice(context.Background(), "mydevicetoken",
			firebase.Message{Notification: firebase.Notification{
				Title: "Hello go firebase admin",
				Body:  "My little Big Notification",
				Color: "#ffcc33"},
			})

		assert.Nil(t, err)
		assert.NotNil(t, response)
		assert.Equal(t, 0, response.Success)
		assert.Equal(t, 1, response.Failure)
		assert.Equal(t, 0, response.CanonicalIDs)
		assert.Equal(t, int64(5438046884136077786), response.MulticastID)
		assert.Equal(t, []firebase.Result{firebase.Result{Error: firebase.ErrInvalidRegistration}}, response.Results)
	})

	t.Run("send=missingDestination", func(t *testing.T) {
		srv := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
			rw.WriteHeader(http.StatusBadRequest)
			rw.Header().Set("Content-Type", "application/json")
			fmt.Fprint(rw, `to`)
		}))
		defer srv.Close()

		app := initApp(t)
		firFCM := app.FCM()
		firFCM.NewFcmSendEndpoint(srv.URL)

		assert.NotNil(t, app)
		assert.NotNil(t, firFCM)

		response, err := firFCM.SendToDevice(context.Background(), "mydevicetoken",
			firebase.Message{Notification: firebase.Notification{
				Title: "Hello go firebase admin",
				Body:  "My little Big Notification",
				Color: "#ffcc33"},
			})

		assert.NotNil(t, err)
		assert.Nil(t, response)
		assert.Equal(t, fmt.Errorf("StatusCode=%d, Desc=%s", http.StatusBadRequest, firebase.ErrInvalidParameters.Error()), err)
	})

	t.Run("send=BadRequest", func(t *testing.T) {
		srv := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
			rw.WriteHeader(http.StatusBadRequest)
		}))
		defer srv.Close()

		app := initApp(t)
		firFCM := app.FCM()
		firFCM.NewFcmSendEndpoint(srv.URL)

		assert.NotNil(t, app)
		assert.NotNil(t, firFCM)

		response, err := firFCM.SendToDevice(context.Background(), "mydevicetoken",
			firebase.Message{Notification: firebase.Notification{
				Title: "Hello go firebase admin",
				Body:  "My little Big Notification",
				Color: "#ffcc33"},
			})

		assert.NotNil(t, err)
		assert.Nil(t, response)
		assert.Equal(t, err, fmt.Errorf("StatusCode=%d, Desc=%s", http.StatusBadRequest, firebase.ErrInvalidParameters.Error()))
	})

	t.Run("send=Unauthorized", func(t *testing.T) {
		srv := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
			rw.WriteHeader(http.StatusUnauthorized)
		}))
		defer srv.Close()

		app := initApp(t)
		firFCM := app.FCM()
		firFCM.NewFcmSendEndpoint(srv.URL)

		assert.NotNil(t, app)
		assert.NotNil(t, firFCM)

		response, err := firFCM.SendToDevice(context.Background(), "mydevicetoken",
			firebase.Message{Notification: firebase.Notification{
				Title: "Hello go firebase admin",
				Body:  "My little Big Notification",
				Color: "#ffcc33"},
			})

		assert.NotNil(t, err)
		assert.Nil(t, response)
		assert.Equal(t, err, fmt.Errorf("StatusCode=%d, Desc=%s", http.StatusUnauthorized, firebase.ErrAuthentication.Error()))
	})
}

func TestSendToTopic(t *testing.T) {

	t.Run("send=success", func(t *testing.T) {
		srv := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
			rw.WriteHeader(http.StatusOK)
			rw.Header().Set("Content-Type", "application/json")
			fmt.Fprint(rw, `{"multicast_id": 5438046884136077786,"success": 1,"failure": 0,"results": [{"message_id":"118218","error": ""}]}`)
		}))
		defer srv.Close()

		app := initApp(t)
		firFCM := app.FCM()
		firFCM.NewFcmSendEndpoint(srv.URL)

		assert.NotNil(t, app)
		assert.NotNil(t, firFCM)

		response, err := firFCM.SendToTopic(context.Background(), "mytopic",
			firebase.Message{Notification: firebase.Notification{
				Title: "Hello go firebase admin",
				Body:  "My little Big Notification",
				Color: "#ffcc33"},
			})

		assert.Nil(t, err)
		assert.NotNil(t, response)
		assert.Equal(t, 1, response.Success)
		assert.Equal(t, 0, response.Failure)
		assert.Equal(t, 0, response.CanonicalIDs)
		assert.Equal(t, int64(5438046884136077786), response.MulticastID)
		assert.Equal(t, []firebase.Result{firebase.Result{MessageID: "118218"}}, response.Results)
	})

	t.Run("send=BadToken", func(t *testing.T) {
		srv := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
			rw.WriteHeader(http.StatusOK)
			rw.Header().Set("Content-Type", "application/json")
			fmt.Fprint(rw, `{"multicast_id": 5438046884136077786,"success": 0,"failure": 1,"results": [{"error": "InvalidRegistration"}]}`)
		}))
		defer srv.Close()

		app := initApp(t)
		firFCM := app.FCM()
		firFCM.NewFcmSendEndpoint(srv.URL)

		assert.NotNil(t, app)
		assert.NotNil(t, firFCM)

		response, err := firFCM.SendToTopic(context.Background(), "mytopic",
			firebase.Message{Notification: firebase.Notification{
				Title: "Hello go firebase admin",
				Body:  "My little Big Notification",
				Color: "#ffcc33"},
			})

		assert.Nil(t, err)
		assert.NotNil(t, response)
		assert.Equal(t, 0, response.Success)
		assert.Equal(t, 1, response.Failure)
		assert.Equal(t, 0, response.CanonicalIDs)
		assert.Equal(t, int64(5438046884136077786), response.MulticastID)
		assert.Equal(t, []firebase.Result{firebase.Result{Error: firebase.ErrInvalidRegistration}}, response.Results)
	})

	t.Run("send=missingDestination", func(t *testing.T) {
		srv := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
			rw.WriteHeader(http.StatusBadRequest)
			rw.Header().Set("Content-Type", "application/json")
			fmt.Fprint(rw, `to`)
		}))
		defer srv.Close()

		app := initApp(t)
		firFCM := app.FCM()
		firFCM.NewFcmSendEndpoint(srv.URL)

		assert.NotNil(t, app)
		assert.NotNil(t, firFCM)

		response, err := firFCM.SendToTopic(context.Background(), "mytopic",
			firebase.Message{Notification: firebase.Notification{
				Title: "Hello go firebase admin",
				Body:  "My little Big Notification",
				Color: "#ffcc33"},
			})

		assert.NotNil(t, err)
		assert.Nil(t, response)
		assert.Equal(t, fmt.Errorf("StatusCode=%d, Desc=%s", http.StatusBadRequest, firebase.ErrInvalidParameters.Error()), err)
	})

	t.Run("send=BadRequest", func(t *testing.T) {
		srv := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
			rw.WriteHeader(http.StatusBadRequest)
		}))
		defer srv.Close()

		app := initApp(t)
		firFCM := app.FCM()
		firFCM.NewFcmSendEndpoint(srv.URL)

		assert.NotNil(t, app)
		assert.NotNil(t, firFCM)

		response, err := firFCM.SendToTopic(context.Background(), "mytopic",
			firebase.Message{Notification: firebase.Notification{
				Title: "Hello go firebase admin",
				Body:  "My little Big Notification",
				Color: "#ffcc33"},
			})

		assert.NotNil(t, err)
		assert.Nil(t, response)
		assert.Equal(t, err, fmt.Errorf("StatusCode=%d, Desc=%s", http.StatusBadRequest, firebase.ErrInvalidParameters.Error()))
	})

	t.Run("send=Unauthorized", func(t *testing.T) {
		srv := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
			rw.WriteHeader(http.StatusUnauthorized)
		}))
		defer srv.Close()

		app := initApp(t)
		firFCM := app.FCM()
		firFCM.NewFcmSendEndpoint(srv.URL)

		assert.NotNil(t, app)
		assert.NotNil(t, firFCM)

		response, err := firFCM.SendToTopic(context.Background(), "mytopic",
			firebase.Message{Notification: firebase.Notification{
				Title: "Hello go firebase admin",
				Body:  "My little Big Notification",
				Color: "#ffcc33"},
			})

		assert.NotNil(t, err)
		assert.Nil(t, response)
		assert.Equal(t, err, fmt.Errorf("StatusCode=%d, Desc=%s", http.StatusUnauthorized, firebase.ErrAuthentication.Error()))
	})

}
