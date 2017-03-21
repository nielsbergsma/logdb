package core

import (
	"sync"
)

type Router struct {
	directory string
	lock      *sync.Mutex
	streams   map[string]*Stream
}

func NewRouter(directory string) *Router {
	return &Router{
		directory: directory,
		lock:      &sync.Mutex{},
		streams:   map[string]*Stream{},
	}
}

func (r *Router) Close() {
	for _, stream := range r.streams {
		stream.Close()
	}
}

func (r *Router) GetStream(name string) (*Stream, error) {
	var err error

	if stream, exist := r.streams[name]; exist {
		return stream, nil
	} else {
		r.lock.Lock()
		defer r.lock.Unlock()

		stream, err = NewStream(r.directory, name)
		if err != nil {
			return nil, err
		}

		r.streams[name] = stream
		return stream, nil
	}
}
