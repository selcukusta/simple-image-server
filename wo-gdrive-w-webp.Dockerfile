ARG APP_VERSION
FROM golang:1.15.0-buster as builder
COPY . $GOPATH/src/github.com/selcukusta/simple-image-server
WORKDIR $GOPATH/src/github.com/selcukusta/simple-image-server/cmd/image-server
RUN go get -d -v
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -a -installsuffix cgo -o $GOPATH/bin/simple-image-server .
FROM golang:1.15.0-alpine as final
ARG APP_VERSION
LABEL maintainer="selcukusta@gmail.com"
COPY --from=builder /go/bin/simple-image-server /go/bin/simple-image-server
RUN apk update && apk upgrade
RUN apk add --no-cache \
    libwebp-tools=1.1.0-r0
ENV APP_VERSION=${APP_VERSION}
ENV ABS_ACCOUNT_KEY=YOUR_ACCOUNT_KEY
ENV ABS_ACCOUNT_NAME=YOUR_ACCOUNT_NAME
ENV ABS_AZURE_URI=YOUR_AZURE_URI
ENV WEBP_ENABLED=1
EXPOSE 8080
ENTRYPOINT ["/go/bin/simple-image-server"]
