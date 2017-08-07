package firebase

import (
	"bytes"
	"encoding/json"
)

// DataSnapshot type
type DataSnapshot struct {
	ref *Reference
	raw []byte
}

// Key returns the location of this DataSnapshot
func (snapshot *DataSnapshot) Key() string {
	panic(ErrNotImplemented)
}

// Ref returns the Reference for the location generated this DataSnapshot
func (snapshot *DataSnapshot) Ref() *Reference {
	return snapshot.ref
}

// Exists returns true if this DataSnapshot contains any data
func (snapshot *DataSnapshot) Exists() bool {
	return bytes.Compare(snapshot.raw, []byte("null")) != 0
}

// Val extracts a value from a DataSnapshot
func (snapshot *DataSnapshot) Val(v interface{}) error {
	return json.NewDecoder(bytes.NewReader(snapshot.raw)).Decode(v)
}

// Bytes returns snapshot raw data
func (snapshot *DataSnapshot) Bytes() []byte {
	return snapshot.raw
}

// ChildSnapshot type
type ChildSnapshot struct {
	PrevChildKey string
}

// OldChildSnapshot type
type OldChildSnapshot struct {
}
