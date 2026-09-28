package main

import (
	"fmt"
	"net/http"
)

func ping(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "pong")
}

func main() {
	http.HandleFunc("/ping", ping)

	fmt.Println("Gateway service started on port 8080")

	http.ListenAndServe(":8080", nil)
}