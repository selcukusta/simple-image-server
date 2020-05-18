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

	"github.com/gorilla/mux"
	constant "github.com/selcukusta/simple-image-server/backend/constant"
	gdrive "github.com/selcukusta/simple-image-server/backend/gdrive-handler"
	helper "github.com/selcukusta/simple-image-server/backend/helper"
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

func main() {
	if !helper.GoogleCredentialIsAvailable() {
		log.Fatal(`Google credential file cannot be found! Please create the file and set the "GOOGLE_APPLICATION_CREDENTIALS" environment variable.`)
	}

	flag.StringVar(&constant.Hostname, "hostname", "127.0.0.1", "Specify the hostname to listen to")
	flag.StringVar(&constant.Port, "port", "8080", "Specify the port to listen to")
	flag.IntVar(&constant.CacheControlMaxAge, "cache_control_max_age", 14, "Specify the max-age for cache-control header")
	flag.Parse()

	var router = mux.NewRouter()
	router.Use(commonMiddleware)
	router.NewRoute().MatcherFunc(func(r *http.Request, rm *mux.RouteMatch) bool {
		patterns := [2]string{
			`/i/(?P<slug>gdrive)/(?P<quality_r_1_100>\d+)/(?P<width_r_0_5000>\d+)x(?P<height_r_0_5000>\d+)/(?P<option>[gts]{1,3})/(?P<path>.*)`,
			`/i/(?P<slug>gdrive)/(?P<quality_r_1_100>\d+)/(?P<width_r_0_5000>\d+)x(?P<height_r_0_5000>\d+)/(?P<path>.*)`,
		}

		for _, pattern := range patterns {
			regex := regexp.MustCompile(pattern)
			if matches := regex.FindStringSubmatch(r.URL.Path); len(matches) > 0 {
				rm.Vars = make(map[string]string)
				for i, name := range regex.SubexpNames() {
					if i != 0 && name != "" {
						if split := strings.Split(name, "_"); len(split) == 4 && split[1] == "r" {
							min, err := strconv.Atoi(split[2])
							if err != nil {
								log.Fatal(err)
								return false
							}

							max, err := strconv.Atoi(split[3])
							if err != nil {
								log.Fatal(err)
								return false
							}

							if validateRangeParams(matches[i], min, max) == false {
								return false
							}

							rm.Vars[split[0]] = matches[i]
						} else {
							rm.Vars[name] = matches[i]
						}
					}
				}
				return true
			}
		}
		return false
	}).HandlerFunc(gdrive.GoogleDriveHandler)

	srv := &http.Server{
		Handler:      router,
		Addr:         fmt.Sprintf("%s:%s", constant.Hostname, constant.Port),
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}
	log.Println(fmt.Sprintf("Server is started: %s", srv.Addr))
	log.Fatal(srv.ListenAndServe())

}
