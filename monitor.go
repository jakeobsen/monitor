package monitor

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"time"
)

type Daemon struct {
	System struct {
		Hostname string `json:"hostname"`
	}
	MQTT struct {
		Server      string `json:"server"`
		Port        int    `json:"port"`
		Username    string `json:"username"`
		Password    string `json:"password"`
		TopicPrefix string `json:"topicPrefix"`
	} `json:"mqtt"`
	Modules struct {
		Disk struct {
			Enabled bool `json:"enabled"`
		} `json:"disk"`
		Memory struct {
			Enabled bool `json:"enabled"`
		} `json:"memory"`
		CPU struct {
			Enabled bool `json:"enabled"`
		} `json:"cpu"`
	} `json:"modules"`
}

func NewDaemon() *Daemon {
	configFilename := flag.String("c", "monitor.json", "Path to config file")
	flag.Parse()

	hostname, err := os.Hostname()
	if err != nil {
		fmt.Printf("Error getting hostname: %v\n", err)
		os.Exit(1)
	}

	file, err := os.ReadFile(*configFilename)
	if err != nil {
	}

	var d Daemon
	if err := json.Unmarshal(file, &d); err != nil {
	}

	d.System.Hostname = hostname

	return &d
}

func (d *Daemon) Run() {
	for {
		if d.Modules.Memory.Enabled == true {
			go d.PublishMemoryUsage()
		}
		if d.Modules.CPU.Enabled == true {
			go d.PublishCPUUsage()
		}
		if d.Modules.Disk.Enabled == true {
			go d.PublishDiskUsage()
		}
		time.Sleep(time.Second * 5)
	}
}
