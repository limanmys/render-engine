package recycle

import (
	"net"
	"sync"
	"time"

	"github.com/limanmys/render-engine/connector"
)

var mut sync.Mutex

//Start Start
func Start() {
	now := time.Now()
	for key, data := range connector.ActiveConnections {
		if now.Sub(data.LastConnection).Seconds() > 266 {
			mut.Lock()
			connector.CloseAllConnections(connector.ActiveConnections[key])
			delete(connector.ActiveConnections, key)
			mut.Unlock()
			continue
		}

		if data.SSH != nil {
			if data.IpAddr != "" && data.Port != "" {
				ipAddress := data.IpAddr
				for i := 0; i < 10; i++ {
					addr, err := net.LookupIP(ipAddress)
					if err == nil {
						ipAddress = addr[0].String()
						break
					}
				}

				addr := net.JoinHostPort(ipAddress, data.Port)
				_, err := net.DialTimeout("tcp", addr, 10*time.Second)
				if err != nil {
					mut.Lock()
					connector.CloseAllConnections(connector.ActiveConnections[key])
					delete(connector.ActiveConnections, key)
					mut.Unlock()
					continue
				}
			}

			ch := make(chan int, 1)
			go func() {
				select {
				case <-time.After(10 * time.Second):
				case <-ch:
					return
				default:
					data.SSH.SendRequest("keepalive@liman.dev", true, nil)
					ch <- 1
				}
			}()

			select {
			case <-ch:
				continue
			case <-time.After(10 * time.Second):
				mut.Lock()
				connector.CloseAllConnections(connector.ActiveConnections[key])
				delete(connector.ActiveConnections, key)
				mut.Unlock()
				continue
			}
		}
	}
	for key, data := range connector.ActiveTunnels {
		if now.Sub(data.LastConnection).Seconds() > 266 {
			mut.Lock()
			connector.ActiveTunnels[key].Tunnel.Stop()
			delete(connector.ActiveTunnels, key)
			mut.Unlock()
		}
	}
	time.Sleep(time.Second)
	go Start()
}
