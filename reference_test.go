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
	rr, err := r1.Push(n)
	if err != nil {
		t.Fatal(err)
	}
	if rr == nil {
		t.Fatalf("expected push result not nil; got nil")
	}
	if rr.Key() == r1.Key() {
		t.Fatalf("expected push result key changed")
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

type dinosaurs struct {
	Appeared int64   `json:"appeared"`
	Height   float32 `json:"height"`
	Length   float32 `json:"length"`
	Order    string  `json:"order"`
	Vanished int64   `json:"vanished"`
	Weight   int     `json:"weight"`
}

func TestOrderBy(t *testing.T) {
	app := initApp()
	db := app.Database()
	var err error

	// feed data
	r := db.Ref("test/path")
	defer r.Remove()
	err = r.Child("bruhathkayosaurus").Set(&dinosaurs{-70000000, 25, 44, "saurischia", -70000000, 135000})
	if err != nil {
		t.Fatal(err)
	}
	err = r.Child("lambeosaurus").Set(&dinosaurs{-76000000, 2.1, 12.5, "ornithischia", -75000000, 5000})
	if err != nil {
		t.Fatal(err)
	}
	err = r.Child("linhenykus").Set(&dinosaurs{-85000000, 0.6, 1, "theropoda", -75000000, 3})
	if err != nil {
		t.Fatal(err)
	}
	err = r.Child("pterodactyl").Set(&dinosaurs{-150000000, 0.6, 0.8, "pterosauria", -148500000, 2})
	if err != nil {
		t.Fatal(err)
	}
	err = r.Child("stegosaurus").Set(&dinosaurs{-155000000, 4, 9, "ornithischia", -150000000, 2500})
	if err != nil {
		t.Fatal(err)
	}
	err = r.Child("triceratops").Set(&dinosaurs{-68000000, 3, 8, "ornithischia", -66000000, 11000})
	if err != nil {
		t.Fatal(err)
	}

	snapshot, err := r.OrderByKey().OnceValue()
	if err != nil {
		t.Fatal(err)
	}
	var d map[string]*dinosaurs
	snapshot.Val(&d)
	if len(d) != 6 {
		t.Fatalf("expected dinosours have 6 len; got %d", len(d))
	}

	snapshot, err = r.OrderByChild("height").EqualTo(0.6).OnceValue()
	if err != nil {
		t.Fatal(err)
	}
	d = nil
	snapshot.Val(&d)
	if len(d) != 2 {
		t.Fatalf("expected dinosours have 2 len; got %d", len(d))
	}
}
