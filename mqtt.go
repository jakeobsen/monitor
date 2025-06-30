package monitor

import (
	"fmt"
	MQTT "github.com/eclipse/paho.mqtt.golang"
)

func (d *Daemon) CreateTopic(topic string) string {
	return fmt.Sprintf("%s/%s/%s", d.MQTT.TopicPrefix, d.System.Hostname, topic)
}

func (d *Daemon) MQTTPublish(topic string, jsonData []byte) {
	opts := MQTT.NewClientOptions()
	opts.AddBroker(fmt.Sprintf("tcp://%s:%d", d.MQTT.Server, d.MQTT.Port))
	opts.SetUsername(d.MQTT.Username)
	opts.SetPassword(d.MQTT.Password)

	//fmt.Printf("Publishing to %s: %s\n", topic, jsonData)

	client := MQTT.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		//fmt.Printf("MQTT connection error: %v\n", token.Error())
		return
	}
	defer client.Disconnect(250)

	if token := client.Publish(topic, 0, false, jsonData); token.Wait() && token.Error() != nil {
		//fmt.Printf("MQTT publish error: %v\n", token.Error())
		return
	}
}
