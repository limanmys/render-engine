package ssh

import (
	"log"
	"strconv"
	"time"

	"github.com/phayes/freeport"
	"github.com/rgzr/sshtun"
)

// //Endpoint Endpoint
// type Endpoint struct {
// 	Host string
// 	Port int
// 	User string
// }

// //Tunnel Tunnel
// type Tunnel struct {
// 	LocalPort  string
// 	RemotePort string
// 	Connection *ssh.Client
// 	Server     sqlite.ServerModel
// }

// //Start Start
// func (tunnel *Tunnel) Start() error {
// 	listener, err := net.Listen("tcp", "127.0.0.1:5555")
// 	if err != nil {
// 		return err
// 	}
// 	defer listener.Close()
// 	tunnel.LocalPort = string(listener.Addr().(*net.TCPAddr).Port)
// 	for {
// 		conn, err := listener.Accept()
// 		if err != nil {
// 			return err
// 		}
// 		go tunnel.forward(conn)
// 	}
// }
// func (tunnel *Tunnel) forward(localConn net.Conn) {
// 	remoteConn, err := tunnel.Connection.Dial("tcp", tunnel.Server.IPAddress+":"+tunnel.RemotePort)
// 	if err != nil {
// 		return
// 	}
// 	copyConn := func(writer, reader net.Conn) {
// 		_, err := io.Copy(writer, reader)
// 		if err != nil {
// 		}
// 	}
// 	go copyConn(localConn, remoteConn)
// 	go copyConn(remoteConn, localConn)
// }

// //NewTunnel NewTunnel
// func NewTunnel(conn *ssh.Client, server sqlite.ServerModel, destinationPort string) *Tunnel {
// 	Tunnel := &Tunnel{
// 		Connection: conn,
// 		RemotePort: destinationPort,
// 		Server:     server,
// 	}
// 	return Tunnel
// }

//CreateTunnel CreateTunnel
func CreateTunnel(remoteHost string, remotePort string, username string, password string) int {
	port, err := freeport.GetFreePort()
	if err != nil {
		log.Fatal(err)
	}
	remotePortInt, _ := strconv.Atoi(remotePort)
	sshTun := sshtun.New(port, remoteHost, remotePortInt)
	sshTun.SetDebug(true)
	sshTun.SetLocalHost("127.0.0.1")
	sshTun.SetPassword(password)
	sshTun.SetUser(username)
	go func() {
		for {
			if err := sshTun.Start(); err != nil {
				log.Printf("SSH tunnel stopped: %s", err.Error())
				time.Sleep(time.Second)
			}
		}
	}()
	return port
}
