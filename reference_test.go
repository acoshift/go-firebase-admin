package admin_test

import (
	"testing"
	"time"
)

func TestReference(t *testing.T) {
	app := initApp()
	db := app.Database()

	r1 := db.Ref("test/path")
	n := time.Now()
	err := r1.Push(n)
	if err != nil {
		t.Fatal(err)
	}
	err = r1.Set(n)
	if err != nil {
		t.Fatal(err)
	}

	snapshot, err := r1.OnceValue()
	if err != nil {
		t.Fatal(err)
	}
	if !snapshot.Exists() {
		t.Fatal("expected snapshot exists")
	}
	var ts time.Time
	err = snapshot.Val(&ts)
	if err != nil {
		t.Fatal(err)
	}
	if !n.Equal(ts) {
		t.Fatalf("expected data to be %v; got %v", n, ts)
	}

	err = r1.Remove()
	if err != nil {
		t.Fatal(err)
	}
	snapshot, err = r1.OnceValue()
	if err != nil {
		t.Fatal(err)
	}
	if snapshot.Exists() {
		t.Fatal("expected snapshot not exists")
	}
}
