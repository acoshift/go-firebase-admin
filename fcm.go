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
	fcmSendEndpoint        = "https://fcm.googleapis.com/fcm/send"
	fcmTopicAddEndpoint    = "https://iid.googleapis.com/iid/v1:batchAdd"
	fcmTopicRemoveEndpoint = "https://iid.googleapis.com/iid/v1:batchRemove"
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

// SubscribeDeviceToTopic subscribe to a device to a topic by providing a registration token for the device to subscribe
// see https://firebase.google.com/docs/cloud-messaging/admin/manage-topic-subscriptions#subscribe_to_a_topic
func (fcm *FCM) SubscribeDeviceToTopic(registrationToken string, topic string) (*Response, error) {
	return fcm.sendFirebaseTopicRequest(fcmTopicAddEndpoint, Topic{To: topic, RegistrationTokens: []string{registrationToken}})
}

// SubscribeDevicesToTopic subscribe devices to a topic by providing a registrationtokens for the devices to subscribe
// see https://firebase.google.com/docs/cloud-messaging/admin/manage-topic-subscriptions#subscribe_to_a_topic
func (fcm *FCM) SubscribeDevicesToTopic(registrationTokens []string, topic string) (*Response, error) {
	return fcm.sendFirebaseTopicRequest(fcmTopicAddEndpoint, Topic{To: topic, RegistrationTokens: registrationTokens})
}

// UnSubscribeDeviceFromTopic Unsubscribe a device to a topic by providing a registration token for the device to unsubscribe
// see https://firebase.google.com/docs/cloud-messaging/admin/manage-topic-subscriptions#unsubscribe_from_a_topic
func (fcm *FCM) UnSubscribeDeviceFromTopic(registrationToken string, topic string) (*Response, error) {
	return fcm.sendFirebaseTopicRequest(fcmTopicRemoveEndpoint, Topic{To: topic, RegistrationTokens: []string{registrationToken}})
}

// UnSubscribeDevicesFromTopic Unsubscribe devices to a topic by providing a registrationtokens for the devices to unsubscribe
// see https://firebase.google.com/docs/cloud-messaging/admin/manage-topic-subscriptions#unsubscribe_from_a_topic
func (fcm *FCM) UnSubscribeDevicesFromTopic(registrationTokens []string, topic string) (*Response, error) {
	return fcm.sendFirebaseTopicRequest(fcmTopicRemoveEndpoint, Topic{To: topic, RegistrationTokens: registrationTokens})
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
	req, err := http.NewRequest("POST", fcmSendEndpoint, bytes.NewBuffer(data))
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

func (fcm *FCM) sendFirebaseTopicRequest(endpoint string, payload Topic) (*Response, error) {

	// validate Topic
	if err := payload.Validate(); err != nil {
		return nil, err
	}

	// marshal Topic
	data, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	// create request
	req, err := http.NewRequest("POST", endpoint, bytes.NewBuffer(data))
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
