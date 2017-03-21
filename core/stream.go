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

type Stream struct {
	log     *biglog.BigLog
	dispose chan struct{}
	write   chan []byte
}

func NewStream(directory, name string) (*Stream, error) {
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

	stream := &Stream{
		log:     log,
		dispose: make(chan struct{}),
		write:   make(chan []byte),
	}

	go stream.run()
	return stream, nil
}

func (stream *Stream) run() {
	for {
		select {
		case data := <-stream.write:
			if written, err := stream.log.Write(data); err != nil {
				log.Printf("Error during write %v", err)
			} else if written < len(data) {
				log.Printf("Error not all bytes written")
			}

		case <-stream.dispose:
			stream.log.Close()
			close(stream.write)
			close(stream.dispose)
		}
	}
}

func (stream *Stream) NewScanner(offset int64) (*biglog.Scanner, error) {
	return biglog.NewScanner(stream.log, offset)
}

func (stream *Stream) NewWatcher() *biglog.Watcher {
	return biglog.NewWatcher(stream.log)
}

func (stream *Stream) Close() {
	stream.dispose <- struct{}{}
}

func (stream *Stream) Write(data []byte) {
	stream.write <- data
}
