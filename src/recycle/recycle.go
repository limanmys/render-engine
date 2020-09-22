package recycle

import (
	"renderer/src/connector"
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
	time.Sleep(time.Second)
	go Start()
}
