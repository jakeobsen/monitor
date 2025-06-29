package monitor

import (
	"flag"
	"fmt"
	"os"
	"strconv"
	"time"
)

type Daemon struct {
	MQTTBroker      string
	MQTTPort        int
	MQTTUsername    string
	MQTTPassword    string
	MQTTTopicPrefix string
	EnvVarName      string
}

func NewDaemon() *Daemon {
	envBrokerAddress := os.Getenv("MQTT_SERVER")
	envBrokerPort := 1883
	if port := os.Getenv("MQTT_PORT"); port != "" {
		if p, err := strconv.Atoi(port); err == nil {
			envBrokerPort = p
		}
	}
	envBrokerUsername := os.Getenv("MQTT_USERNAME")
	envBrokerPassword := os.Getenv("MQTT_PASSWORD")
	envTopicPrefix := os.Getenv("TOPIC_PREFIX")
	envOverrideHostname := os.Getenv("HOSTNAME")

	flagsBrokerAddress := flag.String("b", envBrokerAddress, "Broker address")
	flagsBrokerPort := flag.Int("P", envBrokerPort, "Broker port")
	flagsBrokerUsername := flag.String("u", envBrokerUsername, "Broker username")
	flagsBrokerPassword := flag.String("p", envBrokerPassword, "Broker password")
	flagsTopicPrefix := flag.String("t", envTopicPrefix, "Topic prefix")
	flagsOverrideHostname := flag.String("h", envOverrideHostname, "Hostname override")
	flag.Parse()

	if *flagsBrokerUsername == "" {
		fmt.Println("Error: MQTT username is required")
		os.Exit(1)
	}
	if *flagsBrokerPassword == "" {
		fmt.Println("Error: MQTT password is required")
		os.Exit(1)
	}
	if *flagsTopicPrefix == "" {
		fmt.Println("Error: Topic prefix is required")
		os.Exit(1)
	}

	hostname, err := os.Hostname()
	if err != nil {
		fmt.Printf("Error getting hostname: %v\n", err)
		os.Exit(1)
	}

	if *flagsOverrideHostname != "" {
		hostname = *flagsOverrideHostname
	}

	topicPrefix := fmt.Sprintf("%s/%s", *flagsTopicPrefix, hostname)

	d := Daemon{
		MQTTBroker:      *flagsBrokerAddress,
		MQTTPort:        *flagsBrokerPort,
		MQTTUsername:    *flagsBrokerUsername,
		MQTTPassword:    *flagsBrokerPassword,
		MQTTTopicPrefix: topicPrefix,
	}

	return &d
}

func (d *Daemon) Run() {
	for {
		go d.PublishMemoryUsage()
		go d.PublishCPUUsage()
		time.Sleep(time.Second * 5)
	}
}
