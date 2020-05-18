module github.com/selcukusta/simple-image-server

replace github.com/selcukusta/simple-image-server/gdrive-handler => ./gdrive-handler

replace github.com/selcukusta/simple-image-server/image-processor => ./image-processor

replace github.com/selcukusta/simple-image-server/helper => ./helper

replace github.com/selcukusta/simple-image-server/constant => ./constant

go 1.14

require (
	github.com/gorilla/mux v1.7.4
	github.com/muesli/smartcrop v0.3.0
	github.com/namsral/flag v1.7.4-pre
	github.com/nfnt/resize v0.0.0-20180221191011-83c6a9932646
	golang.org/x/image v0.0.0-20200430140353-33d19683fad8 // indirect
	google.golang.org/api v0.24.0
)
