package web

import (
	"net/http"
	"renderer/src/connector"
	"renderer/src/constants"
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
	username, password, _, keyObject := sqlite.GetServerKey(userID, serverID)
	var val constants.Connection
	flag := false
	if val2, ok := constants.ActiveConnections[userID+serverID]; ok {
		val = val2
	} else if keyObject.Type == "ssh" || keyObject.Type == "ssh_certificate" {
		res := connector.CreateConnection(userID, serverID, server.IPAddress, &val)
		if res == false {
			w.WriteHeader(403)
			_, _ = w.Write([]byte("Cannot create connection"))
			return
		}
	}
	if keyObject.Type == "ssh" || keyObject.Type == "ssh_certificate" {
		if val.SFTP == nil {
			val.SFTP = connector.OpenSFTPConnection(val.SSH)
		}
		flag = connector.PutFileSFTP(val.SFTP, localPath, remotePath)
	} else if keyObject.Type == "winrm" {
		if val.SMB == nil {
			temp, err := connector.OpenSMBConnection(server.IPAddress, username, password)
			if err != nil {
				w.WriteHeader(403)
				_, _ = w.Write([]byte("Cannot create connection"))
				return
			}
			val.SMB = temp
		}
		flag = connector.PutFileSMB(val.SMB, localPath, remotePath)

	} else {
		w.WriteHeader(403)
		_, _ = w.Write([]byte("unsupported key type?"))
		return
	}
	val.LastConnection = time.Now()
	constants.ActiveConnections[userID+serverID] = val

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
	username, password, _, keyObject := sqlite.GetServerKey(userID, serverID)
	var val constants.Connection
	flag := false
	if val2, ok := constants.ActiveConnections[userID+serverID]; ok {
		val = val2
	} else if keyObject.Type == "ssh" || keyObject.Type == "ssh_certificate" {
		res := connector.CreateConnection(userID, serverID, server.IPAddress, &val)
		if res == false {
			w.WriteHeader(403)
			_, _ = w.Write([]byte("Cannot create connection"))
			return
		}
	}
	if keyObject.Type == "ssh" || keyObject.Type == "ssh_certificate" {
		if val.SFTP == nil {
			val.SFTP = connector.OpenSFTPConnection(val.SSH)
		}
		flag = connector.GetFileSFTP(val.SFTP, localPath, remotePath)
	} else if keyObject.Type == "winrm" {
		if val.SMB == nil {
			temp, err := connector.OpenSMBConnection(server.IPAddress, username, password)
			if err != nil {
				w.WriteHeader(403)
				_, _ = w.Write([]byte("Cannot create connection"))
				return
			}
			val.SMB = temp
		}
		flag = connector.GetFileSMB(val.SMB, localPath, remotePath)

	} else {
		w.WriteHeader(403)
		_, _ = w.Write([]byte("unsupported key type?"))
		return
	}
	val.LastConnection = time.Now()
	constants.ActiveConnections[userID+serverID] = val

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
		res := connector.CreateConnection(userID, serverID, server.IPAddress, &val)
		if res == false {
			w.WriteHeader(403)
			_, _ = w.Write([]byte("Cannot create connection"))
			return
		}
	}
	val.LastConnection = time.Now()
	constants.ActiveConnections[userID+serverID] = val

	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(200)
	_, _ = w.Write([]byte(connector.RunCommand(&val, command)))
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
		res := connector.CreateConnection(userID, serverID, server.IPAddress, &val)
		if res == false {
			w.WriteHeader(403)
			_, _ = w.Write([]byte("Cannot create connection"))
			return
		}
	}

	val.LastConnection = time.Now()
	constants.ActiveConnections[userID+serverID] = val

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

	var val constants.Connection
	if val2, ok := constants.ActiveConnections[userID+remoteHost+username]; ok {
		val = val2
	} else {
		res := connector.RawCreateConnection(&val, connectionType, username, password, remoteHost, remotePort)
		if res == false {
			w.WriteHeader(403)
			_, _ = w.Write([]byte("Cannot create connection"))
			return
		}
	}

	val.LastConnection = time.Now()
	constants.ActiveConnections[userID+remoteHost+username] = val
	output := connector.RunCommand(&val, command)
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(200)
	_, _ = w.Write([]byte(output))
}
