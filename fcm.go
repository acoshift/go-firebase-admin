package admin

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

// FCM type
type FCM struct {
	app *App
}

const (
	fcmEndpoint = "https://fcm.googleapis.com/fcm/send"
)

func newFCM(app *App) *FCM {
	return &FCM{
		app: app,
	}
}

// SendToDevice Send Message to individual devices
// see https://firebase.google.com/docs/cloud-messaging/admin/send-messages#send_to_individual_devices
func (fcm *FCM) SendToDevice(registrationToken string, payload Message) (*Response, error) {

	// assign recipient with registrationToken
	payload.To = registrationToken

	// validate Message
	if err := payload.Validate(); err != nil {
		return nil, err
	}

	// marshal message
	data, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	// create request
	req, err := http.NewRequest("POST", fcmEndpoint, bytes.NewBuffer(data))
	if err != nil {
		return nil, err
	}

	// add headers
	req.Header.Set("Authorization", fmt.Sprintf("key=%s", fcm.app.apiKey))
	req.Header.Set("Content-Type", "application/json")

	// TODO use fcm.app.client if possible (not working actually)
	client := &http.Client{}

	// execute request
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		fmt.Printf("HTTP status: %d\n", resp.StatusCode)
		fmt.Printf("HTTP body: %d\n", resp.Body)

		switch resp.StatusCode {
		case http.StatusBadRequest:
			return nil, fmt.Errorf("StatusCode=%d, Desc=%s", resp.StatusCode, ErrInvalidParameters.Error())

		case http.StatusUnauthorized:
			return nil, fmt.Errorf("StatusCode=%d, Desc=%s", resp.StatusCode, ErrAuthentication.Error())

		case http.StatusInternalServerError:
		default:
			return nil, fmt.Errorf("StatusCode=%d, Desc=%s", resp.StatusCode, ErrInternalServerError.Error())
		}
	}

	// build response
	response := new(Response)
	if err := json.NewDecoder(resp.Body).Decode(response); err != nil {
		return nil, err
	}

	return response, nil
}

// SendToDeviceGroup Send Message to a device group
// see https://firebase.google.com/docs/cloud-messaging/admin/send-messages#send_to_a_device_group
func (fcm *FCM) SendToDeviceGroup(notificationKey string, payload Message) (*Response, error) {
	return nil, ErrNotImplemented
}

// SendToTopic TODO NOT IMPLEMENTED
// see https://firebase.google.com/docs/cloud-messaging/admin/send-messages#send_to_a_topic
func (fcm *FCM) SendToTopic(notificationKey string, payload Message) (*Response, error) {
	return nil, ErrNotImplemented
}

// SendToCondition TODO NOT IMPLEMENTED
// see https://firebase.google.com/docs/cloud-messaging/admin/send-messages#send_to_a_condition
func (fcm *FCM) SendToCondition(condition string, payload Message) (*Response, error) {
	return nil, ErrNotImplemented
}

// SubscribeDeviceToTopic TODO NOT IMPLEMENTED
// see https://firebase.google.com/docs/cloud-messaging/admin/manage-topic-subscriptions#subscribe_to_a_topic
func (fcm *FCM) SubscribeDeviceToTopic(registrationToken string, topic string) (*Response, error) {
	return nil, ErrNotImplemented
}

// SubscribeDevicesToTopic TODO NOT IMPLEMENTED
// see https://firebase.google.com/docs/cloud-messaging/admin/manage-topic-subscriptions#subscribe_to_a_topic
func (fcm *FCM) SubscribeDevicesToTopic(registrationTokens []string, topic string) (*Response, error) {
	return nil, ErrNotImplemented
}
