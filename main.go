package main

import (
	"fmt"
	"net/http"
)

func main() {
	http.HandleFunc("/", Welcome)

	fmt.Println("Server started on port 8080")

	http.ListenAndServe(":8080", nil)
}

func Welcome(w http.ResponseWriter, req *http.Request) {
	fmt.Fprintln(w, "Welcome to the Volunteer App!")
}
