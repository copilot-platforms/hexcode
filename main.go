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

	assetsPath := "./web/dist"
	isLocal := os.Getenv("IS_LOCAL")
	if isLocal == "" {
		assetsPath = "/app/web/dist"
	}

	http.Handle("/", http.FileServer(http.Dir(assetsPath)))

	log.Println("listening on", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
