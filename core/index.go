package core

import (
	"log"
	"os"
	"path"

	"github.com/ninibe/netlog/biglog"
)

const (
	maxIndexEntries = 10000
)

type Index struct {
	Log     *biglog.BigLog
	dispose chan struct{}
	write   chan []byte
}

func NewIndex(directory, name string) (*Index, error) {
	var log *biglog.BigLog
	var err error

	filePath := path.Join(directory, name)
	if _, err = os.Stat(filePath); os.IsNotExist(err) {
		log, err = biglog.Create(filePath, maxIndexEntries)
	} else {
		log, err = biglog.Open(filePath)
	}

	if err != nil {
		return nil, err
	}

	index := &Index{
		Log:     log,
		dispose: make(chan struct{}),
		write:   make(chan []byte),
	}

	go index.run()
	return index, nil
}

func (index *Index) run() {
	for {
		select {
		case data := <-index.write:
			if written, err := index.Log.Write(data); err != nil {
				log.Printf("Error during write %v", err)
			} else if written < len(data) {
				log.Printf("Error not all bytes written")
			}

		case <-index.dispose:
			index.Log.Close()
			close(index.write)
			close(index.dispose)
		}
	}
}

func (index *Index) Close() {
	index.dispose <- struct{}{}
}

func (index *Index) Write(data []byte) {
	index.write <- data
}
