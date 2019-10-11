package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"
	"strings"

	"github.com/eclipse/paho.mqtt.golang"
	"github.com/russross/blackfriday"
)

// Server is a HTTP interface
type Server struct {
	token  mqtt.Token
	client mqtt.Client
	Db     Store
	http.Handler
}

// ShellyState is a representation of a state
type ShellyState struct {
	Identifier string `json:"identifier"`
	State      string `json:"state"`
	Type       string `json:"type"`
}

// Store is an interface for datastorage
type Store interface {
	getState(id string) (ShellyState, error)
	storeState(id string, value ShellyState) error
}

// NewServer creates a new server with routing configured
func NewServer(store Store, client mqtt.Client, token mqtt.Token) *Server {
	s := new(Server)

	router := http.NewServeMux()
	router.HandleFunc("/", http.HandlerFunc(s.handleIndex))
	router.HandleFunc("/state/", http.HandlerFunc(s.handleState))

	s.Handler = router
	s.Db = store
	s.client = client
	s.token = token

	return s
}

func (s *Server) onMessageReceived(client mqtt.Client, message mqtt.Message) {
	r, _ := regexp.Compile("shellies/([^\\/]+)/relay/0")
	shelly := strings.Split(r.FindStringSubmatch(message.Topic())[1], "-")
	state := ShellyState{shelly[1], string(message.Payload()), shelly[0]}
	err := s.Db.storeState(shelly[1], state)
	if err != nil {
		fmt.Println("error: ", err)
	}
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

func (s *Server) handleState(w http.ResponseWriter, r *http.Request) {
	trailing := r.URL.Path[len("/state"):]
	id := strings.Split(trailing[1:], "/")[0]
	w.Header().Set("content-type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", r.Header.Get("Origin"))
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Origin, Accept")

	switch r.Method {
	case http.MethodGet:
		s.handleStateGet(w, r, id)
	case http.MethodPost:
		s.handleStatePost(w, r, id)
	case http.MethodOptions:
		w.WriteHeader(http.StatusOK)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func jsonify(w http.ResponseWriter, s ShellyState) ([]byte, error) {
	b, err := json.Marshal(s)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return nil, err
	}
	return b, err
}

func (s *Server) handleStateGet(w http.ResponseWriter, r *http.Request, id string) {
	state, err := s.Db.getState(id)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		_, _ = fmt.Fprint(w, "{\"status\":\"error\",\"message\":\"asset not found\"}")
		return
	}
	b, err := jsonify(w, state)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	_, _ = fmt.Fprint(w, string(b))
}

// just do a mqtt publish and let the state be catched by mqtt.
func (s *Server) handleStatePost(w http.ResponseWriter, r *http.Request, id string) {
	if r.Body == nil {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = fmt.Fprintf(w, "{\"status\":\"error\",\"message\":\"no post body found\"}")
		return
	}
	var state ShellyState
	err := json.NewDecoder(r.Body).Decode(&state)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = fmt.Fprintf(w, "{\"status\":\"error\",\"message\":\"%s\"}", err)
		return
	}
	topic := fmt.Sprintf("shellies/%s-%s/relay/0/command", state.Type, state.Identifier)
	s.client.Publish(topic, byte(0), false, state.State)

	b, err := jsonify(w, state)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	_, _ = fmt.Fprint(w, string(b))
}
