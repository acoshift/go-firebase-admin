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

// SendToDevice Send Message to individual device
// see https://firebase.google.com/docs/cloud-messaging/admin/send-messages#send_to_individual_devices
func (fcm *FCM) SendToDevice(registrationToken string, payload Message) (*Response, error) {

	// assign recipient
	payload.To = registrationToken

	// flush other recipients
	payload.RegistrationIDs = nil
	payload.Condition = ""

	// send request to Firebase
	return fcm.sendFirebaseRequest(payload)
}

// SendToDevices Send multicast Message to a list of devices
// see https://firebase.google.com/docs/cloud-messaging/admin/send-messages#send_to_individual_devices
func (fcm *FCM) SendToDevices(registrationTokens []string, payload Message) (*Response, error) {

	// assign recipient
	payload.RegistrationIDs = registrationTokens

	// flush other recipients
	payload.To = ""
	payload.Condition = ""

	// send request to Firebase
	return fcm.sendFirebaseRequest(payload)
}

// SendToDeviceGroup Send Message to a device group
// see https://firebase.google.com/docs/cloud-messaging/admin/send-messages#send_to_a_device_group
func (fcm *FCM) SendToDeviceGroup(notificationKey string, payload Message) (*Response, error) {
	return fcm.SendToDevice(notificationKey, payload)
}

// SendToTopic Send Message to a topic
// see https://firebase.google.com/docs/cloud-messaging/admin/send-messages#send_to_a_topic
func (fcm *FCM) SendToTopic(notificationKey string, payload Message) (*Response, error) {
	return fcm.SendToDevice(fmt.Sprint("/topic/", notificationKey), payload)
}

// SendToCondition Send a message to devices subscribed to the combination of topics
// specified by the provided condition.
// see https://firebase.google.com/docs/cloud-messaging/admin/send-messages#send_to_a_condition
func (fcm *FCM) SendToCondition(condition string, payload Message) (*Response, error) {

	// assign recipient
	payload.Condition = condition

	// flush other recipients
	payload.To = ""
	payload.RegistrationIDs = nil

	// send request to Firebase
	return fcm.sendFirebaseRequest(payload)
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

// UnSubscribeDeviceFromTopic TODO NOT IMPLEMENTED
// see https://firebase.google.com/docs/cloud-messaging/admin/manage-topic-subscriptions#unsubscribe_from_a_topic
func (fcm *FCM) UnSubscribeDeviceFromTopic(registrationToken string, topic string) (*Response, error) {
	return nil, ErrNotImplemented
}

// UnSubscribeDevicesFromTopic TODO NOT IMPLEMENTED
// see https://firebase.google.com/docs/cloud-messaging/admin/manage-topic-subscriptions#unsubscribe_from_a_topic
func (fcm *FCM) UnSubscribeDevicesFromTopic(registrationTokens []string, topic string) (*Response, error) {
	return nil, ErrNotImplemented
}

func (fcm *FCM) sendFirebaseRequest(payload Message) (*Response, error) {

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
