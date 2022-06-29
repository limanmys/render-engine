package recycle

import (
	"net"
	"time"

	"github.com/limanmys/go/connector"
)

//Start Start
func Start() {
	now := time.Now()
	for key, data := range connector.ActiveConnections {
		if now.Sub(data.LastConnection).Seconds() > 266 {
			connector.CloseAllConnections(connector.ActiveConnections[key])
			delete(connector.ActiveConnections, key)
		}

		if data.SSH != nil {
			addr := net.JoinHostPort(data.IpAddr, data.Port)
			_, err := net.DialTimeout("tcp", addr, 10*time.Second)
			if err != nil {
				connector.CloseAllConnections(connector.ActiveConnections[key])
				delete(connector.ActiveConnections, key)
			}
		}
	}
	for key, data := range connector.ActiveTunnels {
		if now.Sub(data.LastConnection).Seconds() > 266 {
			connector.ActiveTunnels[key].Tunnel.Stop()
			delete(connector.ActiveTunnels, key)
		}
	}
	time.Sleep(time.Second)
	go Start()
}
