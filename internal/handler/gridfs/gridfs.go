package gridfs

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/selcukusta/simple-image-server/internal/util/connection"
	"github.com/selcukusta/simple-image-server/internal/util/constant"
	"github.com/selcukusta/simple-image-server/internal/util/helper"
	"github.com/selcukusta/simple-image-server/internal/util/model"
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
func Handler(w http.ResponseWriter, r *http.Request) {
	path := mux.Vars(r)["path"]
	defer helper.TraceObject{HandlerName: "GridFS", Parameter: path}.TimeTrack(time.Now())

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	conn, err := connection.InitiateMongoClient()
	if err != nil {
		log.Println(fmt.Sprintf(constant.LogErrorFormat, "Unable to connect MongoDB instance", err.Error()))
		w.WriteHeader(http.StatusInternalServerError)
		_, err = w.Write([]byte(constant.ErrorMessage))
		if err != nil {
			log.Println(fmt.Sprintf(constant.LogErrorFormat, constant.LogErrorMessage, err.Error()))
		}
		return
	}

	db := conn.Database("Photos")

	bucket, err := gridfs.NewBucket(db)
	if err != nil {
		log.Println(fmt.Sprintf(constant.LogErrorFormat, "Unable to connect GridFS bucket", err.Error()))
		w.WriteHeader(http.StatusInternalServerError)
		_, err = w.Write([]byte(constant.ErrorMessage))
		if err != nil {
			log.Println(fmt.Sprintf(constant.LogErrorFormat, constant.LogErrorMessage, err.Error()))
		}
		return
	}

	fsFiles := db.Collection("fs.files")
	objectID, err := primitive.ObjectIDFromHex(path)
	if err != nil {
		log.Println(fmt.Sprintf(constant.LogErrorFormat, "Unable to parse value to ObjectID", err.Error()))
		w.WriteHeader(http.StatusInternalServerError)
		_, err = w.Write([]byte(constant.ErrorMessage))
		if err != nil {
			log.Println(fmt.Sprintf(constant.LogErrorFormat, constant.LogErrorMessage, err.Error()))
		}
		return
	}

	var fileInfo gridFile
	err = fsFiles.FindOne(ctx, bson.M{"_id": objectID}).Decode(&fileInfo)
	if err != nil {
		log.Println(fmt.Sprintf(constant.LogErrorFormat, "Unable to get file info", err.Error()))
		w.WriteHeader(http.StatusInternalServerError)
		_, err = w.Write([]byte(constant.ErrorMessage))
		if err != nil {
			log.Println(fmt.Sprintf(constant.LogErrorFormat, constant.LogErrorMessage, err.Error()))
		}
		return
	}

	var buf bytes.Buffer
	_, err = bucket.DownloadToStream(objectID, &buf)
	if err != nil {
		log.Println(fmt.Sprintf(constant.LogErrorFormat, "Unable to download stream", err.Error()))
		w.WriteHeader(http.StatusInternalServerError)
		_, err = w.Write([]byte(constant.ErrorMessage))
		if err != nil {
			log.Println(fmt.Sprintf(constant.LogErrorFormat, constant.LogErrorMessage, err.Error()))
		}
		return
	}

	finalizer := model.HandlerFinalizer{ResponseWriter: w, Headers: nil}
	finalizer.Finalize(mux.Vars(r), buf.Bytes(), fileInfo.Metadata.ContentType)
}
