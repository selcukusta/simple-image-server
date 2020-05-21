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

	flag.StringVar(&connection.Hostname, "mongo_hostname", "127.0.0.1", "Specify the hostname to connect to MongoDB instance")
	flag.StringVar(&connection.Port, "mongo_port", "27017", "Specify the port to connect to MongoDB instance")
	flag.Uint64Var(&connection.MaxPoolSize, "mongo_max_pool_size", 5, "Specify the max pool size for MongoDB connections")
	flag.Parse()

	conn := connection.InitiateMongoClient()
	bucket, err := gridfs.NewBucket(
		conn.Database("Photos"),
	)
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}

	files, err := ioutil.ReadDir("assets")
	if err != nil {
		log.Fatal(err)
	}

	for _, fi := range files {
		data, err := ioutil.ReadFile(path.Join("assets", fi.Name()))
		if err != nil {
			log.Fatal(err)
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
			log.Fatal(err)
		}
	}
	log.Println(`Database seed operation is completed successfully!`)
}
