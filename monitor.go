package monitor

import (
	"flag"
	"fmt"
	"os"
	"time"
)

type Daemon struct {
	MQTTBroker      string
	MQTTPort        int
	MQTTUsername    string
	MQTTPassword    string
	MQTTTopicPrefix string
}

func NewDaemon() *Daemon {
	flagsBrokerAddress := flag.String("b", "", "Broker address")
	flagsBrokerPort := flag.Int("P", 1883, "Broker port")
	flagsBrokerUsername := flag.String("u", "", "Broker username")
	flagsBrokerPassword := flag.String("p", "", "Broker password")
	flagsTopicPrefix := flag.String("t", "", "Topic prefix")
	flag.Parse()

	hostname, err := os.Hostname()
	if err != nil {
		fmt.Printf("Error getting hostname: %v\n", err)
		os.Exit(1)
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
		d.PublishMemoryUsage()

		time.Sleep(time.Second * 5)
	}
}
