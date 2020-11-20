package web

import (
	"net/http"

	"github.com/limanmys/go/helpers"
)

var currentToken = ""

func dnsHandler(w http.ResponseWriter, r *http.Request) {
	target := []string{"server1", "server2", "server3"}
	request, err := extractRequestData(target, r)

	if err != nil {
		w.WriteHeader(403)
		_, _ = w.Write([]byte(err.Error()))
		return
	}

	result := helpers.SetDNSServers(request["server1"], request["server2"], request["server3"])
	if result == true {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("DNS updated!\n"))
	} else {
		w.WriteHeader(http.StatusNotAcceptable)
		_, _ = w.Write([]byte("DNS update failed!\n"))
	}
}

func fixExtensionKeyHandler(w http.ResponseWriter, r *http.Request) {
	target := []string{"extension_id"}
	request, err := extractRequestData(target, r)

	if err != nil {
		w.WriteHeader(403)
		_, _ = w.Write([]byte(err.Error()))
		return
	}

	result := helpers.FixExtensionKeys(request["extension_id"])

	if result == true {
		_, _ = w.Write([]byte("Key permissions updated!\n"))
		w.WriteHeader(http.StatusOK)
	} else {
		_, _ = w.Write([]byte("Key permission update failed!\n"))
		w.WriteHeader(http.StatusNotAcceptable)
	}
}

func userAddHandler(w http.ResponseWriter, r *http.Request) {
	target := []string{"extension_id"}
	request, err := extractRequestData(target, r)

	if err != nil {
		w.WriteHeader(403)
		_, _ = w.Write([]byte(err.Error()))
		return
	}

	result := helpers.AddUser(request["extension_id"])
	if result == true {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("New User Added!\n"))
	} else {
		w.WriteHeader(http.StatusNotAcceptable)
		_, _ = w.Write([]byte("User add failed!\n"))
	}
}

func userRemoveHandler(w http.ResponseWriter, r *http.Request) {
	target := []string{"extension_id"}
	request, err := extractRequestData(target, r)

	if err != nil {
		w.WriteHeader(403)
		_, _ = w.Write([]byte(err.Error()))
		return
	}

	result := helpers.RemoveUser(request["extension_id"])
	if result == true {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("User Removed!\n"))
	} else {
		w.WriteHeader(http.StatusNotAcceptable)
		_, _ = w.Write([]byte("User remove failed!\n"))
	}
}

func permissionFixHandler(w http.ResponseWriter, r *http.Request) {
	target := []string{"extension_id", "extension_name"}
	request, err := extractRequestData(target, r)

	if err != nil {
		w.WriteHeader(403)
		_, _ = w.Write([]byte(err.Error()))
		return
	}

	result := helpers.FixExtensionPermissions(request["extension_id"], request["extension_name"])
	result = helpers.FixExtensionKeys(request["extension_id"])
	if result == true {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("Permissions fixed!\n"))
	} else {
		w.WriteHeader(http.StatusNotAcceptable)
		_, _ = w.Write([]byte("Permission fix failed!\n"))
	}
}

func certificateAddHandler(w http.ResponseWriter, r *http.Request) {
	target := []string{"tmpPath", "targetName"}
	request, err := extractRequestData(target, r)

	if err != nil {
		w.WriteHeader(403)
		_, _ = w.Write([]byte(err.Error()))
		return
	}

	result := helpers.AddSystemCertificate(request["tmpPath"], request["targetName"])
	if result == true {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("New Certificate Added!\n"))
	} else {
		w.WriteHeader(http.StatusNotAcceptable)
		_, _ = w.Write([]byte("Certificate add failed!\n"))
	}
}

func certificateRemoveHandler(w http.ResponseWriter, r *http.Request) {
	target := []string{"targetName"}
	request, err := extractRequestData(target, r)

	if err != nil {
		w.WriteHeader(403)
		_, _ = w.Write([]byte(err.Error()))
		return
	}

	result := helpers.RemoveSystemCertificate(request["targetName"])
	if result == true {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("Certificate removed!\n"))
	} else {
		w.WriteHeader(http.StatusNotAcceptable)
		_, _ = w.Write([]byte("Certificate remove failed!\n"))
	}
}
