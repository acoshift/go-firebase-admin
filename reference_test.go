package admin_test

import (
	"testing"
	"time"
)

func TestReference(t *testing.T) {
	app := initApp()
	db := app.Database()

	r1 := db.Ref("test/path")
	err := r1.Push(time.Now())
	if err != nil {
		t.Fatal(err)
	}
	err = r1.Set(time.Now())
	if err != nil {
		t.Fatal(err)
	}
	err = r1.Remove()
	if err != nil {
		t.Fatal(err)
	}
}
