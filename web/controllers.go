package web

import (
	"encoding/base64"
	"fmt"
	"io"
	"log"
	"net/http"
	"path/filepath"
	"strconv"
	"time"

	"github.com/gorilla/websocket"
	"github.com/limanmys/go/connector"
	"github.com/limanmys/go/postgresql"
)

func putFileHandler(w http.ResponseWriter, r *http.Request) {
	target := []string{"server_id", "remote_path", "local_path"}
	request, err := extractRequestData(target, r)

	if err != nil {
		w.WriteHeader(403)
		_, _ = w.Write([]byte(err.Error()))
		return
	}

	server := postgresql.GetServer(request["server_id"])

	val, err := connector.GetConnection(request["user_id"], request["server_id"], server.IPAddress)

	val.CreateFileConnection(request["user_id"], request["server_id"], server.IPAddress)

	var remotePath string

	if server.Os == "linux" {
		remotePath = "/tmp/" + filepath.Base(request["remote_path"])
	} else {
		remotePath = val.WindowsPath + request["remote_path"]
	}

	flag := val.Put(request["local_path"], remotePath)

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
	target := []string{"server_id", "remote_path", "local_path"}
	request, err := extractRequestData(target, r)

	if err != nil {
		w.WriteHeader(403)
		_, _ = w.Write([]byte(err.Error()))
		return
	}

	server := postgresql.GetServer(request["server_id"])

	val, err := connector.GetConnection(request["user_id"], request["server_id"], server.IPAddress)

	val.CreateFileConnection(request["user_id"], request["server_id"], server.IPAddress)

	var remotePath string

	if server.Os == "linux" {
		remotePath = "/tmp/" + filepath.Base(request["remote_path"])
	} else {
		remotePath = val.WindowsPath + request["remote_path"]
	}

	flag := val.Get(request["local_path"], remotePath)

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
	target := []string{"server_id", "command"}
	request, err := extractRequestData(target, r)

	if err != nil {
		w.WriteHeader(403)
		_, _ = w.Write([]byte(err.Error()))
		return
	}

	server := postgresql.GetServer(request["server_id"])

	val, err := connector.GetConnection(request["user_id"], request["server_id"], server.IPAddress)

	if err != nil {
		w.WriteHeader(403)
		_, _ = w.Write([]byte(err.Error()))
		return
	}

	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(200)
	_, _ = w.Write([]byte(val.Run(request["command"])))
}

func openTunnelHandler(w http.ResponseWriter, r *http.Request) {
	target := []string{"remote_host", "remote_port", "username", "password"}
	request, err := extractRequestData(target, r)

	if err != nil {
		w.WriteHeader(403)
		_, _ = w.Write([]byte(err.Error()))
		return
	}

	port := connector.CreateTunnel(request["remote_host"], request["remote_port"], request["username"], request["password"])
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(200)
	_, _ = w.Write([]byte(strconv.Itoa(port)))
}

func keepTunnelAliveHandler(w http.ResponseWriter, r *http.Request) {
	target := []string{"remote_host", "remote_port", "username"}
	request, err := extractRequestData(target, r)

	if err != nil {
		w.WriteHeader(403)
		_, _ = w.Write([]byte(err.Error()))
		return
	}

	key := request["remote_host"] + ":" + request["remote_port"] + ":" + request["username"]
	if val, ok := connector.ActiveTunnels[key]; ok {
		val.LastConnection = time.Now()
		connector.ActiveTunnels[key] = val
	}

	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(200)
	_, _ = w.Write([]byte("ok"))
}

func runOutsideCommandHandler(w http.ResponseWriter, r *http.Request) {
	target := []string{"command", "connection_type", "remote_host", "remote_port", "username", "password"}
	request, err := extractRequestData(target, r)

	if err != nil {
		w.WriteHeader(403)
		_, _ = w.Write([]byte(err.Error()))
		return
	}

	var val connector.Connection
	if val2, ok := connector.ActiveConnections[request["user_id"]+request["remote_host"]+request["username"]]; ok {
		val = val2
	} else {
		res := val.CreateShellRaw(request["connection_type"], request["username"], request["password"], request["remote_host"], request["remote_port"])
		if res == false {
			w.WriteHeader(403)
			_, _ = w.Write([]byte("Cannot create connection"))
			return
		}
	}

	val.LastConnection = time.Now()
	connector.ActiveConnections[request["user_id"]+request["remote_host"]+request["username"]] = val
	output := val.Run(request["command"])
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(200)
	_, _ = w.Write([]byte(output))
}

