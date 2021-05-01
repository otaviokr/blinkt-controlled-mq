package mqtt

import (
	"fmt"
	"time"

	mq "github.com/eclipse/paho.mqtt.golang"
)

type MqttClient struct {
	clientId string
	Client mq.Client
}

func NewClient(broker string, port int) (*MqttClient, error) {
	clientId := fmt.Sprintf("okr_p_%d", time.Now().Unix())

	// Connecting to MQTT.
	opts := mq.NewClientOptions()
	opts.AddBroker(fmt.Sprintf("tcp://%s:%d", broker, port))
	opts.SetClientID(clientId) // ClientID can be anything, just make sure it is unique
	// opts.SetUsername("username") // Not needed because we are using the public HiveMQ
	// opts.SetPassword("password") // Not needed because we are using the public HiveMQ
	// opts.SetDefaultPublishHandler(getMessageReceivedHandler(device))
	// opts.OnConnect = getOnConnectHandler()
	// opts.OnConnectionLost = getConnectionLostHandler()

	client := mq.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		return &MqttClient{}, token.Error()
	}

	return &MqttClient {
		clientId: clientId,
		Client:   client,
	}, nil
}
