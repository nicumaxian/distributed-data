package mediator

import (
	"github.com/maxiannicu/distributed-data/network_dto"
	"time"
	"github.com/maxiannicu/distributed-data/utils"
	"net"
	"github.com/maxiannicu/distributed-data/network"
)

func (application *Application) findMasterNode() (*network.EndPoint, bool) {
	bytes, err := network_dto.CreateRequest(network_dto.DiscoveryRequestType, network_dto.DiscoveryRequest{
		ResponseEndPoint: application.discoveryUdpListener.LocalEndPoint(),
	})

	if err != nil {
		application.logger.Fatal(err)
	}

	application.logger.Println("Sending discovery command")
	application.discoveryUdpSender.Write(bytes)
	time.Sleep(application.discoveryDuration)

	responses := application.getDiscoveredNodes()

	application.logger.Println("Discovery done. Found", len(responses), "nodes")

	if len(responses) == 0 {
		return nil, false
	} else {
		master := responses[0]
		for _, node := range responses {
			if node.DataSize > master.DataSize {
				master = node
			}
		}

		return &master.ConnectionEndPoint, true
	}
}

func (application *Application) getDiscoveredNodes() []network_dto.DiscoveryResponse {
	responses := make([]network_dto.DiscoveryResponse, 0)
	application.discoveryUdpListener.SetReadTimeOut(time.Now().Add(application.discoveryDuration))
	for {
		bytes, err := application.discoveryUdpListener.Read()
		if err != nil {
			if err, ok := err.(net.Error); ok && err.Timeout() {
				break
			}
			application.logger.Fatal(err)
		} else {
			response := network_dto.DiscoveryResponse{}
			utils.Deserealize(utils.JsonFormat, bytes, &response)
			application.logger.Println("Discovered", response.ConnectionEndPoint)
			responses = append(responses, response)
		}
	}
	return responses
}
