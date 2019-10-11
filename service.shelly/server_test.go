package main

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/russross/blackfriday"
)

func TestShelly(t *testing.T) {
	readme, _ := ioutil.ReadFile("README.md")
	tests := map[string]struct {
		url         string
		method      string
		postbody    io.Reader
		statuscode  int
		contenttype string
		body        string
	}{
		"GET /": {
			url:         "/",
			method:      http.MethodGet,
			postbody:    nil,
			statuscode:  http.StatusOK,
			contenttype: "text/html",
			body:        string(blackfriday.MarkdownCommon([]byte(readme))),
		},
		"GET /state/2C0A55": {
			url:         "/state/2C0A55/",
			method:      http.MethodGet,
			postbody:    nil,
			statuscode:  http.StatusOK,
			contenttype: "application/json",
			body:        "{\"identifier\":\"2C0A55\",\"state\":\"on\",\"type\":\"shelly1\"}",
		},
		"GET /state/notfound": {
			url:         "/state/notfound/",
			method:      http.MethodGet,
			postbody:    nil,
			statuscode:  http.StatusNotFound,
			contenttype: "application/json",
			body:        "{\"status\":\"error\",\"message\":\"asset not found\"}",
		},
		"POST /state/2C0A55": {
			url:         "/state/2C0A55/",
			method:      http.MethodPost,
			postbody:    bytes.NewBuffer([]byte("{\"identifier\":\"2C0A55\",\"state\":\"off\",\"type\":\"shelly1\"}")),
			statuscode:  http.StatusOK,
			contenttype: "application/json",
			body:        "{\"identifier\":\"2C0A55\",\"state\":\"off\",\"type\":\"shelly1\"}",
		},
		"POST /state/notfound": {
			url:         "/state/notfound/",
			method:      http.MethodPost,
			postbody:    bytes.NewBuffer([]byte("{\"identifier\":\"notfound\",\"state\":\"off\",\"type\":\"shelly1\"}")),
			statuscode:  http.StatusOK,
			contenttype: "application/json",
			body:        "{\"identifier\":\"notfound\",\"state\":\"off\",\"type\":\"shelly1\"}",
		},
		"POST /state/nobody": {
			url:         "/state/nobody/",
			method:      http.MethodPost,
			postbody:    nil,
			statuscode:  http.StatusBadRequest,
			contenttype: "application/json",
			body:        "{\"status\":\"error\",\"message\":\"no post body found\"}",
		},
	}

	client := newMockClient()
	token := client.Connect()
	server := NewServer(seedStore(), client, token)

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			request, _ := http.NewRequest(tc.method, tc.url, tc.postbody)
			response := httptest.NewRecorder()

			server.ServeHTTP(response, request)

			assertResponseCode(t, response.Code, tc.statuscode)
			assertContentType(t, response.Header().Get("content-type"), tc.contenttype)
			assertBody(t, response.Body.String(), tc.body)
			if tc.method == http.MethodDelete || tc.method == http.MethodPost {
				fmt.Println("Reseed the store")
				server = NewServer(seedStore(), client, token)
			}
		})
	}
}

func seedStore() Store {
	store := new(InMemoryStore)
	seed := make(map[string]ShellyState)
	seed["2C0A55"] = ShellyState{Identifier: "2C0A55", State: "on", Type: "shelly1"}
	store.States = seed
	return store
}

func assertBody(t *testing.T, got, want string) {
	t.Helper()
	if got != want {
		t.Errorf("wrong result body; got %s, want %s", got, want)
	}
}

func assertResponseCode(t *testing.T, got, want int) {
	t.Helper()
	if got != want {
		t.Errorf("wrong response code; got %d, want %d", got, want)
	}
}

func assertContentType(t *testing.T, got, want string) {
	t.Helper()
	if got != want {
		t.Errorf("wrong content-type; got %s, want %s", got, want)
	}
}
