package clouditor

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"gopkg.in/mgo.v2/bson"
)

type LegacyDatabase struct {
	client   *mongo.Client
	database *mongo.Database
}

func (db *LegacyDatabase) Connect() (err error) {
	if db.client, err = mongo.NewClient(options.Client().ApplyURI("mongodb://localhost:27017")); err != nil {
		return
	}

	if err = db.client.Connect(context.Background()); err != nil {
		return
	}

	db.database = db.client.Database("clouditor")

	return
}

func (db *LegacyDatabase) GetAccountById(id string, account ServiceAccount) (err error) {
	collection := db.database.Collection("accounts")

	return collection.FindOne(context.Background(), bson.M{"_id": id}).Decode(account)
}
