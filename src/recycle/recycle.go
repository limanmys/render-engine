package recycle

import (
	"fmt"
	"renderer/src/constants"
	"time"
)

//Start Start
func Start() {
	now := time.Now()
	for key, data := range constants.ActiveConnections {
		if now.Sub(data.LastConnection).Seconds() > 30 {
			constants.CloseAllConnections(constants.ActiveConnections[key])
			delete(constants.ActiveConnections, key)
			fmt.Printf("Recycling !%s\n", key)
		}
	}
	time.Sleep(time.Second)
	go Start()
}
