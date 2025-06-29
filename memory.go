package monitor

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func (d *Daemon) PublishMemoryUsage() {
	meminfo, err := os.ReadFile("/proc/meminfo")
	if err != nil {
		fmt.Printf("Error reading /proc/meminfo: %v\n", err)
		return
	}

	var total, free, available, cached, shared uint64
	for _, line := range strings.Split(string(meminfo), "\n") {
		fields := strings.Fields(line)
		if len(fields) < 2 {
			continue
		}

		value, err := strconv.ParseUint(fields[1], 10, 64)
		if err != nil {
			continue
		}

		switch fields[0] {
		case "MemTotal:":
			total = value
		case "MemFree:":
			free = value
		case "MemAvailable:":
			available = value
		case "Cached:":
			cached = value
		case "Shmem:":
			shared = value
		}
	}

	stats := struct {
		Total     uint64 `json:"total"`
		Used      uint64 `json:"used"`
		Free      uint64 `json:"free"`
		Cache     uint64 `json:"cache"`
		Available uint64 `json:"available"`
		Shared    uint64 `json:"shared"`
	}{
		Total:     total,
		Used:      total - free - cached,
		Free:      free,
		Cache:     cached,
		Available: available,
		Shared:    shared,
	}
	jsonData, err := json.Marshal(stats)
	if err != nil {
		fmt.Printf("Error marshaling JSON: %v\n", err)
		return
	}

	d.MQTTPublish(d.CreateTopic("memory"), jsonData)
}
