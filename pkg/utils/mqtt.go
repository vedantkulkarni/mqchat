package utils

import (
	"fmt"

	mqtt "github.com/mochi-mqtt/server/v2"
	"github.com/mochi-mqtt/server/v2/packets"
)


func PublishMessage(client *mqtt.Client, payload interface{}) error {
		
	// Publish a message to the client
	err := client.WritePacket(packets.Packet{
		FixedHeader: packets.FixedHeader{
			Type: packets.Publish,
			Qos:  0,
		},
		Payload: []byte(fmt.Sprintf("%v", payload)),
	})

	if err != nil {
		return err
	}


	return nil
}

func PublishError(client *mqtt.Client, err error) error {
	errPacket := packets.Packet{
		FixedHeader: packets.FixedHeader{
			Type: packets.Publish,
			Qos:  0,
		},
		Payload: []byte(fmt.Sprintf("%v", err)),
	}

	err = client.WritePacket(errPacket)

	if err != nil {
		return err
	}

	return nil
}