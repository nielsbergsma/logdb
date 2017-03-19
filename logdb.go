package main

import (
	"flag"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"./api"

	"github.com/gorilla/mux"
)

var (
	dataDirectory     = flag.String("data-directory", "data/", "Path to data directory")
	httpListenAddress = flag.String("http", ":9999", "Http listen address (e.g. :9999)")
)

func main() {
	log.Printf("Service starting, data directory=%v, http=%v", *dataDirectory, *httpListenAddress)
	flag.Parse()

	api.Initialize(*dataDirectory)

	mux := mux.NewRouter()
	mux.HandleFunc("/streams/{stream}", api.InsertMessage).Methods("POST")
	mux.HandleFunc("/streams/{stream}", api.StreamMessage).Methods("GET")
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
	api.Close()
}
