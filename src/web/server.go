package web

import (
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

var (
	port     = 5454
	certPath = "/liman/certs/liman.crt"
	keyPath  = "/liman/certs/liman.key"
)

// CreateWebServer Create Web Server
func CreateWebServer() {
	log.Printf("Starting Server on %d\n", port)

	r := mux.NewRouter()
	r.HandleFunc("/", runExtensionHandler)
	r.Use(loggingMiddleware)
	log.Fatal(http.ListenAndServe("127.0.0.1:"+strconv.Itoa(port), r))
}

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println(r.URL.Path)
		next.ServeHTTP(w, r)
	})
}
