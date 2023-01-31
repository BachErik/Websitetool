package main

import (
	"log"
	"net/http"
)

func main() {
	initializeTemplates()
	initializeRoutes()
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}

func initializeRoutes() {
	http.HandleFunc("/test", testHandler)
}
