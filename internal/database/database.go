package database

import (
	"context"
	"os"
	"time"

	"github.com/s-alad/tiktok-of-alexandria/internal/tikwm"
	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Status string

const (
	Success Status = "success"
	Failure Status = "failure"
)

type MDB struct {
	id     string
	client *mongo.Client
}

func (m *MDB) ID() string {
	return m.id
}

func (m *MDB) Save(uid, mediaID string, location string, date string, status Status, mediaType tikwm.MediaType) error {
	collection := m.client.Database("tiktok-library").Collection("saves")
	parsed, _ := time.Parse("2006-01-02 15:04:05", date)
	filter := bson.D{
		{Key: "uid", Value: uid},
		{Key: "media", Value: mediaID},
	}
	update := bson.D{
		{Key: "$setOnInsert", Value: bson.D{
			{Key: "uid", Value: uid},
			{Key: "media", Value: mediaID},
			{Key: "location", Value: location},
			{Key: "date", Value: parsed},
			{Key: "status", Value: status},
			{Key: "type", Value: mediaType},
		}},
	}

	opts := options.Update().SetUpsert(true)
	_, err := collection.UpdateOne(context.TODO(), filter, update, opts)
	if err != nil {
		return err
	}

	log.WithFields(log.Fields{
		"uid":   uid,
		"media": mediaID,
		"date":  parsed,
	}).Info("saved media")

	return nil
}

func (m *MDB) Saves(uid string) ([]string, error) {
	collection := m.client.Database("tiktok-library").Collection("saves")
	filter := bson.D{{Key: "uid", Value: uid}}

	cursor, err := collection.Find(context.TODO(), filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.TODO())

	var media []string

	var result struct {
		media string `bson:"media"`
	}

	for cursor.Next(context.TODO()) {
		if err := cursor.Decode(&result); err != nil {
			return nil, err
		}
		media = append(media, result.media)
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return media, nil
}

func (m *MDB) init() {
	mdbURI := os.Getenv("MDB_URI")
	log.WithField("uri", mdbURI).Info("connecting to MongoDB")

	opts := options.Client().ApplyURI(mdbURI).SetServerAPIOptions(options.ServerAPI(options.ServerAPIVersion1))
	client, err := mongo.Connect(context.TODO(), opts)
	if err != nil {
		log.Error(err)
		panic(err)
	}

	if err := client.Database("admin").RunCommand(context.TODO(), bson.D{{Key: "ping", Value: 1}}).Err(); err != nil {
		log.Error(err)
		panic(err)
	}

	log.Info("connected to MongoDB")
	m.client = client
}

func (m *MDB) Close() {
	if err := m.client.Disconnect(context.TODO()); err != nil {
		log.Error(err)
		panic(err)
	}
}

func Create(id string) *MDB {
	m := &MDB{id: id}
	m.init()
	return m
}
