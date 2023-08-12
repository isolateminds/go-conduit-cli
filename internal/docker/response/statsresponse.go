package response

import (
	"encoding/json"
	"fmt"
	"math"
)

type FormatedContainerStats struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	CpuUsage    string `json:"cpuUsage"`
	MemoryUsage string `json:"memoryUsage"`
	NetworkIO   string `json:"networkIO"`
	DiskIO      string `json:"diskIO"`
}
type ContainerStats struct {
	Name      string `json:"name"`
	ID        string `json:"id"`
	Read      string `json:"read"`
	Preread   string `json:"preread"`
	PidsStats struct {
		Current int64       `json:"current"`
		Limit   json.Number `json:"limit"`
	} `json:"pidsStats"`
	BlkioStats struct {
		IoServiceBytesRecursive []interface{} `json:"ioServiceBytesRecursive"`
		IoServicedRecursive     []interface{} `json:"ioServicedRecursive"`
		IoQueueRecursive        []interface{} `json:"ioQueueRecursive"`
		IoServiceTimeRecursive  []interface{} `json:"ioServiceTimeRecursive"`
		IoWaitTimeRecursive     []interface{} `json:"ioWaitTimeRecursive"`
		IoMergedRecursive       []interface{} `json:"ioMergedRecursive"`
		IoTimeRecursive         []interface{} `json:"ioTimeRecursive"`
		SectorsRecursive        []interface{} `json:"sectorsRecursive"`
	} `json:"blkioStats"`
	NumProcs     int64    `json:"numProcs"`
	StorageStats struct{} `json:"storageStats"`
	CpuStats     struct {
		CpuUsage struct {
			TotalUsage        int64 `json:"totalUsage"`
			UsageInKernelMode int64 `json:"usageInKernelmode"`
			UsageInUserMode   int64 `json:"usageInUsermode"`
		} `json:"cpuUsage"`
		SystemCPUUsage int64 `json:"systemCpuUsage"`
		OnlineCPUs     int64 `json:"onlineCpus"`
		ThrottlingData struct {
			Periods          int64 `json:"periods"`
			ThrottledPeriods int64 `json:"throttledPeriods"`
			ThrottledTime    int64 `json:"throttledTime"`
		} `json:"throttlingData"`
	} `json:"cpuStats"`
	PreCPUStats struct {
		CpuUsage struct {
			TotalUsage        int64 `json:"totalUsage"`
			UsageInKernelMode int64 `json:"usageInKernelmode"`
			UsageInUserMode   int64 `json:"usageInUsermode"`
		} `json:"cpuUsage"`
		SystemCPUUsage int64 `json:"systemCpuUsage"`
		OnlineCPUs     int64 `json:"onlineCpus"`
		ThrottlingData struct {
			Periods          int64 `json:"periods"`
			ThrottledPeriods int64 `json:"throttledPeriods"`
			ThrottledTime    int64 `json:"throttledTime"`
		} `json:"throttlingData"`
	} `json:"precpuStats"`
	MemoryStats struct {
		Usage int64 `json:"usage"`
		Stats struct {
			ActiveAnon          int64 `json:"activeAnon"`
			ActiveFile          int64 `json:"activeFile"`
			Anon                int64 `json:"anon"`
			AnonTHP             int64 `json:"anonTHP"`
			File                int64 `json:"file"`
			FileDirty           int64 `json:"fileDirty"`
			FileMapped          int64 `json:"fileMapped"`
			FileWriteBack       int64 `json:"fileWriteback"`
			InactiveAnon        int64 `json:"inactiveAnon"`
			InactiveFile        int64 `json:"inactiveFile"`
			KernelStack         int64 `json:"kernelStack"`
			PgActivate          int64 `json:"pgActivate"`
			PgDeactivate        int64 `json:"pgDeactivate"`
			PgFault             int64 `json:"pgFault"`
			PgLazyFree          int64 `json:"pgLazyFree"`
			PgLazyFreed         int64 `json:"pgLazyFreed"`
			PgMajFault          int64 `json:"pgMajFault"`
			PgRefill            int64 `json:"pgRefill"`
			PgScan              int64 `json:"pgScan"`
			PgSteal             int64 `json:"pgSteal"`
			Shmem               int64 `json:"shmem"`
			Slab                int64 `json:"slab"`
			SlabReclaimable     int64 `json:"slabReclaimable"`
			SlabUnreclaimable   int64 `json:"slabUnreclaimable"`
			Sock                int64 `json:"sock"`
			ThpCollapseAlloc    int64 `json:"thpCollapseAlloc"`
			ThpFaultAlloc       int64 `json:"thpFaultAlloc"`
			Unevictable         int64 `json:"unevictable"`
			WorkingSetActivate  int64 `json:"workingsetActivate"`
			WorkingSetNoReclaim int64 `json:"workingsetNodereclaim"`
			WorkingSetRefault   int64 `json:"workingsetRefault"`
		} `json:"stats"`
		Limit int64 `json:"limit"`
	} `json:"memoryStats"`
	Networks struct {
		Eth0 struct {
			RxBytes   int64 `json:"rxBytes"`
			RxPackets int64 `json:"rxPackets"`
			RxErrors  int64 `json:"rxErrors"`
			RxDropped int64 `json:"rxDropped"`
			TxBytes   int64 `json:"txBytes"`
			TxPackets int64 `json:"txPackets"`
			TxErrors  int64 `json:"txErrors"`
			TxDropped int64 `json:"txDropped"`
		} `json:"eth0"`
	} `json:"networks"`
}

