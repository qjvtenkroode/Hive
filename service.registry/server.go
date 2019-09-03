package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/russross/blackfriday"
)

// Server is a HTTP interface
type Server struct {
	Db Store
	http.Handler
}

// Asset defines an asset with all its unique fields
type Asset struct {
	Identifier string `json:"identifier"`
	Name       string `json:"name"`
	Type       string `json:"type"`
}

// Store is an interface for datastorage
type Store interface {
	getAsset(id string) (Asset, error)
	getAllAssets() []Asset
	storeAsset(id string, value Asset) error
	deleteAsset(id string) error
}

// NewServer creates a new server with routing configured
func NewServer(store Store) *Server {
	s := new(Server)

	router := http.NewServeMux()
	router.HandleFunc("/", http.HandlerFunc(s.handleIndex))
	router.HandleFunc("/assets/", http.HandlerFunc(s.handleAssets))

	s.Handler = router
	s.Db = store

	return s
}

func (s *Server) handleIndex(w http.ResponseWriter, r *http.Request) {
	url := r.URL.Path
	if url == "/" {
		w.Header().Set("content-type", "text/html")
		readme, err := ioutil.ReadFile("README.md")
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
		_, err = fmt.Fprintf(w, "%s", blackfriday.MarkdownCommon([]byte(readme)))
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		}
	} else {
		w.WriteHeader(http.StatusNotFound)
	}
}

func (s *Server) handleAssets(w http.ResponseWriter, r *http.Request) {
	trailing := r.URL.Path[len("/assets"):]
	id := strings.Split(trailing[1:], "/")[0]
	w.Header().Set("content-type", "application/json")
	switch r.Method {
	case http.MethodGet:
		s.handleAssetsGet(w, r, id)
	case http.MethodPost:
		s.handleAssetsPost(w, r, id)
	case http.MethodDelete:
		s.handleAssetsDelete(w, r, id)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func jsonify(w http.ResponseWriter, a Asset) ([]byte, error) {
	b, err := json.Marshal(a)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return nil, err
	}
	return b, err
}

func jsonifyArray(w http.ResponseWriter, a []Asset) ([]byte, error) {
	b, err := json.Marshal(a)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return nil, err
	}
	return b, err
}

func (s *Server) handleAssetsGet(w http.ResponseWriter, r *http.Request, id string) {
	if id == "" {
		assets := s.Db.getAllAssets()
		b, err := jsonifyArray(w, assets)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		_, _ = fmt.Fprint(w, string(b))
	} else {
		asset, err := s.Db.getAsset(id)
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			_, _ = fmt.Fprint(w, "{\"status\":\"error\",\"message\":\"asset not found\"}")
			return
		}
		b, err := jsonify(w, asset)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		_, _ = fmt.Fprint(w, string(b))
	}
}

func (s *Server) handleAssetsPost(w http.ResponseWriter, r *http.Request, id string) {
	var a Asset
	if r.Body == nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = fmt.Fprint(w, "{\"status\":\"error\",\"message\":\"oooops...\"}")
		return
	}
	err := json.NewDecoder(r.Body).Decode(&a)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = fmt.Fprint(w, "{\"status\":\"error\",\"message\":\"oooops...\"}")
		return
	}
	err = s.Db.storeAsset(a.Identifier, a)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = fmt.Fprint(w, "{\"status\":\"error\",\"message\":\"oooops...\"}")
		return
	}
	w.WriteHeader(http.StatusCreated)
	_, _ = fmt.Fprint(w, "{\"status\":\"ok\",\"message\":\"asset created\"}")
}

func (s *Server) handleAssetsDelete(w http.ResponseWriter, r *http.Request, id string) {
	err := s.Db.deleteAsset(id)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		_, _ = fmt.Fprint(w, "{\"status\":\"error\",\"message\":\"asset not found\"}")
	} else {
		w.WriteHeader(http.StatusNoContent)
		_, _ = fmt.Fprint(w, "{\"status\":\"ok\",\"message\":\"asset deleted succesfully\"}")
	}
}
