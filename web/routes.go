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
	r.HandleFunc("/dns", dnsHandler)
	r.HandleFunc("/userAdd", userAddHandler)
	r.HandleFunc("/userRemove", userRemoveHandler)
	r.HandleFunc("/fixPermissions", permissionFixHandler)
	r.HandleFunc("/certificateAdd", certificateAddHandler)
	r.HandleFunc("/certificateRemove", certificateRemoveHandler)
	r.HandleFunc("/fixExtensionKeysPermission", fixExtensionKeyHandler)

	r.Use(loggingMiddleware)
	r.Use(permissionsMiddleware)

	targetHost := "127.0.0.1"

	if !helpers.ListenInternally {
		targetHost = "0.0.0.0"
	}
	log.Printf("Starting Server on %v:%v\n", targetHost, port)
	log.Fatal(http.ListenAndServeTLS(targetHost+":"+strconv.Itoa(port), helpers.CertsPath+"liman.crt", helpers.CertsPath+"liman.key", r))
}
