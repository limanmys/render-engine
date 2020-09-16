package web

import (
	"net/http"
	"renderer/src/constants"
	"renderer/src/sqlite"
	"renderer/src/ssh"
	"strconv"
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

	var val constants.Connection
	if val2, ok := constants.ActiveConnections[userID+serverID]; ok {
		val = val2
	} else {
		temp, res := ssh.CreateConnection(userID, serverID, server.IPAddress)
		if res == false {
			w.WriteHeader(403)
			_, _ = w.Write([]byte("Cannot create connection"))
			return
		}
		val.SSH = temp
	}
	if val.SFTP == nil {
		val.SFTP = ssh.OpenSFTPConnection(val.SSH)
	}

	constants.ActiveConnections[userID+serverID] = val

	flag := ssh.PutFile(val.SFTP, localPath, remotePath)
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

	var val constants.Connection
	if val2, ok := constants.ActiveConnections[userID+serverID]; ok {
		val = val2
	} else {
		temp, res := ssh.CreateConnection(userID, serverID, server.IPAddress)
		if res == false {
			w.WriteHeader(403)
			_, _ = w.Write([]byte("Cannot create connection"))
			return
		}
		val.SSH = temp
	}
	if val.SFTP == nil {
		val.SFTP = ssh.OpenSFTPConnection(val.SSH)
	}

	constants.ActiveConnections[userID+serverID] = val

	flag := ssh.GetFile(val.SFTP, localPath, remotePath)
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

	var val constants.Connection
	if val2, ok := constants.ActiveConnections[userID+serverID]; ok {
		val = val2
	} else {
		temp, res := ssh.CreateConnection(userID, serverID, server.IPAddress)
		if res == false {
			w.WriteHeader(403)
			_, _ = w.Write([]byte("Cannot create connection"))
			return
		}
		val.SSH = temp

		constants.ActiveConnections[userID+serverID] = val
	}

	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(200)
	_, _ = w.Write([]byte(ssh.RunCommand(val.SSH, command)))
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

	var val constants.Connection
	if val2, ok := constants.ActiveConnections[userID+serverID]; ok {
		val = val2
	} else {
		temp, res := ssh.CreateConnection(userID, serverID, server.IPAddress)
		if res == false {
			w.WriteHeader(403)
			_, _ = w.Write([]byte("Cannot create connection"))
			return
		}
		val.SSH = temp
	}

	constants.ActiveConnections[userID+serverID] = val
	port := ssh.CreateTunnel(remoteHost, remotePort, username, password)
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

	var val constants.Connection
	if val2, ok := constants.ActiveConnections[userID+remoteHost+username]; ok {
		val = val2
	} else {
		temp, res := ssh.RawCreateConnection(connectionType, username, password, remoteHost, remotePort)
		if res == false {
			w.WriteHeader(403)
			_, _ = w.Write([]byte("Cannot create connection"))
			return
		}
		val.SSH = temp

		constants.ActiveConnections[userID+remoteHost+username] = val
	}

	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(200)
	_, _ = w.Write([]byte(ssh.RunCommand(val.SSH, command)))
}
