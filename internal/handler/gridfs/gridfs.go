package gridfs

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/selcukusta/simple-image-server/internal/processor"
	"github.com/selcukusta/simple-image-server/internal/util/connection"
	"github.com/selcukusta/simple-image-server/internal/util/constant"
	"github.com/selcukusta/simple-image-server/internal/util/helper"
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

//GridFSHandler is using connect to MongoDB and get the image
func GridFSHandler(w http.ResponseWriter, r *http.Request) {
	path := mux.Vars(r)["path"]
	defer helper.TimeTrack(time.Now(), path)

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

	result, errMessage, err := processor.ImageProcess(mux.Vars(r), buf.Bytes(), fileInfo.Metadata.ContentType)
	if err != nil {
		log.Println(fmt.Sprintf(constant.LogErrorFormat, errMessage, err.Error()))
		w.WriteHeader(http.StatusInternalServerError)
		_, err = w.Write([]byte(constant.ErrorMessage))
		if err != nil {
			log.Println(fmt.Sprintf(constant.LogErrorFormat, constant.LogErrorMessage, err.Error()))
		}
		return
	}

	if constant.CacheControlMaxAge != -1 {
		maxAge := constant.CacheControlMaxAge * 24 * 60 * 60
		w.Header().Add("Cache-Control", fmt.Sprintf("public, max-age=%d", maxAge))
	}
	_, err = w.Write(result)
	if err != nil {
		log.Println(fmt.Sprintf(constant.LogErrorFormat, constant.LogErrorMessage, err.Error()))
	}
}
