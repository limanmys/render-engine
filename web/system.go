package web

import (
	"net/http"

	"github.com/limanmys/go/replications"
)

func dnsHandler(w http.ResponseWriter, r *http.Request) {
	replications.Dns()
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte("Message received!"))
}

func extensionHandler(w http.ResponseWriter, r *http.Request) {
	replications.Extension()
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte("Message received!"))
}

func certificateHandler(w http.ResponseWriter, r *http.Request) {
	replications.Certificate()
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte("Message received!"))
}
