package docker

import (
	"encoding/json"
	"fmt"
)

type ContainerStats struct {
	Read      string `json:"read"`
	Preread   string `json:"preread"`
	PidsStats struct {
		Current int64       `json:"current"`
		Limit   json.Number `json:"limit"`
	} `json:"pids_stats"`
	BlkioStats struct {
		IoServiceBytesRecursive []interface{} `json:"io_service_bytes_recursive"`
		IoServicedRecursive     []interface{} `json:"io_serviced_recursive"`
		IoQueueRecursive        []interface{} `json:"io_queue_recursive"`
		IoServiceTimeRecursive  []interface{} `json:"io_service_time_recursive"`
		IoWaitTimeRecursive     []interface{} `json:"io_wait_time_recursive"`
		IoMergedRecursive       []interface{} `json:"io_merged_recursive"`
		IoTimeRecursive         []interface{} `json:"io_time_recursive"`
		SectorsRecursive        []interface{} `json:"sectors_recursive"`
	} `json:"blkio_stats"`
	NumProcs     int64    `json:"num_procs"`
	StorageStats struct{} `json:"storage_stats"`
	CpuStats     struct {
		CpuUsage struct {
			TotalUsage        int64 `json:"total_usage"`
			UsageInKernelMode int64 `json:"usage_in_kernelmode"`
			UsageInUserMode   int64 `json:"usage_in_usermode"`
		} `json:"cpu_usage"`
		SystemCPUUsage int64 `json:"system_cpu_usage"`
		OnlineCPUs     int64 `json:"online_cpus"`
		ThrottlingData struct {
			Periods          int64 `json:"periods"`
			ThrottledPeriods int64 `json:"throttled_periods"`
			ThrottledTime    int64 `json:"throttled_time"`
		} `json:"throttling_data"`
	} `json:"cpu_stats"`
	PreCPUStats struct {
		CpuUsage struct {
			TotalUsage        int64 `json:"total_usage"`
			UsageInKernelMode int64 `json:"usage_in_kernelmode"`
			UsageInUserMode   int64 `json:"usage_in_usermode"`
		} `json:"cpu_usage"`
		SystemCPUUsage int64 `json:"system_cpu_usage"`
		OnlineCPUs     int64 `json:"online_cpus"`
		ThrottlingData struct {
			Periods          int64 `json:"periods"`
			ThrottledPeriods int64 `json:"throttled_periods"`
			ThrottledTime    int64 `json:"throttled_time"`
		} `json:"throttling_data"`
	} `json:"precpu_stats"`
	MemoryStats struct {
		Usage int64 `json:"usage"`
		Stats struct {
			ActiveAnon          int64 `json:"active_anon"`
			ActiveFile          int64 `json:"active_file"`
			Anon                int64 `json:"anon"`
			AnonTHP             int64 `json:"anon_thp"`
			File                int64 `json:"file"`
			FileDirty           int64 `json:"file_dirty"`
			FileMapped          int64 `json:"file_mapped"`
			FileWriteBack       int64 `json:"file_writeback"`
			InactiveAnon        int64 `json:"inactive_anon"`
			InactiveFile        int64 `json:"inactive_file"`
			KernelStack         int64 `json:"kernel_stack"`
			PgActivate          int64 `json:"pgactivate"`
			PgDeactivate        int64 `json:"pgdeactivate"`
			PgFault             int64 `json:"pgfault"`
			PgLazyFree          int64 `json:"pglazyfree"`
			PgLazyFreed         int64 `json:"pglazyfreed"`
			PgMajFault          int64 `json:"pgmajfault"`
			PgRefill            int64 `json:"pgrefill"`
			PgScan              int64 `json:"pgscan"`
			PgSteal             int64 `json:"pgsteal"`
			Shmem               int64 `json:"shmem"`
			Slab                int64 `json:"slab"`
			SlabReclaimable     int64 `json:"slab_reclaimable"`
			SlabUnreclaimable   int64 `json:"slab_unreclaimable"`
			Sock                int64 `json:"sock"`
			ThpCollapseAlloc    int64 `json:"thp_collapse_alloc"`
			ThpFaultAlloc       int64 `json:"thp_fault_alloc"`
			Unevictable         int64 `json:"unevictable"`
			WorkingSetActivate  int64 `json:"workingset_activate"`
			WorkingSetNoReclaim int64 `json:"workingset_nodereclaim"`
			WorkingSetRefault   int64 `json:"workingset_refault"`
		} `json:"stats"`
		Limit int64 `json:"limit"`
	} `json:"memory_stats"`
	Name     string `json:"name"`
	ID       string `json:"id"`
	Networks struct {
		Eth0 struct {
			RxBytes   int64 `json:"rx_bytes"`
			RxPackets int64 `json:"rx_packets"`
			RxErrors  int64 `json:"rx_errors"`
			RxDropped int64 `json:"rx_dropped"`
			TxBytes   int64 `json:"tx_bytes"`
			TxPackets int64 `json:"tx_packets"`
			TxErrors  int64 `json:"tx_errors"`
			TxDropped int64 `json:"tx_dropped"`
		} `json:"eth0"`
	} `json:"networks"`
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
