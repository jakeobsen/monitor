package monitor

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"syscall"
)

type DiskStats struct {
	Mountpoint    string  `json:"mountpoint"`
	PartitionName string  `json:"partition_name"`
	TotalSpace    uint64  `json:"total_space"`
	UsedSpace     uint64  `json:"used_space"`
	FreeSpace     uint64  `json:"free_space"`
	UsagePercent  float64 `json:"usage_percent"`
	TotalInodes   uint64  `json:"total_inodes"`
	UsedInodes    uint64  `json:"used_inodes"`
	FreeInodes    uint64  `json:"free_inodes"`
	InodesPercent float64 `json:"inodes_percent"`
}

func isPhysicalDisk(fstype, mountpoint string) bool {
	// List of filesystem types to include
	validFS := map[string]bool{
		"ext2":  true,
		"ext3":  true,
		"ext4":  true,
		"xfs":   true,
		"btrfs": true,
		"ntfs":  true,
		"vfat":  true,
	}

	// Skip special filesystems and mount points
	if strings.HasPrefix(mountpoint, "/dev") ||
		strings.HasPrefix(mountpoint, "/sys") ||
		strings.HasPrefix(mountpoint, "/proc") ||
		strings.HasPrefix(mountpoint, "/run") {
		return false
	}

	return validFS[fstype]
}

func (d *Daemon) PublishDiskUsage() {

	content, err := os.ReadFile("/proc/mounts")
	if err != nil {
		fmt.Printf("Error reading /proc/mounts: %v\n", err)
		return
	}

	var stats []DiskStats
	for _, line := range strings.Split(string(content), "\n") {
		fields := strings.Fields(line)
		if len(fields) < 3 {
			continue
		}

		fstype := fields[2]
		mountpoint := fields[1]

		// Skip non-physical and special filesystems
		if !isPhysicalDisk(fstype, mountpoint) {
			continue
		}
		var statfs syscall.Statfs_t
		if err := syscall.Statfs(mountpoint, &statfs); err != nil {
			continue
		}

		totalSpace := statfs.Blocks * uint64(statfs.Bsize)
		freeSpace := statfs.Bfree * uint64(statfs.Bsize)
		usedSpace := totalSpace - freeSpace
		usagePercent := 0.0
		if totalSpace > 0 {
			usagePercent = float64(usedSpace) / float64(totalSpace) * 100
		}

		totalInodes := statfs.Files
		freeInodes := statfs.Ffree
		usedInodes := totalInodes - freeInodes
		inodesPercent := 0.0
		if totalInodes > 0 {
			inodesPercent = float64(usedInodes) / float64(totalInodes) * 100
		}

		deviceName := fields[0]
		label, err := os.ReadFile(fmt.Sprintf("/dev/disk/by-label/%s", strings.TrimPrefix(deviceName, "/dev/")))
		partitionName := mountpoint
		if strings.Contains(deviceName, "/dev/mapper/") {
			parts := strings.Split(deviceName, "/")
			lvName := parts[len(parts)-1]
			if strings.Contains(lvName, "-") {
				partitionName = lvName
			}
		} else if err == nil {
			partitionName = string(label)
		} else if mountpoint == "/" {
			partitionName = "root"
		}
		stats = append(stats, DiskStats{
			Mountpoint:    mountpoint,
			PartitionName: partitionName,
			TotalSpace:    totalSpace,
			UsedSpace:     usedSpace,
			FreeSpace:     freeSpace,
			UsagePercent:  usagePercent,
			TotalInodes:   totalInodes,
			UsedInodes:    usedInodes,
			FreeInodes:    freeInodes,
			InodesPercent: inodesPercent,
		})
	}

	// Publish individual disk stats
	for _, stat := range stats {
		diskData, err := json.Marshal(stat)
		if err != nil {
			fmt.Printf("Error marshaling disk JSON: %v\n", err)
			continue
		}
		sanitizedName := strings.ReplaceAll(stat.PartitionName, "/", "")
		d.MQTTPublish(d.CreateTopic("disk/"+sanitizedName), diskData)
	}
}
