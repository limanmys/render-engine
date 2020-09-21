package web

import (
	"net/http"
	"renderer/src/connector"
	"renderer/src/sqlite"
	"strconv"
	"time"
)

func putFileHandler(w http.ResponseWriter, r *http.Request) {
	serverID := r.FormValue("server_id")
	if serverID == "" {
		w.WriteHeader(403)
		_, _ = w.Write([]byte("ServerID?"))
		return
	}

	server := sqlite.GetServer(serverID)

	token := r.FormValue("token")
	if token == "" {
		w.WriteHeader(403)
		_, _ = w.Write([]byte("Token?"))
		return
	}

	userID := sqlite.GetUserIDFromToken(token)
	if userID == "" {
		w.WriteHeader(403)
		_, _ = w.Write([]byte("Token?"))
		return
	}

	remotePath := r.FormValue("remotePath")
	if remotePath == "" {
		w.WriteHeader(403)
		_, _ = w.Write([]byte("remotePath?"))
		return
	}

	localPath := r.FormValue("localPath")
	if localPath == "" {
		w.WriteHeader(403)
		_, _ = w.Write([]byte("localPath?"))
		return
	}

	var val connector.Connection

	if val2, ok := connector.ActiveConnections[userID+serverID]; ok {
		val = val2
	}
	val.CreateFileConnection(userID, serverID, server.IPAddress)

	flag := val.Put(localPath, remotePath)

	val.LastConnection = time.Now()
	connector.ActiveConnections[userID+serverID] = val

	w.Header().Set("Content-Type", "text/plain")
	if flag == true {
		w.WriteHeader(200)
		_, _ = w.Write([]byte("ok"))
	} else {
		w.WriteHeader(201)
		_, _ = w.Write([]byte("nok"))
	}

}

func getFileHandler(w http.ResponseWriter, r *http.Request) {
	serverID := r.FormValue("server_id")
	if serverID == "" {
		w.WriteHeader(403)
		_, _ = w.Write([]byte("ServerID?"))
		return
	}

	server := sqlite.GetServer(serverID)

	token := r.FormValue("token")
	if token == "" {
		w.WriteHeader(403)
		_, _ = w.Write([]byte("Token?"))
		return
	}

	userID := sqlite.GetUserIDFromToken(token)
	if userID == "" {
		w.WriteHeader(403)
		_, _ = w.Write([]byte("Token?"))
		return
	}

	remotePath := r.FormValue("remotePath")
	if remotePath == "" {
		w.WriteHeader(403)
		_, _ = w.Write([]byte("remotePath?"))
		return
	}

	localPath := r.FormValue("localPath")
	if localPath == "" {
		w.WriteHeader(403)
		_, _ = w.Write([]byte("localPath?"))
		return
	}
	var val connector.Connection

	if val2, ok := connector.ActiveConnections[userID+serverID]; ok {
		val = val2
	}
	val.CreateFileConnection(userID, serverID, server.IPAddress)

	flag := val.Get(localPath, remotePath)

	val.LastConnection = time.Now()
	connector.ActiveConnections[userID+serverID] = val

	w.Header().Set("Content-Type", "text/plain")
	if flag == true {
		w.WriteHeader(200)
		_, _ = w.Write([]byte("ok"))
	} else {
		w.WriteHeader(201)
		_, _ = w.Write([]byte("nok"))
	}
}

func runCommandHandler(w http.ResponseWriter, r *http.Request) {
	serverID := r.FormValue("server_id")
	if serverID == "" {
		w.WriteHeader(403)
		_, _ = w.Write([]byte("ServerID?"))
		return
	}

	server := sqlite.GetServer(serverID)

	token := r.FormValue("token")
	if token == "" {
		w.WriteHeader(403)
		_, _ = w.Write([]byte("Token?"))
		return
	}

	command := r.FormValue("command")
	if command == "" {
		w.WriteHeader(403)
		_, _ = w.Write([]byte("command?"))
		return
	}

	userID := sqlite.GetUserIDFromToken(token)
	if userID == "" {
		w.WriteHeader(403)
		_, _ = w.Write([]byte("Token?"))
		return
	}

	var val connector.Connection
	if val2, ok := connector.ActiveConnections[userID+serverID]; ok {
		val = val2
	} else {
		res := val.CreateShell(userID, serverID, server.IPAddress)
		if res == false {
			w.WriteHeader(403)
			_, _ = w.Write([]byte("Cannot create connection"))
			return
		}
	}
	val.LastConnection = time.Now()
	connector.ActiveConnections[userID+serverID] = val

	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(200)
	_, _ = w.Write([]byte(val.Run(command)))
}

