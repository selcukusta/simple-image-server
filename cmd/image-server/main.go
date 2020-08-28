package main

import (
	"fmt"
	"log"
	"time"

	"github.com/namsral/flag"
	"github.com/selcukusta/simple-image-server/internal/handler/abs"
	"github.com/selcukusta/simple-image-server/internal/handler/googledrive"
	"github.com/selcukusta/simple-image-server/internal/handler/gridfs"
	"github.com/selcukusta/simple-image-server/internal/handler/version"
	"github.com/selcukusta/simple-image-server/internal/util/connection"
	"github.com/selcukusta/simple-image-server/internal/util/constant"
	"github.com/selcukusta/simple-image-server/internal/util/helper"
	"github.com/selcukusta/simple-image-server/internal/util/logger"
	"github.com/selcukusta/simple-image-server/internal/util/middleware"
	"github.com/valyala/fasthttp"
)

func requestHandler(ctx *fasthttp.RequestCtx) {
	path := string(ctx.Path())

	switch path {
	case "/version":
		version.Handler(ctx)
		return
	default:
		patterns := [2]string{
			`/i/(?P<slug>gdrive|gridfs|abs)/(?P<quality_r_1_100>\d+)/(?P<width_r_0_5000>\d+)x(?P<height_r_0_5000>\d+)/(?P<option>[gtc]{1,3})/(?P<path>.*)`,
			`/i/(?P<slug>gdrive|gridfs|abs)/(?P<quality_r_1_100>\d+)/(?P<width_r_0_5000>\d+)x(?P<height_r_0_5000>\d+)/(?P<path>.*)`,
		}

		available, vars := helper.IsRouteFit(patterns, path)
		if !available {
			ctx.SetStatusCode(fasthttp.StatusNotFound)
			return
		}

		slug := vars["slug"]
		defer helper.TraceObject{HandlerName: slug, Parameter: path}.TimeTrack(time.Now())
		switch slug {
		case "gdrive":
			googledrive.Handler(ctx, vars)
			return
		case "gridfs":
			gridfs.Handler(ctx, vars)
		case "abs":
			abs.Handler(ctx, vars)
			return
		}
	}
}

func main() {
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
	flag.Parse()

	handler := requestHandler
	handler = middleware.CommonMiddleware(handler)

	addr := fmt.Sprintf("%s:%s", constant.Hostname, constant.Port)
	logger.WriteLog(logger.Log{Level: logger.INFO, Message: fmt.Sprintf("Server is started: %s", addr)})
	log.Fatal(fasthttp.ListenAndServe(addr, handler))
}
