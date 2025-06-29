package monitor

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

func (d *Daemon) PublishCPUUsage() {

	readCPUStats := func() (map[string][]uint64, error) {
		stat, err := os.ReadFile("/proc/stat")
		if err != nil {
			return nil, fmt.Errorf("error reading /proc/stat: %v", err)
		}

		cpuStats := make(map[string][]uint64)
		for _, line := range strings.Split(string(stat), "\n") {
			fields := strings.Fields(line)
			if len(fields) < 2 || !strings.HasPrefix(fields[0], "cpu") {
				continue
			}

			core := fields[0]
			values := make([]uint64, len(fields)-1)
			for i, v := range fields[1:] {
				values[i], _ = strconv.ParseUint(v, 10, 64)
			}
			cpuStats[core] = values
		}
		return cpuStats, nil
	}

	stats1, err := readCPUStats()
	if err != nil {
		fmt.Printf("Error reading CPU stats: %v\n", err)
		return
	}

	// Sleep for a second to get delta
	time.Sleep(time.Second)

	stats2, err := readCPUStats()
	if err != nil {
		fmt.Printf("Error reading CPU stats: %v\n", err)
		return
	}

	cpuUsage := make(map[string]float64)
	for core, s2 := range stats2 {
		s1, ok := stats1[core]
		if !ok {
			continue
		}

		var total2, total1, idle2, idle1 uint64
		for i, v := range s2 {
			total2 += v
			if i == 3 || i == 4 { // idle and iowait
				idle2 += v
			}
		}
		for i, v := range s1 {
			total1 += v
			if i == 3 || i == 4 { // idle and iowait
				idle1 += v
			}
		}

		totalDelta := total2 - total1
		idleDelta := idle2 - idle1

		usage := 100.0 * (1.0 - float64(idleDelta)/float64(totalDelta))
		if core == "cpu" {
			cpuUsage["Total"] = usage
		} else {
			cpuUsage[strings.ToUpper(core)] = usage
		}
	}

	jsonData, err := json.Marshal(cpuUsage)
	if err != nil {
		fmt.Printf("Error marshaling JSON: %v\n", err)
		return
	}
	d.MQTTPublish(d.CreateTopic("cpu"), jsonData)
}
