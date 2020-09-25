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
	if val, ok := ActiveTunnels[remoteHost+":"+remotePort+":"+username]; ok {
		val.LastConnection = time.Now()
		return val.Port
	}

	port, err := freeport.GetFreePort()
	if err != nil {
		log.Fatal(err)
		return 0
	}
	remotePortInt, _ := strconv.Atoi(remotePort)
	sshTun := sshtun.New(port, remoteHost, remotePortInt)
	sshTun.SetLocalHost("127.0.0.1")
	sshTun.SetPassword(password)
	sshTun.SetUser(username)
	tunnelState := sshtun.StateStarting
	sshTun.SetConnState(func(tun *sshtun.SSHTun, state sshtun.ConnState) {
		tunnelState = state
	})

	go func() {
		for {
			if err := sshTun.Start(); err != nil {
				log.Printf("SSH tunnel stopped: %s", err.Error())
				delete(ActiveTunnels, remoteHost+":"+remotePort+":"+username)
				break
			}
		}
	}()

	for {
		if tunnelState == sshtun.StateStarted {
			break
		}
	}

	tunnel := Tunnel{Tunnel: sshTun, Port: port, LastConnection: time.Now()}
	ActiveTunnels[remoteHost+":"+remotePort+":"+username] = tunnel
	return port
}
