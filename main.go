package main

import (
	"fmt"
	"sync"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

// Creating Options for mqtt client
func createClientOptions(BrokerURL string) *mqtt.ClientOptions {

	opts := mqtt.NewClientOptions().AddBroker(BrokerURL)
	opts.SetKeepAlive(time.Second * 5)
	opts.SetCleanSession(true)
	opts.SetAutoReconnect(true)

	return opts
}

// Connecting to mqtt client
func connect(opts *mqtt.ClientOptions, BrokerURL string) mqtt.Client {

	client := mqtt.NewClient(opts)

	token := client.Connect()
	if token.Wait() && token.Error() != nil {
		panic(token.Error())
	} else {
		fmt.Println("Mqtt server is ready: ", BrokerURL)
	}

	return client
}

func main() {

	const TOPIC = "myTopic"
	const BrokerURL = "tcp://localhost:1883"

	// Creating mqtt client options
	opts := createClientOptions(BrokerURL)

	// Connect to mqtt client
	client := connect(opts, BrokerURL)

	var wg sync.WaitGroup
	wg.Add(1)

	/* ===== Subscribe ===== */
	if token := client.Subscribe(TOPIC, 0, func(client mqtt.Client, msg mqtt.Message) {

		fmt.Printf("Subscriber Message recived:  %s\n", msg.Payload())
		// wg.Done()

		// Publish
		if token := client.Publish("received", 0, false, "repeat publish"); token.Wait() && token.Error() != nil {
			panic(token.Error())
		}

	}); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}

	/* ======= Publish ======= */
	// if token := client.Publish("received", 0, false, "sending again"); token.Wait() && token.Error() != nil {
	// 	panic(token.Error())
	// }

	wg.Wait()
}