func (stats *ContainerStats) FormatCpuUsagePercentage() string {
	// Calculate the total CPU time used by the container
	totalCPUUsage := float64(stats.CpuStats.CpuUsage.TotalUsage - stats.PreCPUStats.CpuUsage.TotalUsage)

	// Calculate the system CPU time
	systemCPUUsage := float64(stats.CpuStats.SystemCPUUsage - stats.PreCPUStats.SystemCPUUsage)

	// Calculate the number of online CPUs
	onlineCPUs := float64(stats.CpuStats.OnlineCPUs)

	// Calculate the CPU usage percentage
	cpuUsagePercentage := (totalCPUUsage / systemCPUUsage) * onlineCPUs * 100.0
	if math.IsNaN(cpuUsagePercentage) {
		return "0.00%"
	}
	return fmt.Sprintf("%.2f%%", cpuUsagePercentage)
}
func (stats *ContainerStats) FormatMemoryUsage() string {
	// Get the memory usage and limit in bytes
	memoryUsage := stats.MemoryStats.Usage
	memoryLimit := stats.MemoryStats.Limit

	// Convert the memory usage and limit to human-readable strings
	memoryUsageStr := bytesToHumanReadable(memoryUsage)
	memoryLimitStr := bytesToHumanReadable(memoryLimit)

	// Combine the strings and return the result
	return fmt.Sprintf("%s / %s", memoryUsageStr, memoryLimitStr)
}
func (stats *ContainerStats) FormatNetworkIO() string {
	// Get the network I/O values in bytes
	rxBytes := stats.Networks.Eth0.RxBytes
	txBytes := stats.Networks.Eth0.TxBytes

	// Convert the network I/O values to human-readable strings
	rxBytesStr := bytesToHumanReadable(rxBytes)
	txBytesStr := bytesToHumanReadable(txBytes)

	// Combine the strings and return the result
	return fmt.Sprintf("%s / %s", rxBytesStr, txBytesStr)
}
func (stats *ContainerStats) FormatDiskIO() string {
	// Get the disk read/write values in bytes
	readBytes := int64(0)
	writeBytes := int64(0)

	if len(stats.BlkioStats.IoServiceBytesRecursive) >= 2 {
		if readVal, ok := stats.BlkioStats.IoServiceBytesRecursive[0].(float64); ok {
			readBytes = int64(readVal)
		}
		if writeVal, ok := stats.BlkioStats.IoServiceBytesRecursive[1].(float64); ok {
			writeBytes = int64(writeVal)
		}
	}
	// Convert the disk read/write values to human-readable strings
	readBytesStr := bytesToHumanReadable(readBytes)
	writeBytesStr := bytesToHumanReadable(writeBytes)

	// Combine the strings and return the result
	return fmt.Sprintf("%s / %s", readBytesStr, writeBytesStr)
}
func bytesToHumanReadable(bytes int64) string {
	// Define the units and their corresponding values in bytes
	units := []string{"B", "KB", "MB", "GB", "TB", "PB", "EB"}

	// If the value is already in bytes or less than 1KB, return it directly
	if bytes < 1024 {
		return fmt.Sprintf("%d %s", bytes, units[0])
	}

	// Calculate the index to get the appropriate unit from the "units" array
	index := 0
	value := float64(bytes)
	for value >= 1024 && index < len(units)-1 {
		value /= 1024
		index++
	}

	// Format the value with 2 decimal places and return the human-readable string
	return fmt.Sprintf("%.2f %s", value, units[index])
}
