package mediator

import "github.com/maxiannicu/distributed-data/network"

type ApplicationConfig struct {
	ListenEndPoint             network.EndPoint
	MulticastEndPoint          network.EndPoint
	DiscoveryDuration          int
}
