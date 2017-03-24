package main

import (
	"flag"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"./api"
	"./core"

	"github.com/gorilla/mux"
)

var (
	dataDirectory     = flag.String("data-directory", "data/", "Path to data directory")
	httpListenAddress = flag.String("http", ":9999", "Http listen address (e.g. :9999)")
)

func main() {
	log.Printf("Service starting, data directory=%v, http=%v", *dataDirectory, *httpListenAddress)
	flag.Parse()

	router := core.NewRouter(*dataDirectory)
	defer router.Close()

	mux := mux.NewRouter()
	mux.HandleFunc("/health", api.GetHealth).Methods("GET")
	mux.HandleFunc("/streams/{id}", api.InsertMessage).Methods("POST")
	mux.HandleFunc("/streams/{id}", api.StreamMessage).Methods("GET")
	api.Initialize(router)
	http.Handle("/", mux)

	go func() {
		if err := http.ListenAndServe(*httpListenAddress, nil); err != nil {
			log.Fatal(err)
			os.Exit(-1)
		}
	}()

	exit := make(chan os.Signal, 1)
	signal.Notify(exit, syscall.SIGINT, syscall.SIGTERM)
	<-exit

	log.Printf("Exiting...")
}
