package httputil

import (
	"encoding/json"
	"log"
	"net/http"
)

// Content type names
const (
	JSON       = "application/json"
	PLAIN_TEXT = "text/plain"
)

var (
	StatusOK = []byte(http.StatusText(http.StatusOK))
)

// SendOK Sends an OK status to the requestor.
func SendOK(w http.ResponseWriter) {
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", PLAIN_TEXT)
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

// SendJSON Marshals a json body and sends as response.
func SendJSON(w http.ResponseWriter, v interface{}) error {
	js, err := json.Marshal(v)
	if err != nil {
		log.Println(err)
		return InternalServerError
	}
	w.Header().Set("Content-Type", JSON)
	w.WriteHeader(http.StatusOK)
	w.Write(js)
	return nil
}
