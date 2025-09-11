package db

import (
	"context"
	"errors"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
	"go.mongodb.org/mongo-driver/v2/mongo/readpref"
	"time"
)

type DB struct {
	Client *mongo.Client
}

func NewDB(uri string) (*DB, error) {
	db := &DB{}

	if _, err := db.Connection(uri); err != nil {
		return nil, err
	}

	return db, nil
}

func (d *DB) Connection(uri string) (*mongo.Client, error) {
	client, err := connect(uri)
	if err != nil {
		return nil, err
	}

	err = ping(client)
	if err != nil {
		return nil, err
	}

	println("database is connected")

	d.Client = client

	return d.Client, nil
}

func connect(uri string) (*mongo.Client, error) {
	if uri == "" {
		return nil, errors.New("database connection not provided")
	}

	client, err := mongo.Connect(options.Client().ApplyURI(uri))
	if err != nil {
		return nil, err
	}

	return client, nil
}

func ping(client *mongo.Client) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := client.Ping(ctx, readpref.Primary()); err != nil {
		return err
	}

	return nil
}

func (d *DB) Disconnect() error {
	if err := d.Client.Disconnect(context.TODO()); err != nil {
		return err
	}

	println("database connection closed")
	return nil
}
