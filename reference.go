package admin

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	_path "path"
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

func (ref *Reference) invokeRequest(method string, body io.Reader) ([]byte, error) {
	url, err := ref.url()
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json; charset=utf-8")
	resp, err := ref.database.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("firebasedatabase: %s", resp.Status)
	}
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return bytes.TrimSpace(b), nil
}

// Set writes data to current location
func (ref *Reference) Set(value interface{}) error {
	buf := bytes.NewBuffer([]byte{})
	err := json.NewEncoder(buf).Encode(value)
	if err != nil {
		return err
	}
	_, err = ref.invokeRequest(http.MethodPut, buf)
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
	_, err = ref.invokeRequest(http.MethodPost, buf)
	if err != nil {
		return err
	}
	return nil
}

// Remove removes data from current location
func (ref *Reference) Remove() error {
	_, err := ref.invokeRequest(http.MethodDelete, nil)
	if err != nil {
		return err
	}
	return nil
}

// Key returns the last path of Reference
func (ref *Reference) Key() string {
	_, p := _path.Split(ref.path)
	return p
}

// Ref returns a copy
func (ref Reference) Ref() *Reference {
	return &ref
}

// Root returns the root location of database
func (ref *Reference) Root() *Reference {
	return &Reference{
		database: ref.database,
	}
}

// Child returns a Reference for relative path
func (ref Reference) Child(path string) *Reference {
	ref.path = _path.Join(ref.path, path)
	return &ref
}

// EndAt implements Query interface
func (ref *Reference) EndAt(value interface{}, key string) Query {
	panic(ErrNotImplement)
}

// StartAt implements Query interface
func (ref *Reference) StartAt(value interface{}, key string) Query {
	panic(ErrNotImplement)
}

// EqualTo implements Query interface
func (ref *Reference) EqualTo(value interface{}, key string) Query {
	panic(ErrNotImplement)
}

// IsEqual implements Query interface
func (ref *Reference) IsEqual(other interface{}) Query {
	panic(ErrNotImplement)
}

// LimitToFirst implements Query interface
func (ref *Reference) LimitToFirst(limit int) Query {
	panic(ErrNotImplement)
}

// LimitToLast implements Query interface
func (ref *Reference) LimitToLast(limit int) Query {
	panic(ErrNotImplement)
}

// OrderByChild implements Query interface
func (ref *Reference) OrderByChild(path string) Query {
	panic(ErrNotImplement)
}

// OrderByKey implements Query interface
func (ref *Reference) OrderByKey() Query {
	panic(ErrNotImplement)
}

// OrderByPriority implements Query interface
func (ref *Reference) OrderByPriority() Query {
	panic(ErrNotImplement)
}

// OrderByValue implements Query interface
func (ref *Reference) OrderByValue() Query {
	panic(ErrNotImplement)
}

// OnValue implements Query interface
func (ref *Reference) OnValue(event chan *DataSnapshot) CancelFunc {
	panic(ErrNotImplement)
}

// OnChildAdded implements Query interface
func (ref *Reference) OnChildAdded(event chan *ChildSnapshot) CancelFunc {
	panic(ErrNotImplement)
}

// OnChildRemoved implements Query interface
func (ref *Reference) OnChildRemoved(event chan *OldChildSnapshot) CancelFunc {
	panic(ErrNotImplement)
}

// OnChildChanged implements Query interface
func (ref *Reference) OnChildChanged(event chan *ChildSnapshot) CancelFunc {
	panic(ErrNotImplement)
}

// OnChildMoved implements Query interface
func (ref *Reference) OnChildMoved(event chan *ChildSnapshot) CancelFunc {
	panic(ErrNotImplement)
}

// OnceValue implements Query interface
func (ref *Reference) OnceValue() (*DataSnapshot, error) {
	// TODO: find from cached first
	b, err := ref.invokeRequest(http.MethodGet, nil)
	if err != nil {
		return nil, err
	}
	return &DataSnapshot{
		ref: ref,
		raw: b,
	}, nil
}

// OnceChildAdded implements Query interface
func (ref *Reference) OnceChildAdded() *ChildSnapshot {
	panic(ErrNotImplement)
}

// OnceChildRemove implements Query interface
func (ref *Reference) OnceChildRemove() *OldChildSnapshot {
	panic(ErrNotImplement)
}

// OnceChanged implements Query interface
func (ref *Reference) OnceChanged() *ChildSnapshot {
	panic(ErrNotImplement)
}

// OnceMoved implements Query interface
func (ref *Reference) OnceMoved() *ChildSnapshot {
	panic(ErrNotImplement)
}

func (ref *Reference) String() string {
	panic(ErrNotImplement)
}
