package main

import (
	"StK8s/apiserver"
	"fmt"
	"log"
	"net/http"
	"time"
)

var (
	address   = "127.0.0.1"
	port      = 8899
	apiPrefix = "/api/v1beta1"
)

func main() {
	var Rdata apiserver.RESTStorageData = "12"
	storage := map[string]apiserver.RESTStorage{
		"tasks": &Rdata,
	}

	s := &http.Server{
		Addr:           fmt.Sprintf("%s:%d", address, port),
		Handler:        apiserver.New(storage, apiPrefix),
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	log.Fatal(s.ListenAndServe())
}
