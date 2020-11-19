package web

import (
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/limanmys/go/helpers"
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

	r.Use(loggingMiddleware)
	r.Use(permissionsMiddleware)

	log.Fatal(http.ListenAndServeTLS("127.0.0.1:"+strconv.Itoa(port), helpers.CertsPath+"liman.crt", helpers.CertsPath+"liman.key", r))
}
