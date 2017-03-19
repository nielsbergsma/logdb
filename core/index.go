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
	log     *biglog.BigLog
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
		log:     log,
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
			if written, err := index.log.Write(data); err != nil {
				log.Printf("Error during write %v", err)
			} else if written < len(data) {
				log.Printf("Error not all bytes written")
			}

		case <-index.dispose:
			index.log.Close()
			close(index.write)
			close(index.dispose)
		}
	}
}

func (index *Index) NewScanner(offset int64) (*biglog.Scanner, error) {
	return biglog.NewScanner(index.log, offset)
}

func (index *Index) NewWatcher() *biglog.Watcher {
	return biglog.NewWatcher(index.log)
}

func (index *Index) Close() {
	index.dispose <- struct{}{}
}

func (index *Index) Write(data []byte) {
	index.write <- data
}
