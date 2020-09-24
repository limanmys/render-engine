package recycle

import (
	"github.com/limanmys/go/connector"
	"time"
)

//Start Start
func Start() {
	now := time.Now()
	for key, data := range connector.ActiveConnections {
		if now.Sub(data.LastConnection).Seconds() > 30 {
			connector.CloseAllConnections(connector.ActiveConnections[key])
			delete(connector.ActiveConnections, key)
		}
	}
	for key, data := range connector.ActiveTunnels {
		if now.Sub(data.LastConnection).Seconds() > 30 {
			connector.ActiveTunnels[key].Tunnel.Stop()
			delete(connector.ActiveTunnels, key)
		}
	}
	time.Sleep(time.Second)
	go Start()
}
