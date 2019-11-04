package main

import (
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/russross/blackfriday"
)

func TestAssets(t *testing.T) {
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
		"GET /state/2": {
			url:         "/",
			method:      http.MethodGet,
			postbody:    nil,
			statuscode:  http.StatusOK,
			contenttype: "application/json",
			body:        "{\"identifier\":\"2\",\"state\":\"on\",\"type\":\"test\",\"lat_update\":\"\"}"
		},
		"GET /state/notfound": {
			url:         "/state/notfound/",
			method:      http.MethodGet,
			postbody:    nil,
			statuscode:  http.StatusNotFound,
			contenttype: "application/json",
			body:        "{\"status\":\"error\",\"message\":\"asset not found\"}",
		},


	}

	bridge := new(HueBridge)
	server := NewServer(seedStore(), *bridge)

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			request, _ := http.NewRequest(tc.method, tc.url, tc.postbody)
			response := httptest.NewRecorder()

			server.ServeHTTP(response, request)

			assertResponseCode(t, response.Code, tc.statuscode)
			assertContentType(t, response.Header().Get("content-type"), tc.contenttype)
			assertBody(t, response.Body.String(), tc.body)
		})
	}
}

func seedStore() Store {
	store := new(InMemoryStore)
	seed := make(map[string]HueLightState)
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
