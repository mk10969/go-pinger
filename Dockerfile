FROM golang:1.15

WORKDIR /go/src/app

COPY . .

ENV GO111MODULE=on

RUN go build

ENTRYPOINT [ "go-pinger" ]
