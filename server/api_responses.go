package main

import "net/http"

var noBody = []byte{}

func sendAPIResponse(w http.ResponseWriter, statusCode int, body []byte) {
	w.WriteHeader(statusCode)
	w.Write(body)
}

// sendBadRequest sends a bad request to the client
func sendBadRequest(w http.ResponseWriter) {
	sendAPIResponse(w, http.StatusBadRequest, noBody)
}

// sendInternalServerError sends a internal server error responseto the client
func sendInternalServerError(w http.ResponseWriter) {
	sendAPIResponse(w, http.StatusInternalServerError, noBody)
}