func openTunnelHandler(w http.ResponseWriter, r *http.Request) {
	serverID := r.FormValue("server_id")
	if serverID == "" {
		w.WriteHeader(403)
		_, _ = w.Write([]byte("ServerID?"))
		return
	}

	server := sqlite.GetServer(serverID)

	token := r.FormValue("token")
	if token == "" {
		w.WriteHeader(403)
		_, _ = w.Write([]byte("Token?"))
		return
	}

	userID := sqlite.GetUserIDFromToken(token)
	if userID == "" {
		w.WriteHeader(403)
		_, _ = w.Write([]byte("Token?"))
		return
	}

	remoteHost := r.FormValue("remote_host")
	if remoteHost == "" {
		w.WriteHeader(403)
		_, _ = w.Write([]byte("remoteHost?"))
		return
	}

	remotePort := r.FormValue("remote_port")
	if remotePort == "" {
		w.WriteHeader(403)
		_, _ = w.Write([]byte("remotePort?"))
		return
	}

	username := r.FormValue("username")
	if username == "" {
		w.WriteHeader(403)
		_, _ = w.Write([]byte("username?"))
		return
	}

	password := r.FormValue("password")
	if password == "" {
		w.WriteHeader(403)
		_, _ = w.Write([]byte("password?"))
		return
	}

	var val connector.Connection
	if val2, ok := connector.ActiveConnections[userID+serverID]; ok {
		val = val2
	} else {
		res := val.CreateShell(userID, serverID, server.IPAddress)
		if res == false {
			w.WriteHeader(403)
			_, _ = w.Write([]byte("Cannot create connection"))
			return
		}
	}

	val.LastConnection = time.Now()
	connector.ActiveConnections[userID+serverID] = val

	port := connector.CreateTunnel(remoteHost, remotePort, username, password)
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(200)
	_, _ = w.Write([]byte(strconv.Itoa(port)))
}

func runOutsideCommandHandler(w http.ResponseWriter, r *http.Request) {

	token := r.FormValue("token")
	if token == "" {
		w.WriteHeader(403)
		_, _ = w.Write([]byte("Token?"))
		return
	}

	command := r.FormValue("command")
	if command == "" {
		w.WriteHeader(403)
		_, _ = w.Write([]byte("command?"))
		return
	}

	userID := sqlite.GetUserIDFromToken(token)
	if userID == "" {
		w.WriteHeader(403)
		_, _ = w.Write([]byte("Token?"))
		return
	}

	connectionType := r.FormValue("connection_type")
	if connectionType == "" {
		w.WriteHeader(403)
		_, _ = w.Write([]byte("connectionType?"))
		return
	}

	remoteHost := r.FormValue("remote_host")
	if remoteHost == "" {
		w.WriteHeader(403)
		_, _ = w.Write([]byte("remoteHost?"))
		return
	}

	remotePort := r.FormValue("remote_port")
	if remotePort == "" {
		w.WriteHeader(403)
		_, _ = w.Write([]byte("remotePort?"))
		return
	}

	username := r.FormValue("username")
	if username == "" {
		w.WriteHeader(403)
		_, _ = w.Write([]byte("username?"))
		return
	}

	password := r.FormValue("password")
	if password == "" {
		w.WriteHeader(403)
		_, _ = w.Write([]byte("password?"))
		return
	}

	var val connector.Connection
	if val2, ok := connector.ActiveConnections[userID+remoteHost+username]; ok {
		val = val2
	} else {
		res := val.CreateShellRaw(connectionType, username, password, remoteHost, remotePort)
		if res == false {
			w.WriteHeader(403)
			_, _ = w.Write([]byte("Cannot create connection"))
			return
		}
	}

	val.LastConnection = time.Now()
	connector.ActiveConnections[userID+remoteHost+username] = val
	output := val.Run(command)
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(200)
	_, _ = w.Write([]byte(output))
}

func verifyHandler(w http.ResponseWriter, r *http.Request) {
	IPAddress := r.FormValue("ip_address")
	if IPAddress == "" {
		w.WriteHeader(403)
		_, _ = w.Write([]byte("IPAddress?"))
		return
	}

	username := r.FormValue("username")
	if username == "" {
		w.WriteHeader(403)
		_, _ = w.Write([]byte("username?"))
		return
	}

	password := r.FormValue("password")
	if password == "" {
		w.WriteHeader(403)
		_, _ = w.Write([]byte("password?"))
		return
	}

	port := r.FormValue("port")
	if port == "" {
		w.WriteHeader(403)
		_, _ = w.Write([]byte("port?"))
		return
	}

	keyType := r.FormValue("keyType")
	if keyType == "" {
		w.WriteHeader(403)
		_, _ = w.Write([]byte("keyType?"))
		return
	}

	flag := connector.VerifyAuth(username, password, IPAddress, port, keyType)

	w.Header().Set("Content-Type", "text/plain")
	if flag == true {
		w.WriteHeader(200)
		_, _ = w.Write([]byte("ok"))
	} else {
		w.WriteHeader(201)
		_, _ = w.Write([]byte("nok"))
	}
}

func runScriptHandler(w http.ResponseWriter, r *http.Request) {
	IPAddress := r.FormValue("ip_address")
	if IPAddress == "" {
		w.WriteHeader(403)
		_, _ = w.Write([]byte("IPAddress?"))
		return
	}

	username := r.FormValue("username")
	if username == "" {
		w.WriteHeader(403)
		_, _ = w.Write([]byte("username?"))
		return
	}

	password := r.FormValue("password")
	if password == "" {
		w.WriteHeader(403)
		_, _ = w.Write([]byte("password?"))
		return
	}

	port := r.FormValue("port")
	if port == "" {
		w.WriteHeader(403)
		_, _ = w.Write([]byte("port?"))
		return
	}

	keyType := r.FormValue("keyType")
	if keyType == "" {
		w.WriteHeader(403)
		_, _ = w.Write([]byte("keyType?"))
		return
	}

	flag := true

	w.Header().Set("Content-Type", "text/plain")
	if flag == true {
		w.WriteHeader(200)
		_, _ = w.Write([]byte("ok"))
	} else {
		w.WriteHeader(201)
		_, _ = w.Write([]byte("nok"))
	}
}
