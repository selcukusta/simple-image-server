package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/namsral/flag"
	"github.com/selcukusta/simple-image-server/internal/handler/googledrive"
	"github.com/selcukusta/simple-image-server/internal/handler/gridfs"
	"github.com/selcukusta/simple-image-server/internal/util/connection"
	"github.com/selcukusta/simple-image-server/internal/util/constant"
	"github.com/selcukusta/simple-image-server/internal/util/helper"
	"github.com/selcukusta/simple-image-server/internal/util/middleware"

	"github.com/gorilla/mux"
)

func selectHandler(w http.ResponseWriter, r *http.Request) {
	slug := mux.Vars(r)["slug"]
	switch slug {
	case "gdrive":
		googledrive.Handler(w, r)
		return
	case "gridfs":
		gridfs.Handler(w, r)
		return
	}
}

func main() {

	flag.StringVar(&constant.Hostname, "hostname", "127.0.0.1", "Specify the hostname to listen to")
	flag.StringVar(&constant.Port, "port", "8080", "Specify the port to listen to")
	flag.IntVar(&constant.CacheControlMaxAge, "cache_control_max_age", 14, "Specify the max-age for cache-control header")
	flag.StringVar(&connection.ConnectionString, "mongo_connection_str", "mongodb://127.0.0.1:27017", "Specify the connection string to connect to MongoDB instance")
	flag.Uint64Var(&connection.MaxPoolSize, "mongo_max_pool_size", 5, "Specify the max pool size for MongoDB connections")
	flag.Parse()

	var router = mux.NewRouter()
	router.Use(middleware.CommonMiddleware)
	router.NewRoute().MatcherFunc(func(r *http.Request, rm *mux.RouteMatch) bool {
		patterns := [2]string{
			`/i/(?P<slug>gdrive|gridfs)/(?P<quality_r_1_100>\d+)/(?P<width_r_0_5000>\d+)x(?P<height_r_0_5000>\d+)/(?P<option>[gts]{1,3})/(?P<path>.*)`,
			`/i/(?P<slug>gdrive|gridfs)/(?P<quality_r_1_100>\d+)/(?P<width_r_0_5000>\d+)x(?P<height_r_0_5000>\d+)/(?P<path>.*)`,
		}

		available, vars := helper.IsRouteFit(patterns, r.URL.Path)
		if available {
			rm.Vars = vars
		}
		return available
	}).HandlerFunc(selectHandler)

	srv := &http.Server{
		Handler:      router,
		Addr:         fmt.Sprintf("%s:%s", constant.Hostname, constant.Port),
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}
	log.Println(fmt.Sprintf("Server is started: %s", srv.Addr))
	log.Fatal(srv.ListenAndServe())
}
