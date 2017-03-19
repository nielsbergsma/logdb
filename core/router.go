package core

import (
	"sync"
)

type Router struct {
	directory string
	lock      *sync.Mutex
	indexes   map[string]*Index
}

func NewRouter(directory string) *Router {
	return &Router{
		directory: directory,
		lock:      &sync.Mutex{},
		indexes:   map[string]*Index{},
	}
}

func (r *Router) Close() {
	for _, index := range r.indexes {
		index.Close()
	}
}

func (r *Router) GetIndex(name string) (*Index, error) {
	var err error

	if index, exist := r.indexes[name]; exist {
		return index, nil
	} else {
		r.lock.Lock()
		defer r.lock.Unlock()

		index, err = NewIndex(r.directory, name)
		if err != nil {
			return nil, err
		}

		r.indexes[name] = index
		return index, nil
	}
}
