FROM golang:1.15.0-buster as builder
COPY . $GOPATH/src/github.com/selcukusta/db-seed
WORKDIR $GOPATH/src/github.com/selcukusta/db-seed/cmd/mongo-seed
RUN go get -d -v
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -a -installsuffix cgo -o $GOPATH/bin/db-seed .
CMD ["/go/bin/db-seed"]
