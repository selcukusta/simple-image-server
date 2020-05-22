package connection

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

//InitiateMongoClient is used to connect MongoDB instance
func InitiateMongoClient() (*mongo.Client, error) {
	var err error
	var client *mongo.Client

	opts := options.Client()
	opts.ApplyURI(ConnectionString)
	opts.SetMaxPoolSize(MaxPoolSize)
	if client, err = mongo.Connect(context.Background(), opts); err != nil {
		return nil, err
	}
	return client, nil
}
