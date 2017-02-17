package admin

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
)

// Reference represents a specific location in Database
type Reference struct {
	database *Database
	path     string
}

func (ref *Reference) url() (string, error) {
	if ref.database.app.tokenSource == nil {
		return ref.database.app.databaseURL + "/" + ref.path + ".json", nil
	}
	tk, err := ref.database.app.tokenSource.Token()
	if err != nil {
		return "", err
	}
	token := tk.AccessToken
	return ref.database.app.databaseURL + "/" + ref.path + ".json?access_token=" + token, nil
}

func (ref *Reference) invokeRequest(method string, body io.Reader, v interface{}) error {
	url, err := ref.url()
	if err != nil {
		return err
	}
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json; charset=utf-8")
	resp, err := ref.database.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("firebasedatabase: %s", resp.Status)
	}
	if v == nil {
		io.Copy(ioutil.Discard, resp.Body)
		return nil
	}
	err = json.NewDecoder(resp.Body).Decode(v)
	if err != nil {
		return err
	}
	return nil
}

// Set writes data to current location
func (ref *Reference) Set(value interface{}) error {
	buf := bytes.NewBuffer([]byte{})
	err := json.NewEncoder(buf).Encode(value)
	if err != nil {
		return err
	}
	err = ref.invokeRequest(http.MethodPut, buf, nil)
	if err != nil {
		return err
	}
	return nil
}

// Push pushs data to current location
func (ref *Reference) Push(value interface{}) error {
	buf := bytes.NewBuffer([]byte{})
	err := json.NewEncoder(buf).Encode(value)
	if err != nil {
		return err
	}
	err = ref.invokeRequest(http.MethodPost, buf, nil)
	if err != nil {
		return err
	}
	return nil
}

// Remove removes data from current location
func (ref *Reference) Remove() error {
	err := ref.invokeRequest(http.MethodDelete, nil, nil)
	if err != nil {
		return err
	}
	return nil
}
