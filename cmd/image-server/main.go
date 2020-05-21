package main

import (
	"fmt"
	"log"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/namsral/flag"
	"github.com/selcukusta/simple-image-server/internal/handler/googledrive"
	"github.com/selcukusta/simple-image-server/internal/handler/gridfs"
	"github.com/selcukusta/simple-image-server/internal/util/connection"
	"github.com/selcukusta/simple-image-server/internal/util/constant"

	"github.com/gorilla/mux"
)

func validateRangeParams(value string, minValue int, maxValue int) bool {
	numeric, _ := strconv.Atoi(value)
	return numeric >= minValue && numeric <= maxValue
}

func commonMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Access-Control-Allow-Origin", "*")
		w.Header().Add("Access-Control-Allow-Methods", "GET")
		next.ServeHTTP(w, r)
	})
}

func isRouteAvailable(patterns [2]string, url string) (bool, map[string]string) {
	variables := make(map[string]string)
	for _, pattern := range patterns {
		regex := regexp.MustCompile(pattern)
		if matches := regex.FindStringSubmatch(url); len(matches) > 0 {
			for i, name := range regex.SubexpNames() {
				if i != 0 && name != "" {
					if split := strings.Split(name, "_"); len(split) == 4 && split[1] == "r" {
						min, err := strconv.Atoi(split[2])
						if err != nil {
							log.Fatal(err)
							return false, nil
						}

						max, err := strconv.Atoi(split[3])
						if err != nil {
							log.Fatal(err)
							return false, nil
						}

						if !validateRangeParams(matches[i], min, max) {
							return false, nil
						}

						variables[split[0]] = matches[i]
					} else {
						variables[name] = matches[i]
					}
				}
			}
			return true, variables
		}
	}
	return false, nil
}

func main() {

	flag.StringVar(&constant.Hostname, "hostname", "127.0.0.1", "Specify the hostname to listen to")
	flag.StringVar(&constant.Port, "port", "8080", "Specify the port to listen to")
	flag.IntVar(&constant.CacheControlMaxAge, "cache_control_max_age", 14, "Specify the max-age for cache-control header")
	flag.StringVar(&connection.Hostname, "mongo_hostname", "127.0.0.1", "Specify the hostname to connect to MongoDB instance")
	flag.StringVar(&connection.Port, "mongo_port", "27017", "Specify the port to connect to MongoDB instance")
	flag.Uint64Var(&connection.MaxPoolSize, "mongo_max_pool_size", 5, "Specify the max pool size for MongoDB connections")
	flag.Parse()

	var router = mux.NewRouter()
	router.Use(commonMiddleware)
	router.NewRoute().MatcherFunc(func(r *http.Request, rm *mux.RouteMatch) bool {
		patterns := [2]string{
			`/i/gdrive/(?P<quality_r_1_100>\d+)/(?P<width_r_0_5000>\d+)x(?P<height_r_0_5000>\d+)/(?P<option>[gts]{1,3})/(?P<path>.*)`,
			`/i/gdrive/(?P<quality_r_1_100>\d+)/(?P<width_r_0_5000>\d+)x(?P<height_r_0_5000>\d+)/(?P<path>.*)`,
		}

		available, vars := isRouteAvailable(patterns, r.URL.Path)
		if available {
			rm.Vars = vars
		}
		return available
	}).HandlerFunc(googledrive.GoogleDriveHandler)
	router.NewRoute().MatcherFunc(func(r *http.Request, rm *mux.RouteMatch) bool {
		patterns := [2]string{
			`/i/gridfs/(?P<quality_r_1_100>\d+)/(?P<width_r_0_5000>\d+)x(?P<height_r_0_5000>\d+)/(?P<option>[gts]{1,3})/(?P<path>.*)`,
			`/i/gridfs/(?P<quality_r_1_100>\d+)/(?P<width_r_0_5000>\d+)x(?P<height_r_0_5000>\d+)/(?P<path>.*)`,
		}

		available, vars := isRouteAvailable(patterns, r.URL.Path)
		if available {
			rm.Vars = vars
		}
		return available
	}).HandlerFunc(gridfs.GridFSHandler)

	srv := &http.Server{
		Handler:      router,
		Addr:         fmt.Sprintf("%s:%s", constant.Hostname, constant.Port),
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}
	log.Println(fmt.Sprintf("Server is started: %s", srv.Addr))
	log.Fatal(srv.ListenAndServe())

}
