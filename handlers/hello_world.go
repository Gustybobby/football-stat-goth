package handlers

import (
	"net/http"
)

func HandleHelloWorld(w http.ResponseWriter, r *http.Request) error {
	w.Write([]byte("Hello, World!"))
	return nil
}