func verifyHandler(w http.ResponseWriter, r *http.Request) {
	target := []string{"ip_address", "username", "password", "port", "keyType"}
	request, err := extractRequestData(target, r)

	if err != nil {
		w.WriteHeader(403)
		_, _ = w.Write([]byte(err.Error()))
		return
	}

	flag := connector.VerifyAuth(request["username"], request["password"], request["ip_address"], request["port"], request["keyType"])

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
	target := []string{"local_path", "root", "parameters", "server_id"}
	request, err := extractRequestData(target, r)

	if err != nil {
		w.WriteHeader(403)
		_, _ = w.Write([]byte(err.Error()))
		return
	}
	server := postgresql.GetServer(request["server_id"])

	val, err := connector.GetConnection(request["user_id"], request["server_id"], server.IPAddress)

	if err != nil {
		w.WriteHeader(403)
		_, _ = w.Write([]byte(err.Error()))
		return
	}

	val.CreateFileConnection(request["user_id"], request["server_id"], server.IPAddress)

	var remotePath string

	if server.Os == "linux" {
		remotePath = "/tmp/" + filepath.Base(request["local_path"])
		val.Run("rm " + remotePath)
	} else {
		remotePath = val.WindowsPath + filepath.Base(request["local_path"]) + ".ps1"
	}
	flag := val.Put(request["local_path"], remotePath)

	if flag == false {
		w.WriteHeader(201)
		_, _ = w.Write([]byte("Cannot send the script!"))
		return
	}

	var output string
	if server.Os == "linux" {
		val.Run("chmod +x " + remotePath)
		if request["root"] == "yes" {
			_, password, _, keyObj := postgresql.GetServerKey(request["user_id"], request["server_id"])
			if keyObj.Type == "ssh" {
				encoded := base64.StdEncoding.EncodeToString([]byte(password))
				sudo := "echo " + encoded + " | base64 -d | sudo -S -p ' ' id 2>/dev/null 1>/dev/null; sudo "
				remotePath = sudo + remotePath
			} else if keyObj.Type == "ssh_certificate" {
				remotePath = "sudo " + remotePath
			}
		}
		output = val.Run(remotePath + " " + request["parameters"])
	} else {
		output = val.Run(val.WindowsLetter + ":\\" + remotePath + " " + request["parameters"])
	}

	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(200)
	_, _ = w.Write([]byte(output))
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  8096,
	WriteBufferSize: 8096,
}

func terminalHandler(w http.ResponseWriter, r *http.Request) {
	target := []string{"ip_address", "server_id", "token"}
	request, err := extractGetRequestData(target, r)

	if err != nil {
		w.WriteHeader(403)
		fmt.Println(err.Error())
		_, _ = w.Write([]byte(err.Error()))
		return
	}

	request["user_id"] = postgresql.GetUserIDFromToken(request["token"])

	username, password, keyPort, keyObj := postgresql.GetServerKey(request["user_id"], request["server_id"])

	data, ok := r.URL.Query()["connection_type"]
	if !ok || len(data[0]) < 1 {
		request["connection_type"] = keyObj.Type
	} else {
		request["connection_type"] = data[0]
	}

	data, ok = r.URL.Query()["connection_port"]
	if !ok || len(data[0]) < 1 {
		request["connection_port"] = keyPort
	} else {
		request["connection_port"] = data[0]
	}

	val, err := connector.GetConnectionRaw(request["user_id"], request["connection_type"], username, password, request["ip_address"], request["connection_port"])

	if err != nil {
		w.WriteHeader(403)
		fmt.Println(err.Error())
		_, _ = w.Write([]byte(err.Error()))
		return
	}

	upgrader.CheckOrigin = func(r *http.Request) bool { return true }

	ws, err := upgrader.Upgrade(w, r, nil)

	if err != nil {
		w.WriteHeader(403)
		fmt.Println(err.Error())
		_, _ = w.Write([]byte(err.Error()))
		return
	}

	err = val.Magic()
	if err != nil {
		w.WriteHeader(403)
		fmt.Println(err.Error())
		_, _ = w.Write([]byte(err.Error()))
		return
	}

	socketReader(ws, val.In, val.Out)
	return
}

func socketReader(conn *websocket.Conn, in io.WriteCloser, out io.Reader) {
	for {
		go socketReadLoop(conn, out)
		_, p, err := conn.ReadMessage()
		_, err = in.Write(p)
		if err != nil {
			log.Println(err)
			return
		}
	}
}

func socketReadLoop(conn *websocket.Conn, out io.Reader) {
	for {
		data := make([]byte, 1024)

		receivedSize, _ := out.Read(data)
		if receivedSize != 1 {
			fmt.Println("BEKLE")
			time.Sleep(time.Second * 2)
		} else {
			fmt.Println("DEVAM")
		}
		fmt.Println(strconv.Itoa(receivedSize))

		if err := conn.WriteMessage(2, data); err != nil {
			log.Println(err)
			return
		}
	}
}
