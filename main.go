package main

import (
	"log"
	"net/http"
)

func main() {
	mys := newServer()
	s := &http.Server{
		Handler: mys,
		Addr:    ":8080",
	}
	log.Print("Server is starting")
	log.Fatal(s.ListenAndServe())
}
