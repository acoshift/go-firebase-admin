package admin_test

import (
	"log"
	"testing"
	"time"

	. "github.com/acoshift/go-firebase-admin"
)

func TestDatabase(t *testing.T) {
	app, _ := InitializeApp(AppOptions{
		DatabaseURL: "https://acourse-acoshift.firebaseio.com",
	})

	db := app.Database()
	start := time.Now()
	db.Ref("temp").Push(map[string]string{"test": "test abc"})
	end := time.Now()
	log.Println(end.Sub(start).String())
	start = time.Now()
	db.Ref("temp").Push(map[string]string{"test": "test abc"})
	end = time.Now()
	log.Println(end.Sub(start).String())
	start = time.Now()
	db.Ref("temp").Push(map[string]string{"test": "test abc"})
	end = time.Now()
	log.Println(end.Sub(start).String())
}
