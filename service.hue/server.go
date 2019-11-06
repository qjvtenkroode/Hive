package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/russross/blackfriday"
)

// Server is a HTTP interface
type Server struct {
	http.Handler
	Db     Store
	Bridge HueBridge
}

// HueBridge represents a Hue Bridge for REST API interactions
type HueBridge struct {
	Token   string
	Address string
}

// HueLightState is a representation of state
type HueLightState struct {
	Identifier string `json:"identifier"`
	State      string `json:"state"`
	Type       string `json:"type"`
	LastUpdate string `json:"last_update"`
}

// Store is an interface for datastorage
type Store interface {
	getState(id string) (HueLightState, error)
	storeState(id string, value HueLightState) error
}

// NewServer creates a new server with routing configured
func NewServer(store Store, h HueBridge) *Server {
	s := new(Server)

	router := http.NewServeMux()
	router.HandleFunc("/", http.HandlerFunc(s.handleIndex))
	router.HandleFunc("/state/", http.HandlerFunc(s.handleState))

	s.Handler = router
	s.Db = store
	s.Bridge = h

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

func jsonify(w http.ResponseWriter, h HueLightState) ([]byte, error) {
	b, err := json.Marshal(h)
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

func (s *Server) handleStatePost(w http.ResponseWriter, r *http.Request, id string) {
	if r.Body == nil {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = fmt.Fprintf(w, "{\"status\":\"error\",\"message\":\"no post body found\"}")
		return
	}
	var state HueLightState
	err := json.NewDecoder(r.Body).Decode(&state)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = fmt.Fprintf(w, "{\"status\":\"error\",\"message\":\"%s\"}", err)
		return
	}

	err = s.Db.storeState(state.Identifier, state)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = fmt.Fprintf(w, "{\"status\":\"error\",\"message\":\"%s\"}", err)
		return
	}

	// Set the state at the Hue Bridge
	payload := make(map[string]bool)
	if state.State == "on" {
		payload["on"] = true
	} else {
		payload["on"] = false
	}

	statePayload, err := json.Marshal(payload)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = fmt.Fprintf(w, "{\"status\":\"error\",\"message\":\"%s\"}", err)
		return
	}

	req, err := http.NewRequest("PUT", "http://"+s.Bridge.Address+"/api/"+s.Bridge.Token+"/lights/"+state.Identifier+"/state", bytes.NewBuffer(statePayload))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = fmt.Fprintf(w, "{\"status\":\"error\",\"message\":\"%s\"}", err)
		return
	}

	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	// @TODO add more checks for bridge failures
	_, err = client.Do(req)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = fmt.Fprintf(w, "{\"status\":\"error\",\"message\":\"%s\"}", err)
		return
	}

	b, err := jsonify(w, state)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = fmt.Fprintf(w, "{\"status\":\"error\",\"message\":\"%s\"}", err)
		return
	}
	w.WriteHeader(http.StatusOK)
	_, _ = fmt.Fprint(w, string(b))
}

func (s *Server) pollState() {
	resp, err := http.Get("http://" + s.Bridge.Address + "/api/" + s.Bridge.Token + "/lights/")
	if err != nil {
		fmt.Println("Error getting Hue Bridge data")
		fmt.Println(err)
		return
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error parsing the results body")
		fmt.Println(err)
		return
	}

	m := map[string]interface{}{}
	err = json.Unmarshal(body, &m)
	if err != nil {
		fmt.Println("Error converting JSON")
		fmt.Println(err)
		return
	}
	for key, value := range m {
		var lightState string
		value, _ := value.(map[string]interface{})
		state, _ := value["state"].(map[string]interface{})
		on, _ := state["on"].(bool)
		if on == true {
			lightState = "on"
		} else {
			lightState = "off"
		}
		modelid, _ := value["modelid"].(string)
		hueState := HueLightState{key, lightState, modelid, ""}

		err = s.Db.storeState(hueState.Identifier, hueState)
		if err != nil {
			fmt.Println("Error storing initial states")
			return
		}
	}
}
