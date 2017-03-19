package api

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"sync"
	"time"

	"../core"

	"github.com/gorilla/mux"
)

var (
	directory         string
	heartbeatInterval = 30 * time.Second
	heartbeatMessage  = []byte("\r\n")
	indexesLock       = &sync.Mutex{}
	indexes           = map[string]*core.Index{}
)

func Initialize(directory_ string) {
	directory = directory_
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
		http.Error(w, "Corrupted body", 400)
		return
	}

	index, exist := indexes[stream]
	if !exist {
		indexesLock.Lock()
		index, err = core.NewIndex(directory, stream)

		if err != nil {
			http.Error(w, fmt.Sprintf("Cannot open index: %v", err), 400)
			return
		}

		indexes[stream] = index
		indexesLock.Unlock()
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

	index, exist := indexes[stream]
	if !exist {
		indexesLock.Lock()
		index, err = core.NewIndex(directory, stream)

		if err != nil {
			http.Error(w, fmt.Sprintf("Cannot open index: %v", err), 400)
			return
		}

		indexes[stream] = index
		indexesLock.Unlock()
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

func Close() {
	for _, index := range indexes {
		index.Close()
	}
}
