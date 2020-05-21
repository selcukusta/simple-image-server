package connection

import (
	"context"
	"fmt"
	"log"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

//InitiateMongoClient is used to connect MongoDB instance
func InitiateMongoClient() *mongo.Client {
	var err error
	var client *mongo.Client

	uri := fmt.Sprintf("mongodb://%s:%s", Hostname, Port)
	opts := options.Client()
	opts.ApplyURI(uri)
	opts.SetMaxPoolSize(MaxPoolSize)
	if client, err = mongo.Connect(context.Background(), opts); err != nil {
		log.Fatal(err.Error())
	}
	return client
}
