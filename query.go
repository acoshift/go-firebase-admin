package firebase

// Query is the query interface
type Query interface {
	Ref() *Reference
	EndAt(value interface{}) Query
	EqualTo(value interface{}) Query
	IsEqual(other Query) bool
	LimitToFirst(limit int) Query
	LimitToLast(limit int) Query
	OrderByChild(path interface{}) Query
	OrderByKey() Query
	OrderByPriority() Query
	OrderByValue() Query
	StartAt(value interface{}) Query
	OnValue(event chan *DataSnapshot) CancelFunc
	OnChildAdded(event chan *ChildSnapshot) CancelFunc
	OnChildRemoved(event chan *OldChildSnapshot) CancelFunc
	OnChildChanged(event chan *ChildSnapshot) CancelFunc
	OnChildMoved(event chan *ChildSnapshot) CancelFunc
	OnceValue() (*DataSnapshot, error)
	OnceChildAdded() *ChildSnapshot
	OnceChildRemove() *OldChildSnapshot
	OnceChildChanged() *ChildSnapshot
	OnceChildMoved() *ChildSnapshot
	String() string
}

// CancelFunc is the function for cancel "On"
type CancelFunc func()
