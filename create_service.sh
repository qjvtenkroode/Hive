#!/bin/bash

if [ $# -eq 0 ]
then
    echo "No arguments supplied"
    exit 0
elif [ $# -gt 1 ]
then
    echo "Too many arguments supplied"
    exit 0
elif [ -z "$1" ]
then
    echo "Service name cannot be empty"
    exit 0
fi

mkdir $1
touch $1/README.md
echo """# $1

## Endpoints

**Definition**
""" > $1/README.md

touch $1/Dockerfile
echo """FROM golang:latest AS builder

WORKDIR $GOPATH/src/hive/
COPY . . 
RUN go get -d .
RUN CGO_ENABLED=0 go build --tags netgo --ldflags \'-w -extldflags \"-static\"\' -a -o /go/app

FROM scratch

COPY --from=builder /go/app /
COPY --from=builder /go/src/hive/README.md /""" > $1/Dockerfile

touch $1/server.go
echo """package main

import (
	\"fmt\"
	\"io/ioutil\"
	\"net/http\"

	\"github.com/russross/blackfriday\"
)

// Server is a HTTP interface
type Server struct {
	http.Handler
}

// NewServer creates a new server with routing configured
func NewServer() *Server {
	s := new(Server)

	router := http.NewServeMux()
	router.HandleFunc(\"/\", http.HandlerFunc(s.handleIndex))

	s.Handler = router

	return s
}

func (s *Server) handleIndex(w http.ResponseWriter, r *http.Request) {
	url := r.URL.Path
	if url == \"/\" {
		w.Header().Set(\"content-type\", \"text/html\")
		readme, err := ioutil.ReadFile(\"README.md\")
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
		_, err = fmt.Fprintf(w, \"%s\", blackfriday.MarkdownCommon([]byte(readme)))
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		}
	} else {
		w.WriteHeader(http.StatusNotFound)
	}
}
""" > $1/server.go

touch $1/server_test.go
echo """package main

import (
	\"io\"
	\"io/ioutil\"
	\"net/http\"
	\"net/http/httptest\"
	\"testing\"

	\"github.com/russross/blackfriday\"
)

func TestAssets(t *testing.T) {
	readme, _ := ioutil.ReadFile(\"README.md\")
	tests := map[string]struct {
		url         string
		method      string
		postbody    io.Reader
		statuscode  int
		contenttype string
		body        string
	}{
		\"GET /\": {
			url:         \"/\",
			method:      http.MethodGet,
			postbody:    nil,
			statuscode:  http.StatusOK,
			contenttype: \"text/html\",
			body:        string(blackfriday.MarkdownCommon([]byte(readme))),
		},
	}

	server := NewServer()

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			request, _ := http.NewRequest(tc.method, tc.url, tc.postbody)
			response := httptest.NewRecorder()

			server.ServeHTTP(response, request)

			assertResponseCode(t, response.Code, tc.statuscode)
			assertContentType(t, response.Header().Get(\"content-type\"), tc.contenttype)
			assertBody(t, response.Body.String(), tc.body)
		})
	}
}

func assertBody(t *testing.T, got, want string) {
	t.Helper()
	if got != want {
		t.Errorf(\"wrong result body; got %s, want %s\", got, want)
	}
}

func assertResponseCode(t *testing.T, got, want int) {
	t.Helper()
	if got != want {
		t.Errorf(\"wrong response code; got %d, want %d\", got, want)
	}
}

func assertContentType(t *testing.T, got, want string) {
	t.Helper()
	if got != want {
		t.Errorf(\"wrong content-type; got %s, want %s\", got, want)
	}
}
""" > $1/server_test.go

touch $1/main.go
echo """package main

import (
	\"log\"
	\"net/http\"
)

func main() {
    server := NewServer()

	if err := http.ListenAndServe(\":80\", server); err != nil {
		log.Fatalf(\"could not listen on port 80 %v\", err)
	}
}
""" > $1/main.go
