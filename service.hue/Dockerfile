FROM golang:latest AS builder

WORKDIR $GOPATH/src/hive/
COPY . . 
RUN go get -d .
RUN CGO_ENABLED=0 go build --tags netgo --ldflags '-w -extldflags "-static"' -a -o /go/app

FROM scratch

COPY --from=builder /go/app /
COPY --from=builder /go/src/hive/README.md /
