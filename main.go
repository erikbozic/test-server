package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

var (
	port            int
	chunkSize       int
	readerMaxMemory int
	serviceBaseUrl  string
	serviceCallPath string
)

func main() {
	flag.IntVar(&port, "port", 8888, "the port the application will listen on")
	flag.IntVar(&chunkSize, "chunkSize", 1000000, "size in bytes used to chunk the response in the download handler")
	flag.IntVar(&readerMaxMemory, "maxMemory", 50*1000*1000, "max size of memory used when reading files in the upload handler")
	flag.StringVar(&serviceBaseUrl, "serviceBaseUrl", "http://localhost:8888", "base address of the service called in the service handler")
	flag.StringVar(&serviceCallPath, "serviceCallPath", "/headers", "path to call in the service handler")
	flag.Parse()

	r := mux.NewRouter()
	r.HandleFunc("/upload", Upload).Methods(http.MethodPost)
	r.HandleFunc("/download", Download).Methods(http.MethodGet)
	r.HandleFunc("/headers", Headers).Methods(http.MethodGet)
	r.HandleFunc("/service", Service).Methods(http.MethodGet)
	log.Println("start listening on port:", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), r))
}
