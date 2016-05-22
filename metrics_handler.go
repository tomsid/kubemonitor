package main

import (
	"encoding/json"
	"fmt"
	"github.com/tomsid/kubemonitor/castats"
	"net/http"
)

type Filesystem struct {
	Device   string `json:"device"`
	Capacity uint64 `json:"capacity"`
	Used     uint64 `json:"used"`
}

type NetworkDevice struct {
	Name    string `json:"name"`
	Enabled bool   `json:"enabled"`
}

type Metrics struct {
	SystemUUID     string          `json:"system_uuid"`
	NumCores       int8            `json:"num_cores"`
	CpuUsage       uint64          `json:"cpu_usage"`
	MemoryCapacity uint64          `json:"memory_capacity"`
	MemoryFree     uint64          `json:"memory_free"`
	Filesystems    []Filesystem    `json:"filesystems"`
	NetworkDevices []NetworkDevice `json:"network_devices"`
}

type MetricsRetriever interface {
	Containers() (castats.ContainersStatistics, error)
	Machine() (castats.MachineStatistics, error)
	Log(string)
}

func GetMetrics(m MetricsRetriever) string {
	var metrics Metrics

	containersStats, err := m.Containers()

	if err != nil {
		errMsg := fmt.Sprintf("Error while getting containers info. Error: %s", err)
		m.Log(errMsg)
		return errMsg
	} else if len(containersStats.Stats) <= 0 {
		errMsg := "Error: No containers to get info from."
		m.Log(errMsg)
		return errMsg
	}

	machineStats, err := m.Machine()

	if err != nil {
		errMsg := fmt.Sprintf("Error while getting machine info. Error: %s", err)
		m.Log(errMsg)
		return errMsg
	}

	stats := containersStats.Stats[len(containersStats.Stats)-1]

	var filesystems = make([]Filesystem, len(stats.Filesystem))

	for index, filesystemInfo := range stats.Filesystem {
		filesystems[index].Device = filesystemInfo.Device
		filesystems[index].Capacity = filesystemInfo.Capacity
		filesystems[index].Used = filesystemInfo.Usage
	}

	var networkDevices = make([]NetworkDevice, len(machineStats.NetworkDevices))

	for index, machineNetworkDevice := range machineStats.NetworkDevices {
		networkDevices[index].Name = machineNetworkDevice.Name
		networkDevices[index].Enabled = true
	}

	metrics = Metrics{
		SystemUUID:     machineStats.SystemUUID,
		NumCores:       machineStats.NumCores,
		CpuUsage:       stats.Cpu.Usage.Total,
		MemoryCapacity: machineStats.MemoryCapacity,
		MemoryFree:     machineStats.MemoryCapacity - stats.Memory.Usage,
		Filesystems:    filesystems,
		NetworkDevices: networkDevices,
	}

	if responseBody, err := json.Marshal(metrics); err == nil {
		return string(responseBody)
	} else {
		return fmt.Sprintf("Error: %s", err)
	}
}

func metricsHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, GetMetrics(castats.CAMetricsRetriever{}))
}
