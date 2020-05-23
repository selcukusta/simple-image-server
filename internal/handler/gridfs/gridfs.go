package gridfs

import (
	"bytes"
	"context"
	"time"

	"github.com/selcukusta/simple-image-server/internal/util/connection"
	"github.com/selcukusta/simple-image-server/internal/util/model"
	"github.com/valyala/fasthttp"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/gridfs"
)

type gridFile struct {
	ID       primitive.ObjectID `json:"_id" bson:"_id"`
	FileName string             `json:"filename" bson:"filename"`
	Metadata struct {
		ContentType string `json:"content-type" bson:"content-type"`
	} `json:"metadata" bson:"metadata"`
}

//Handler is using connect to MongoDB and get the image
func Handler(ctx *fasthttp.RequestCtx, vars map[string]string) {
	path := vars["path"]

	_, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	conn, err := connection.InitiateMongoClient()
	if err != nil {
		customError := model.CustomError{Message: "Unable to connect MongoDB instance", Detail: err}
		model.FailedFinalizer{ResponseWriter: ctx, StdOut: &customError}.Finalize()
		return
	}

	db := conn.Database("Photos")

	bucket, err := gridfs.NewBucket(db)
	if err != nil {
		customError := model.CustomError{Message: "Unable to connect GridFS bucket", Detail: err}
		model.FailedFinalizer{ResponseWriter: ctx, StdOut: &customError}.Finalize()
		return
	}

	fsFiles := db.Collection("fs.files")
	objectID, err := primitive.ObjectIDFromHex(path)
	if err != nil {
		customError := model.CustomError{Message: "Unable to parse value to ObjectID", Detail: err}
		model.FailedFinalizer{ResponseWriter: ctx, StdOut: &customError}.Finalize()
		return
	}

	var fileInfo gridFile
	err = fsFiles.FindOne(ctx, bson.M{"_id": objectID}).Decode(&fileInfo)
	if err != nil {
		customError := model.CustomError{Message: "Unable to get file info", Detail: err}
		model.FailedFinalizer{ResponseWriter: ctx, StdOut: &customError}.Finalize()
		return
	}

	var buf bytes.Buffer
	_, err = bucket.DownloadToStream(objectID, &buf)
	if err != nil {
		customError := model.CustomError{Message: "Unable to download stream", Detail: err}
		model.FailedFinalizer{ResponseWriter: ctx, StdOut: &customError}.Finalize()
		return
	}

	model.SucceededFinalizer{ResponseWriter: ctx, ContentType: fileInfo.Metadata.ContentType}.Finalize(vars, buf.Bytes())
}
