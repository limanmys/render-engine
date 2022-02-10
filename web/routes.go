package web

import (
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

// CreateWebServer Create Web Server
func CreateWebServer() {
	port := 5454
	log.Printf("Starting Server on %d\n", port)

	r := mux.NewRouter()
	r.HandleFunc("/", runExtensionHandler)
	r.HandleFunc("/sendLog", extensionLogHandler)
	r.HandleFunc("/backgroundJob", backgroundJobHandler)
	r.HandleFunc("/externalAPI", externalAPIHandler)
	r.HandleFunc("/runCommand", runCommandHandler)
	r.HandleFunc("/runScript", runScriptHandler)
	r.HandleFunc("/runOutsideCommand", runOutsideCommandHandler)
	r.HandleFunc("/putFile", putFileHandler)
	r.HandleFunc("/getFile", getFileHandler)
	r.HandleFunc("/openTunnel", openTunnelHandler)
	r.HandleFunc("/keepTunnelAlive", keepTunnelAliveHandler)
	r.HandleFunc("/verify", verifyHandler)
	r.HandleFunc("/setExtensionDb", setExtensionDb)

	r.Use(loggingMiddleware)
	r.Use(permissionsMiddleware)

	log.Fatal(http.ListenAndServeTLS("127.0.0.1:"+strconv.Itoa(port), "/liman/certs/liman.crt", "/liman/certs/liman.key", r))
}
