package api

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"

	"../core"

	"github.com/gorilla/mux"
)

var (
	router            *core.Router
	heartbeatInterval = 30 * time.Second
	heartbeatMessage  = []byte("\r\n")
)

type (
	HealthResult struct {
		Success bool `json:"success"`
	}
)

func Initialize(router_ *core.Router) {
	router = router_
}

func InsertMessage(w http.ResponseWriter, r *http.Request) {
	var err error
	params := mux.Vars(r)
	id := params["id"]

	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, fmt.Sprintf("Cannot read body: %v", err), 500)
		return
	}
	defer r.Body.Close()

	if len(data) == 0 {
		http.Error(w, "Corrupt body", 400)
		return
	}

	stream, err := router.GetStream(id)
	if err != nil {
		http.Error(w, fmt.Sprintf("Cannot open stream: %v", err), 400)
		return
	}

	data = append(data, []byte("\r\n")...)
	stream.Write(data)
}

func StreamMessage(w http.ResponseWriter, r *http.Request) {
	var err error
	params := mux.Vars(r)
	id := params["id"]

	stream, err := router.GetStream(id)
	if err != nil {
		http.Error(w, fmt.Sprintf("Cannot open stream: %v", err), 400)
		return
	}

	offset := 1 + stream.LastOffset()
	offsetParameter := r.URL.Query().Get("offset")
	if len(offsetParameter) > 0 {
		offset, _ = strconv.ParseInt(offsetParameter, 10, 64)
	}

	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "Cannot create flusher", 500)
		return
	}

	scanner, err := stream.NewScanner(offset)
	if err != nil {
		http.Error(w, fmt.Sprintf("Cannot read log %v", err), 500)
		return
	}
	defer scanner.Close()

	for scanner.Scan() {
		if scanner.Err() != nil {
			http.Error(w, fmt.Sprintf("Cannot read log %v", scanner.Err()), 500)
			return
		}

		if _, err := w.Write(scanner.Bytes()); err != nil {
			return
		}
		flusher.Flush()
	}

	watcher := stream.NewWatcher()
	defer watcher.Close()
	modifications := watcher.Watch()

	heartbeat := time.NewTicker(heartbeatInterval)
	defer heartbeat.Stop()

	for {
		select {
		case <-modifications:
			for scanner.Scan() {
				if scanner.Err() != nil {
					http.Error(w, fmt.Sprintf("Cannot read log %v", scanner.Err()), 500)
					return
				}

				if _, err := w.Write(scanner.Bytes()); err != nil {
					return
				}
				flusher.Flush()
			}

		case <-heartbeat.C:
			if _, err := w.Write(heartbeatMessage); err != nil {
				return
			}
			flusher.Flush()
		}
	}
}

func GetHealth(w http.ResponseWriter, r *http.Request) {
	var url = "http://" + r.Host + "/streams/health"
	var message = []byte(time.Now().String() + "\r\n")

	go func() {
		time.Sleep(50 * time.Millisecond)
		http.Post(url, "text/plain", bytes.NewBuffer(message))
	}()

	client := http.Client{
		Timeout: time.Duration(1 * time.Second),
	}

	response, err := client.Get(url)
	if err != nil {
		body, _ := json.Marshal(&HealthResult{false})
		w.Header().Set("Content-Type", "application/json")
		w.Write(body)
		return
	}

	defer response.Body.Close()
	reader := bufio.NewReader(response.Body)
	result, _ := reader.ReadBytes('\n')
	ok := string(result) == string(message)

	body, _ := json.Marshal(&HealthResult{ok})
	w.Header().Set("Content-Type", "application/json")
	w.Write(body)
}
