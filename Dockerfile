FROM golang:1.15.3-alpine3.12

RUN apk update && apk upgrade

RUN apk --no-cache add bash wget curl git;

RUN mkdir /go/src/app

WORKDIR /go/src/app

RUN go get -u github.com/cosmtrek/air
