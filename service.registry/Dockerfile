FROM golang:1.12.6-alpine

ARG DEV

RUN if [ "$DEV" == "1" ] ; then apk add git gcc musl-dev; else apk add git ; fi

WORKDIR /go/src/app
COPY . .

RUN go get -d -v ./...
RUN go install -v ./...
