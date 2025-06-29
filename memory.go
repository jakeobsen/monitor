package monitor

import (
	"encoding/json"
	"fmt"
	"runtime"
)

func (d *Daemon) PublishMemoryUsage() {
	var memStats runtime.MemStats
	runtime.ReadMemStats(&memStats)

	stats := struct {
		Total     uint64 `json:"total"`
		Used      uint64 `json:"used"`
		Free      uint64 `json:"free"`
		Cache     uint64 `json:"cache"`
		Available uint64 `json:"available"`
	}{
		Total:     memStats.Sys,
		Used:      memStats.Alloc,
		Free:      memStats.Sys - memStats.Alloc,
		Cache:     memStats.HeapIdle,
		Available: memStats.Sys - memStats.Alloc,
	}

	jsonData, err := json.Marshal(stats)
	if err != nil {
		fmt.Printf("Error marshaling JSON: %v\n", err)
		return
	}

	d.MQTTPublish(d.CreateTopic("memory"), jsonData)
}
