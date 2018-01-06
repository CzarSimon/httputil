package httputil

import (
	"log"
	"net/http"
)

const (
	Json      = "application/json"
	PlainText = "text/plain"
)

var (
	StatusOK = []byte(http.StatusText(http.StatusOK))
)

// SendOK Sends an OK status to the requestor.
func SendOK(w http.ResponseWriter) {
	w.WriteHeader(http.StatusOK)
	w.Write(StatusOK)
}

// Ping Responds to and logs a ping.
func Ping(w http.ResponseWriter, r *http.Request) {
	log.Println("Ping recieved")
	SendOK(w)
}

// SendErr Sends an error message and status to the requestor.
func SendErr(w http.ResponseWriter, err Error) {
	http.Error(w, err.Error(), err.Status)
}
