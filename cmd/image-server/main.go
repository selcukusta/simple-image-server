package main

import (
	"fmt"
	"log"

	"time"

	"github.com/namsral/flag"
	"github.com/selcukusta/simple-image-server/internal/handler/abs"
	"github.com/selcukusta/simple-image-server/internal/handler/googledrive"
	"github.com/selcukusta/simple-image-server/internal/handler/gridfs"
	"github.com/selcukusta/simple-image-server/internal/handler/s3"
	"github.com/selcukusta/simple-image-server/internal/handler/url"
	"github.com/selcukusta/simple-image-server/internal/handler/version"
	"github.com/selcukusta/simple-image-server/internal/util/connection"
	"github.com/selcukusta/simple-image-server/internal/util/constant"
	"github.com/selcukusta/simple-image-server/internal/util/helper"
	"github.com/selcukusta/simple-image-server/internal/util/logger"
	"github.com/selcukusta/simple-image-server/internal/util/middleware"
	"github.com/sirupsen/logrus"
	"github.com/valyala/fasthttp"
)

func requestHandler(ctx *fasthttp.RequestCtx) {
	path := string(ctx.Path())

	switch path {
	case "/version":
		version.Handler(ctx)
		return
	default:
		available, vars := helper.IsRouteFit(constant.Patterns, path)
		if !available {
			ctx.SetStatusCode(fasthttp.StatusNotFound)
			return
		}

		slug := vars["slug"]
		defer helper.TraceObject{HandlerName: slug, Parameter: path, Rq: ctx}.TimeTrack(time.Now())
		switch slug {
		case "gdrive":
			googledrive.Handler(ctx, vars)
			return
		case "gridfs":
			gridfs.Handler(ctx, vars)
			return
		case "abs":
			abs.Handler(ctx, vars)
		case "s3":
			s3.Handler(ctx, vars)
			return
		case "url":
		 	url.Handler(ctx, vars)
			return
		}
	}
}

func main() {
	//Set logger

	logrus.SetFormatter(&logrus.JSONFormatter{
		FieldMap: logrus.FieldMap{
			logrus.FieldKeyLevel: "Level",
			logrus.FieldKeyTime:  "@timestamp",
			logrus.FieldKeyMsg:   "Message",
		},
	})
	// Hosting flags
	flag.StringVar(&constant.Hostname, "hostname", "127.0.0.1", "Specify the hostname to listen to")
	flag.StringVar(&constant.Port, "port", "8080", "Specify the port to listen to")
	flag.IntVar(&constant.CacheControlMaxAge, "cache_control_max_age", 14, "Specify the max-age for cache-control header")

	// MongoDB flags
	flag.StringVar(&connection.ConnectionString, "mongo_connection_str", "mongodb://127.0.0.1:27017", "Specify the connection string to connect to MongoDB instance")
	flag.StringVar(&connection.DBName, "mongo_db_name", "Photos", "Specify the DB name to determine which database will be used to store the images")
	flag.Uint64Var(&connection.MaxPoolSize, "mongo_max_pool_size", 5, "Specify the max pool size for MongoDB connections")

	// Azure Blog Storage flags
	flag.StringVar(&connection.AccountKey, "abs_account_key", "", "Specify the account key to connect Azure Blob Storage account")
	flag.StringVar(&connection.AccountName, "abs_account_name", "", "Specify the account name to connect Azure Blob Storage account")
	flag.StringVar(&connection.AzureURI, "abs_azure_uri", "", "Specify the URI to connect Azure Blob Storage account")

	// s3 flags
	flag.StringVar(&connection.S3Name, "s3_name", "", "Specify the S3 name to connect S3 Storage account")
	flag.StringVar(&connection.S3Region, "s3_region", "", "Specify the S3 region to connect S3 Storage account")

	// URL flags
	flag.StringVar(&connection.URL, "url", "", "Specify the url address")

	flag.Parse()

	handler := requestHandler
	handler = middleware.CommonMiddleware(handler)

	addr := fmt.Sprintf("%s:%s", constant.Hostname, constant.Port)
	logger.Init().Info(fmt.Sprintf("Server is started: %s", addr))
	log.Fatal(fasthttp.ListenAndServe(addr, handler))
}
