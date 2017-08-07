package firebase

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
)

// FCM type
type FCM struct {
	app    *App
	client *http.Client
}

var (
	fcmSendEndpoint        = "https://fcm.googleapis.com/fcm/send"
	fcmTopicAddEndpoint    = "https://iid.googleapis.com/iid/v1:batchAdd"
	fcmTopicRemoveEndpoint = "https://iid.googleapis.com/iid/v1:batchRemove"
)

func newFCM(app *App) *FCM {
	return &FCM{
		app:    app,
		client: &http.Client{},
	}
}

// NewFcmSendEndpoint set fcmSendEndpoint URL
func (fcm *FCM) NewFcmSendEndpoint(endpoint string) {
	fcmSendEndpoint = endpoint
}

// NewFcmTopicAddEndpoint set fcmTopicAddEndpoint URL
func (fcm *FCM) NewFcmTopicAddEndpoint(endpoint string) {
	fcmTopicAddEndpoint = endpoint
}

// NewFcmTopicRemoveEndpoint set fcmTopicRemoveEndpoint URL
func (fcm *FCM) NewFcmTopicRemoveEndpoint(endpoint string) {
	fcmTopicRemoveEndpoint = endpoint
}

// SendToDevice Send Message to individual device
// see https://firebase.google.com/docs/cloud-messaging/admin/send-messages#send_to_individual_devices
func (fcm *FCM) SendToDevice(ctx context.Context, registrationToken string, payload Message) (*Response, error) {

	// assign recipient
	payload.To = registrationToken

	// flush other recipients
	payload.RegistrationIDs = nil
	payload.Condition = ""

	// send request to Firebase
	return fcm.sendFirebaseRequest(ctx, payload)
}

// SendToDevices Send multicast Message to a list of devices
// see https://firebase.google.com/docs/cloud-messaging/admin/send-messages#send_to_individual_devices
func (fcm *FCM) SendToDevices(ctx context.Context, registrationTokens []string, payload Message) (*Response, error) {

	// assign recipient
	payload.RegistrationIDs = registrationTokens

	// flush other recipients
	payload.To = ""
	payload.Condition = ""

	// send request to Firebase
	return fcm.sendFirebaseRequest(ctx, payload)
}

// SendToDeviceGroup Send Message to a device group
// see https://firebase.google.com/docs/cloud-messaging/admin/send-messages#send_to_a_device_group
func (fcm *FCM) SendToDeviceGroup(ctx context.Context, notificationKey string, payload Message) (*Response, error) {
	return fcm.SendToDevice(ctx, notificationKey, payload)
}

// SendToTopic Send Message to a topic
// see https://firebase.google.com/docs/cloud-messaging/admin/send-messages#send_to_a_topic
func (fcm *FCM) SendToTopic(ctx context.Context, notificationKey string, payload Message) (*Response, error) {
	return fcm.SendToDevice(ctx, normalizeTopicName(notificationKey), payload)
}

// SendToCondition Send a message to devices subscribed to the combination of topics
// specified by the provided condition.
// see https://firebase.google.com/docs/cloud-messaging/admin/send-messages#send_to_a_condition
func (fcm *FCM) SendToCondition(ctx context.Context, condition string, payload Message) (*Response, error) {

	// assign recipient
	payload.Condition = condition

	// flush other recipients
	payload.To = ""
	payload.RegistrationIDs = nil

	// send request to Firebase
	return fcm.sendFirebaseRequest(ctx, payload)
}

// SubscribeDeviceToTopic subscribe to a device to a topic by providing a registration token for the device to subscribe
// see https://firebase.google.com/docs/cloud-messaging/admin/manage-topic-subscriptions#subscribe_to_a_topic
func (fcm *FCM) SubscribeDeviceToTopic(ctx context.Context, registrationToken string, topic string) (*Response, error) {
	return fcm.sendFirebaseTopicRequest(ctx, fcmTopicAddEndpoint, Topic{To: normalizeTopicName(topic), RegistrationTokens: []string{registrationToken}})
}

// SubscribeDevicesToTopic subscribe devices to a topic by providing a registrationtokens for the devices to subscribe
// see https://firebase.google.com/docs/cloud-messaging/admin/manage-topic-subscriptions#subscribe_to_a_topic
func (fcm *FCM) SubscribeDevicesToTopic(ctx context.Context, registrationTokens []string, topic string) (*Response, error) {
	return fcm.sendFirebaseTopicRequest(ctx, fcmTopicAddEndpoint, Topic{To: normalizeTopicName(topic), RegistrationTokens: registrationTokens})
}

// UnSubscribeDeviceFromTopic Unsubscribe a device to a topic by providing a registration token for the device to unsubscribe
// see https://firebase.google.com/docs/cloud-messaging/admin/manage-topic-subscriptions#unsubscribe_from_a_topic
func (fcm *FCM) UnSubscribeDeviceFromTopic(ctx context.Context, registrationToken string, topic string) (*Response, error) {
	return fcm.sendFirebaseTopicRequest(ctx, fcmTopicRemoveEndpoint, Topic{To: normalizeTopicName(topic), RegistrationTokens: []string{registrationToken}})
}

// UnSubscribeDevicesFromTopic Unsubscribe devices to a topic by providing a registrationtokens for the devices to unsubscribe
// see https://firebase.google.com/docs/cloud-messaging/admin/manage-topic-subscriptions#unsubscribe_from_a_topic
func (fcm *FCM) UnSubscribeDevicesFromTopic(ctx context.Context, registrationTokens []string, topic string) (*Response, error) {
	return fcm.sendFirebaseTopicRequest(ctx, fcmTopicRemoveEndpoint, Topic{To: normalizeTopicName(topic), RegistrationTokens: registrationTokens})
}

func (fcm *FCM) sendFirebaseRequest(ctx context.Context, payload Message) (*Response, error) {

	// validate Message
	if err := payload.Validate(); err != nil {
		return nil, err
	}

	// Encode Message
	buf := bytes.NewBuffer([]byte{})
	err := json.NewEncoder(buf).Encode(payload)
	if err != nil {
		return nil, err
	}

	// create request
	req, err := http.NewRequest("POST", fcmSendEndpoint, buf)
	if err != nil {
		return nil, err
	}

	// add context
	req = req.WithContext(ctx)

	// add headers
	req.Header.Set("Authorization", fmt.Sprintf("key=%s", fcm.app.apiKey))
	req.Header.Set("Content-Type", "application/json")

	// execute request
	resp, err := fcm.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	defer io.Copy(ioutil.Discard, resp.Body)

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

func (fcm *FCM) sendFirebaseTopicRequest(ctx context.Context, endpoint string, payload Topic) (*Response, error) {

	// validate Topic
	if err := payload.Validate(); err != nil {
		return nil, err
	}

	// Encode Topic
	buf := bytes.NewBuffer([]byte{})
	err := json.NewEncoder(buf).Encode(payload)
	if err != nil {
		return nil, err
	}

	// create request
	req, err := http.NewRequest("POST", endpoint, buf)
	if err != nil {
		return nil, err
	}

	// add context
	req = req.WithContext(ctx)

	// add headers
	req.Header.Set("Authorization", fmt.Sprintf("key=%s", fcm.app.apiKey))
	req.Header.Set("Content-Type", "application/json")

	// execute request
	resp, err := fcm.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	defer io.Copy(ioutil.Discard, resp.Body)

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

// normalizeTopicName with prefix /topics/
func normalizeTopicName(topic string) (result string) {
	if strings.HasPrefix(topic, "/topics/") {
		return topic
	}
	return fmt.Sprint("/topics/", topic)
}
