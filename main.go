package main

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/otaviokr/blinkt-controlled-mq/blinkt"
	log "github.com/sirupsen/logrus"
)

var (
	blinktDevice *blinkt.Dev

	messagePubHandler mqtt.MessageHandler = func(client mqtt.Client, msg mqtt.Message) {
		log.Debugf("Received message: %s from topic: %s\n", msg.Payload(), msg.Topic())
		UpdateLed(string(msg.Payload()))
	}

	connectHandler mqtt.OnConnectHandler = func(client mqtt.Client) {
		log.Info("Connected to MQTT")
	}

	connectLostHandler mqtt.ConnectionLostHandler = func(client mqtt.Client, err error) {
		log.Errorf("Connect lost: %v\n", err)
	}
)

func main() {
	// MQTT connection details.
	// broker := "broker.emqx.io"
	// port := 1883
	// username := "emqx"
	// password := "public"
	broker := "broker.hivemq.com"
	port := 1883
	clientID := "blinkt_rpi0_01"
	topic := "test_okr"

	// The mandatory init call for Blinkt.
	err := blinkt.Init()
	if err != nil {
		panic(err)
	}

	// Instantiating a new device for Blinkt.
	blinktDevice, err = blinkt.NewDev()
	if err != nil {
		panic(err)
	}
	blinktDevice.SetClearOnExit(true)
	log.Info("Connected to Blinkt")

	opts := mqtt.NewClientOptions()
	opts.AddBroker(fmt.Sprintf("tcp://%s:%d", broker, port))
	opts.SetClientID(clientID)
	// opts.SetUsername(username)
	// opts.SetPassword(password)
	opts.SetDefaultPublishHandler(messagePubHandler)
	opts.OnConnect = connectHandler
	opts.OnConnectionLost = connectLostHandler
	client := mqtt.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}
	log.Info("Connected to MQTT")

	sub(client, topic)

	log.Info("Processing incoming messages")
	for {
		time.Sleep(3 * time.Second)
	}
}

// ParseColor will convert the hexadecimal in the string into a valid int.
func ParseColor(c string) int {
	r, err := strconv.ParseInt(c, 16, 32)
	if err != nil {
		log.WithFields(
			log.Fields{
				"ColorComponent": c,
			},
		).Error("Could not parse the color component")
		return 0
	}
	return int(r)
}

// SetPixel will parse the incomming details for a LED (a "pixel") to configure it correctly.
func SetPixel(led, rgb string) (int, int, int, int) {
	index, err := strconv.Atoi(led[len(led)-1:])
	if err != nil {
		log.WithFields(
			log.Fields{
				"led": led,
				"rgb": rgb,
			},
		).Error("Could not convert LED index")
		return -1, 0, 0, 0
	}

	re := regexp.MustCompile("#([0-9a-fA-F]{2})([0-9a-fA-F]{2})([0-9a-fA-F]{2})")
	match := re.FindStringSubmatch(rgb)

	red := 0
	green := 0
	blue := 0
	if len(match) == 4 {
		red = ParseColor(match[1])
		green = ParseColor(match[2])
		blue = ParseColor(match[3])
	} else {
		log.WithFields(
			log.Fields{
				"led": led,
				"rgb": rgb,
			},
		).Error("Could not parse RGB")
	}

	return index - 1, red, green, blue
}

// UpdateLed will collect the LED values and brightness and send them to Blinkt.
func UpdateLed(m string) {
	leds := strings.Split(m, " ")

	for i, l := range leds {
		c := strings.Split(l, "#")

		red, err := strconv.Atoi(c[1])
		if err != nil {
			log.WithFields(
				log.Fields{
					"value": c,
					"parsed": c[1],
					"component": "red",
					"index": i,
					"error": err.Error(),
				}).Error("Failed to parse component of LED")
				red = 0
		}

		green, err := strconv.Atoi(c[2])
		if err != nil {
			log.WithFields(
				log.Fields{
					"value": c,
					"parsed": c[2],
					"component": "green",
					"index": i,
					"error": err.Error(),
				}).Error("Failed to parse component of LED")
				green = 0
		}

		blue, err := strconv.Atoi(c[3])
		if err != nil {
			log.WithFields(
				log.Fields{
					"value": c,
					"parsed": c[3],
					"component": "blue",
					"index": i,
					"error": err.Error(),
				}).Error("Failed to parse component of LED")
				blue = 0
		}

		brightness, err := strconv.ParseFloat(c[4], 32)
		if err != nil {
			log.WithFields(
				log.Fields{
					"value": c,
					"parsed": c[3],
					"component": "brightness",
					"index": i,
					"error": err.Error(),
				}).Error("Failed to parse component of LED")
				brightness = 0.5
		}

		blinktDevice.SetPixelWithBright(i, red, green, blue, brightness)
	}
	blinktDevice.Show()
}

func sub(client mqtt.Client, topic string) {
	token := client.Subscribe(topic, 1, nil)
	token.Wait()
	log.WithFields(
		log.Fields{
			"topic": topic,
		}).Info("Subscribed to topic")
}