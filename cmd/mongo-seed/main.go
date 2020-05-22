package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"

	"github.com/namsral/flag"

	"github.com/selcukusta/simple-image-server/internal/util/connection"
	"go.mongodb.org/mongo-driver/mongo/gridfs"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/x/bsonx"
)

func main() {

	flag.StringVar(&connection.ConnectionString, "mongo_connection_str", "mongodb://127.0.0.1:27017", "Specify the connection string to connect to MongoDB instance")
	flag.Uint64Var(&connection.MaxPoolSize, "mongo_max_pool_size", 5, "Specify the max pool size for MongoDB connections")
	flag.Parse()

	conn, err := connection.InitiateMongoClient()
	if err != nil {
		log.Fatal(err.Error())
		os.Exit(1)
	}

	bucket, err := gridfs.NewBucket(
		conn.Database("Photos"),
	)
	if err != nil {
		log.Fatal(err.Error())
		os.Exit(1)
	}

	files, err := ioutil.ReadDir("assets")
	if err != nil {
		log.Fatal(err.Error())
	}

	for _, fi := range files {
		data, err := ioutil.ReadFile(path.Join("assets", fi.Name()))
		if err != nil {
			log.Fatal(err.Error())
		}

		opts := options.GridFSUpload()
		opts.SetMetadata(bsonx.Doc{{Key: "Content-Type", Value: bsonx.String("image/jpeg")}})

		uploadStream, err := bucket.OpenUploadStream(fi.Name(), opts)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		defer uploadStream.Close()

		_, err = uploadStream.Write(data)
		if err != nil {
			log.Fatal(err.Error())
		}
	}
	log.Println(`Database seed operation is completed successfully!`)
}
