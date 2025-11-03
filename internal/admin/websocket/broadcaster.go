package websocket

import (
	"time"
)

type StatusGetter interface {
	GetStatus() interface{}
}

func StartStatusBroadcaster(hub *Hub, getter interface{}, interval time.Duration) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for range ticker.C {
		if sg, ok := getter.(interface{ GetStatus() interface{} }); ok && sg != nil {
			status := sg.GetStatus()
			hub.BroadcastStatus(status)
		}
	}
}
