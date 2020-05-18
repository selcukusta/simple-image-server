module github.com/selcukusta/simple-image-server/backend

replace github.com/selcukusta/simple-image-server/backend/gdrive-handler => ./gdrive-handler

replace github.com/selcukusta/simple-image-server/backend/image-processor => ./image-processor

replace github.com/selcukusta/simple-image-server/backend/helper => ./helper

replace github.com/selcukusta/simple-image-server/backend/constant => ./constant

go 1.14

require (
	github.com/Azure/azure-storage-blob-go v0.8.0
	github.com/Azure/go-autorest/autorest/adal v0.8.3 // indirect
	github.com/gorilla/mux v1.7.4
	github.com/muesli/smartcrop v0.3.0
	github.com/namsral/flag v1.7.4-pre
	github.com/nfnt/resize v0.0.0-20180221191011-83c6a9932646
	github.com/pkg/errors v0.9.1 // indirect
	golang.org/x/image v0.0.0-20200430140353-33d19683fad8 // indirect
	google.golang.org/api v0.24.0
)
