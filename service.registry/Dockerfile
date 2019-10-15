FROM golang:1.12.6-alpine

ARG DEV

RUN if [ "$DEV" == "1" ] ; then apk add git gcc musl-dev tzdata; else apk add git tzdata; fi && cp /usr/share/zoneinfo/Europe/Amsterdam /etc/localtime && apk del tzdata

WORKDIR /go/src/app
COPY . .

RUN go get -d -v ./...
RUN go install -v ./...
