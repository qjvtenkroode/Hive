package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/russross/blackfriday"
)

func TestAssets(t *testing.T) {
	readme, _ := ioutil.ReadFile("README.md")
	posttest, _ := json.Marshal(map[string]string{"identifier": "posttest", "name": "Switch post", "type": "Shelly 1"})
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
		"GET /assets": {
			url:         "/assets/",
			method:      http.MethodGet,
			postbody:    nil,
			statuscode:  http.StatusOK,
			contenttype: "application/json",
			body:        "[{\"identifier\":\"deletetest\",\"name\":\"Switch somewhere\",\"type\":\"Shelly 1\"},{\"identifier\":\"test\",\"name\":\"Switch bedroom\",\"type\":\"Shelly 1\"}]",
		},
		"GET /assets/test": {
			url:         "/assets/test",
			method:      http.MethodGet,
			postbody:    nil,
			statuscode:  http.StatusOK,
			contenttype: "application/json",
			body:        "{\"identifier\":\"test\",\"name\":\"Switch bedroom\",\"type\":\"Shelly 1\"}",
		},
		"GET /assets/notfound": {
			url:         "/assets/notfound",
			method:      http.MethodGet,
			postbody:    nil,
			statuscode:  http.StatusNotFound,
			contenttype: "application/json",
			body:        "{\"status\":\"error\",\"message\":\"asset not found\"}",
		},
		"POST /assets": {
			url:         "/assets/",
			method:      http.MethodPost,
			postbody:    bytes.NewBuffer(posttest),
			statuscode:  http.StatusCreated,
			contenttype: "application/json",
			body:        "{\"status\":\"ok\",\"message\":\"asset created\"}",
		},
		"DELETE /assets/deletetest": {
			url:         "/assets/deletetest",
			method:      http.MethodDelete,
			postbody:    nil,
			statuscode:  http.StatusNoContent,
			contenttype: "application/json",
			body:        "{\"status\":\"ok\",\"message\":\"asset deleted succesfully\"}",
		},
		"DELETE /assets/notfound": {
			url:         "/assets/notfound",
			method:      http.MethodDelete,
			postbody:    nil,
			statuscode:  http.StatusNotFound,
			contenttype: "application/json",
			body:        "{\"status\":\"error\",\"message\":\"asset not found\"}",
		},
	}

	server := NewServer(seedStore())

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
				server = NewServer(seedStore())
			}
		})
	}
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

func seedStore() Store {
	store := new(InMemoryStore)
	seed := make(map[string]Asset)
	seed["test"] = Asset{Identifier: "test", Name: "Switch bedroom", Type: "Shelly 1"}
	seed["deletetest"] = Asset{Identifier: "deletetest", Name: "Switch somewhere", Type: "Shelly 1"}
	store.Assets = seed
	return store
}
