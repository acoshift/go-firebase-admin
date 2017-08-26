package firebase_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestReference(t *testing.T) {
	app := initApp(t)
	db := app.Database()

	err := db.Ref(".settings/rules").Set(map[string]interface{}{
		"rules": map[string]interface{}{
			".read":  false,
			".write": false,
			"test": map[string]interface{}{
				"path": map[string]interface{}{
					".indexOn": []string{"height"},
				},
			},
		},
	})
	assert.NoError(t, err)

	r1 := db.Ref("test/path")
	n := time.Now()
	rr, err := r1.Push(n)
	assert.NoError(t, err)
	assert.NotNil(t, rr, "expected push result not nil")
	assert.NotEqual(t, rr.Key(), r1.Key(), "expected push result key changed")

	err = r1.Set(n)
	assert.NoError(t, err)

	snapshot, err := r1.OnceValue()
	assert.NoError(t, err)
	assert.True(t, snapshot.Exists())

	var ts time.Time
	err = snapshot.Val(&ts)
	assert.NoError(t, err)
	assert.True(t, n.Equal(ts), "expected data to be %v; got %v", n, ts)

	err = r1.Remove()
	assert.NoError(t, err)

	snapshot, err = r1.OnceValue()
	assert.NoError(t, err)
	assert.False(t, snapshot.Exists(), "expected snapshot not exists")
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
	app := initApp(t)
	db := app.Database()
	var err error

	// feed data
	r := db.Ref("test/path")
	defer r.Remove()
	err = r.Child("bruhathkayosaurus").Set(&dinosaurs{-70000000, 25, 44, "saurischia", -70000000, 135000})
	assert.NoError(t, err)

	err = r.Child("lambeosaurus").Set(&dinosaurs{-76000000, 2.1, 12.5, "ornithischia", -75000000, 5000})
	assert.NoError(t, err)

	err = r.Child("linhenykus").Set(&dinosaurs{-85000000, 0.6, 1, "theropoda", -75000000, 3})
	assert.NoError(t, err)

	err = r.Child("pterodactyl").Set(&dinosaurs{-150000000, 0.6, 0.8, "pterosauria", -148500000, 2})
	assert.NoError(t, err)

	err = r.Child("stegosaurus").Set(&dinosaurs{-155000000, 4, 9, "ornithischia", -150000000, 2500})
	assert.NoError(t, err)

	err = r.Child("triceratops").Set(&dinosaurs{-68000000, 3, 8, "ornithischia", -66000000, 11000})
	assert.NoError(t, err)

	snapshot, err := r.OrderByKey().OnceValue()
	assert.NoError(t, err)

	var d map[string]*dinosaurs
	snapshot.Val(&d)
	assert.Len(t, d, 6, "expected dinosours have 6 len; got %d", len(d))

	snapshot, err = r.OrderByChild("height").EqualTo(0.6).OnceValue()
	assert.NoError(t, err)

	d = nil
	snapshot.Val(&d)
	assert.Len(t, d, 2, "expected dinosours have 2 len; got %d", len(d))
}
