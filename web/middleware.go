package web

import (
	"errors"
	"log"
	"net/http"
	"strings"

	"github.com/limanmys/render-engine/helpers"
	"github.com/limanmys/render-engine/postgresql"
)

func extractRequestData(target []string, r *http.Request) (map[string]string, error) {
	request := make(map[string]string)
	target = append(target, "token")
	for _, value := range target {
		temp := r.FormValue(value)
		if temp == "" {
			return nil, errors.New(value + " is missing.")
		}
		if value == "token" {
			//Try to get UserID from Token
			userID := postgresql.GetUserIDFromToken(temp)
			if userID == "" {
				return nil, errors.New("token is not valid")
			}
			request["user_id"] = userID
			continue
		}
		request[value] = temp
	}
	return request, nil
}

func permissionsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		token := r.FormValue("token")

		var userID, extensionID, serverID string

		if r.Header.Get("liman-token") != "" {
			userID = postgresql.GetUserIDFromLimanToken(r.Header.Get("liman-token"))
			if userID == "" {
				w.WriteHeader(403)
				_, _ = w.Write([]byte("nope4"))
				return
			}
		} else {
			userID = postgresql.GetUserIDFromToken(token)
			if userID == "" {
				w.WriteHeader(403)
				_, _ = w.Write([]byte("nope5"))
				return
			}
		}

		userObj := postgresql.GetUser(userID)
		if userObj.Status == 1 {
			next.ServeHTTP(w, r)
			return
		}

		if strings.TrimSpace(executeCommand("cat /liman/server/.env | grep 'LIMAN_RESTRICTED=true' >/dev/null && echo 1 || echo 0")) == "1" {
			next.ServeHTTP(w, r)
			return
		}

		permissions := postgresql.GetObjPermissions(userID)
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
				extensionID = postgresql.GetExtensionFromName(extensionID).ID
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
