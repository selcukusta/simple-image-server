ARG APP_VERSION
FROM golang:1.15.0-buster as builder
COPY . $GOPATH/src/github.com/selcukusta/simple-image-server
WORKDIR $GOPATH/src/github.com/selcukusta/simple-image-server/cmd/image-server
RUN go get -d -v
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -a -installsuffix cgo -o $GOPATH/bin/simple-image-server .
COPY ./gcloud-image-server-cred.json $GOPATH/bin
FROM scratch as final
ARG APP_VERSION
COPY --from=builder /go/bin/simple-image-server /go/bin/simple-image-server
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /go/bin/gcloud-image-server-cred.json /etc/
ENV APP_VERSION=${APP_VERSION}
ENV GOOGLE_APPLICATION_CREDENTIALS=/etc/gcloud-image-server-cred.json
ENV ABS_ACCOUNT_KEY=YOUR_ACCOUNT_KEY
ENV ABS_ACCOUNT_NAME=YOUR_ACCOUNT_NAME
ENV ABS_AZURE_URI=YOUR_AZURE_URI
EXPOSE 8080
ENTRYPOINT ["/go/bin/simple-image-server"]
