package db

import (
	"context"
	"fmt"
	"os"

	firebase "firebase.google.com/go/v4"
	"firebase.google.com/go/v4/db"
	"google.golang.org/api/option"
)

var (
	app       *firebase.App
	FbDB      *db.Client
	appMulti  map[string]*firebase.App
	FbDBMulti map[string]*db.Client
	Ctx       context.Context
)

func NewFirebaseData() {
	NewFirebase()
	NewFirebaseMultiInstances("STAGES")
}

func NewFirebase() error {
	var (
		err    error
		dbName string
	)
	contxt := context.Background()
	filename := fmt.Sprintf("%s/firebase.json", os.Getenv("FIREBASE_CRED_PATH"))
	dbName = os.Getenv("FIREBASE_DATABASE")
	if dbName == "" {
		if os.Getenv("FIREBASE_SERVER") == "production" {
			dbName = "https://bina-taruna-wiratama.firebaseio.com"
		} else {
			dbName = "https://btw-push.firebaseio.com"
		}
	}
	opt := option.WithCredentialsFile(filename)
	conf := &firebase.Config{
		DatabaseURL: dbName,
	}

	app, err = firebase.NewApp(contxt, conf, opt)
	if err != nil {
		return err
	}
	FbDB, err = app.Database(contxt)

	if err != nil {
		return err
	}
	Ctx = contxt

	return nil
}

func NewFirebaseMultiInstances(instanceName string) error {
	var (
		err    error
		dbName string
	)
	contxt := context.Background()
	filename := fmt.Sprintf("%s/firebase.json", os.Getenv("FIREBASE_CRED_PATH"))
	dbName = os.Getenv(fmt.Sprintf("FIREBASE_DATABASE_%s", instanceName))
	if dbName == "" {
		dbName = os.Getenv("FIREBASE_DATABASE")
		if dbName == "" {
			if os.Getenv("FIREBASE_SERVER") == "production" {
				dbName = "https://bina-taruna-wiratama.firebaseio.com"
			} else {
				dbName = "https://btw-push.firebaseio.com"
			}
		}
	}
	opt := option.WithCredentialsFile(filename)
	conf := &firebase.Config{
		DatabaseURL: dbName,
	}
	if appMulti == nil {
		appMulti = map[string]*firebase.App{}
	}
	if FbDBMulti == nil {
		FbDBMulti = map[string]*db.Client{}
	}

	appMulti[instanceName], err = firebase.NewApp(contxt, conf, opt)
	if err != nil {
		return err
	}
	FbDBMulti[instanceName], err = appMulti[instanceName].Database(contxt)

	if err != nil {
		return err
	}
	// Ctx = contxt

	return nil
}
