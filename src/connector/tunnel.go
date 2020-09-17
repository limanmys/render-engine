package connector

import (
	"log"
	"strconv"
	"time"

	"github.com/phayes/freeport"
	"github.com/rgzr/sshtun"
)

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
