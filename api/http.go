package api

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"

	"../core"

	"github.com/gorilla/mux"
)

var (
	directory         string
	router            *core.Router
	heartbeatInterval = 30 * time.Second
	heartbeatMessage  = []byte("\r\n")
)

func Initialize(router_ *core.Router) {
	router = router_
}

func InsertMessage(w http.ResponseWriter, r *http.Request) {
	var err error
	params := mux.Vars(r)
	stream := params["stream"]

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

	index, err := router.GetIndex(stream)
	if err != nil {
		http.Error(w, fmt.Sprintf("Cannot open index: %v", err), 400)
		return
	}

	data = append(data, []byte("\r\n")...)
	index.Write(data)
}

func StreamMessage(w http.ResponseWriter, r *http.Request) {
	var err error
	params := mux.Vars(r)
	stream := params["stream"]

	offset := int64(0)
	offsetParameter := r.URL.Query().Get("offset")
	if len(offsetParameter) > 0 {
		offset, _ = strconv.ParseInt(offsetParameter, 10, 64)
	}

	index, err := router.GetIndex(stream)
	if err != nil {
		http.Error(w, fmt.Sprintf("Cannot open index: %v", err), 400)
		return
	}

	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "Cannot create flusher", 500)
		return
	}

	scanner, err := index.NewScanner(offset)
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

	watcher := index.NewWatcher()
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
