package web

import (
	"log"
	"net/http"
	"renderer/src/helpers"
	"renderer/src/sqlite"
	"strconv"
	"strings"

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
	r.Use(loggingMiddleware)
	r.Use(permissionsMiddleware)
	log.Fatal(http.ListenAndServe("127.0.0.1:"+strconv.Itoa(port), r))
}

func permissionsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		token := r.FormValue("token")

		var userID, extensionID, serverID string

		if r.Header.Get("liman-token") != "" {
			userID = sqlite.GetUserIDFromLimanToken(r.Header.Get("liman-token"))
			if userID == "" {
				w.WriteHeader(403)
				_, _ = w.Write([]byte("nope4"))
				return
			}
		} else {
			userID = sqlite.GetUserIDFromToken(token)
			if userID == "" {
				w.WriteHeader(403)
				_, _ = w.Write([]byte("nope5"))
				return
			}
		}

		userObj := sqlite.GetUser(userID)
		if userObj.Status == 1 {
			next.ServeHTTP(w, r)
			return
		}

		if strings.TrimSpace(executeCommand("cat /liman/server/.env | grep 'LIMAN_RESTRICTED=true' >/dev/null && echo 1 || echo 0")) == "1" {
			next.ServeHTTP(w, r)
			return
		}

		permissions := sqlite.GetObjPermissions(userID)
		if r.FormValue("server_id") != "" {
			serverID = r.FormValue("server_id")
			if !helpers.Contains(permissions, serverID) {
				w.WriteHeader(403)
				_, _ = w.Write([]byte("nope6"))
				return
			}
		}

		if r.FormValue("extension_id") != "" {
			extensionID = r.FormValue("extension_id")
			if helpers.IsValidUUID(extensionID) == false {
				extensionID = sqlite.GetExtensionFromName(extensionID).ID
			}
			if !helpers.Contains(permissions, extensionID) || extensionID == "" {
				w.WriteHeader(403)
				_, _ = w.Write([]byte("nope7"))
				return
			}
		}
		next.ServeHTTP(w, r)
	})
}

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println(r.URL.Path)
		next.ServeHTTP(w, r)
	})
}
