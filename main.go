package main

import (
	"log"
	"net/http"
	"os"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	http.Handle("/", http.FileServer(http.Dir("./web/dist")))

	log.Println("listening on", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
