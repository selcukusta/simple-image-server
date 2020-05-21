FROM golang:1.14-alpine as builder
COPY . $GOPATH/src/github.com/selcukusta/simple-image-server
WORKDIR $GOPATH/src/github.com/selcukusta/simple-image-server/cmd/image-server
RUN go get -d -v
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -a -installsuffix cgo -o $GOPATH/bin/simple-image-server .
COPY ./gcloud-image-server-cred.json $GOPATH/bin
FROM scratch as final
COPY --from=builder /go/bin/simple-image-server /go/bin/simple-image-server
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /go/bin/gcloud-image-server-cred.json /etc/
ENV GOOGLE_APPLICATION_CREDENTIALS=/etc/gcloud-image-server-cred.json
EXPOSE 8080
ENTRYPOINT ["/go/bin/simple-image-server"]
