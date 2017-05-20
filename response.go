package main

import (
	"fmt"
	"net/http"
)

func errorHandler(w http.ResponseWriter, r *http.Request, statusCode int) {
	w.WriteHeader(statusCode)
	if statusCode == 404 {
		fmt.Print(w, "Image asset not found on server.")
	}
}
