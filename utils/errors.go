package utils

import (
	"net/http"
)

//InternalServerError show internal server error
func InternalServerError(w http.ResponseWriter) {
	w.WriteHeader(http.StatusInternalServerError)
	w.Write([]byte("Internal server error"))
}
