package web

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/limanmys/go/connector"
	tusd "github.com/tus/tusd/pkg/handler"
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

	sftpstore := connector.SftpStore{}
	composer := tusd.NewStoreComposer()
	sftpstore.UseIn(composer)
	tusdHandler, err := tusd.NewHandler(tusd.Config{
		BasePath:              "/sftpUpload/",
		StoreComposer:         composer,
		NotifyCompleteUploads: true,
	})
	if err != nil {
		panic(fmt.Errorf("unable to create tusd handler: %s", err))
	}

	go func() {
		for {
			event := <-tusdHandler.CompleteUploads
			fmt.Printf("Upload %s finished\n", event.Upload.ID)
		}
	}()

	r.PathPrefix("/sftpUpload/").Handler(http.StripPrefix("/sftpUpload/", tusdHandler))

	r.Use(loggingMiddleware)
	r.Use(permissionsMiddleware)

	log.Fatal(http.ListenAndServeTLS("127.0.0.1:"+strconv.Itoa(port), "/liman/certs/liman.crt", "/liman/certs/liman.key", r))
}
