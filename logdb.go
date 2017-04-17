package main

import (
	"flag"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"./api"
	"./core"

	"github.com/gorilla/mux"
	"google.golang.org/grpc"
)

var (
	dataDirectory     = flag.String("data-directory", "data/", "Path to data directory")
	httpListenAddress = flag.String("http", ":9999", "Http listen address (e.g. :9999)")
	grpcListenAddress = flag.String("grpc", ":9998", "GRPC listen address (e.g. :9998)")
)

func main() {
	log.Printf("Service starting, data directory=%v, http=%v, grpc=%v", *dataDirectory, *httpListenAddress, *grpcListenAddress)
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

	go func() {
		grpcListen, err := net.Listen("tcp", *grpcListenAddress)
		if err != nil {
			log.Fatal(err)
			os.Exit(-1)
		}

		grpcServer := grpc.NewServer()
		api.RegisterStreamApiServer(grpcServer, &api.GrpcApi{})
		grpcServer.Serve(grpcListen)
	}()

	exit := make(chan os.Signal, 1)
	signal.Notify(exit, syscall.SIGINT, syscall.SIGTERM)
	<-exit

	log.Printf("Exiting...")
}
