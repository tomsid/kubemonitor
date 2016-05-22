package castats

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
)

const machineEndpoint string = "machine"

type MachineFilesystem struct {
	Device   string `json:"device"`
	Capacity uint64 `json:"capacity"`
	Type     string `json:"type"`
	Inodes   uint64 `json:"inodes"`
}
type MachineNetworkDevice struct {
	Name string `json:"name"`
}

type MachineStatistics struct {
	NumCores       int8                   `json:"num_cores"`
	MemoryCapacity uint64                 `json:"memory_capacity"`
	SystemUUID     string                 `json:"system_uuid"`
	Filesystems    []MachineFilesystem    `json:"filesystems"`
	NetworkDevices []MachineNetworkDevice `json:"network_devices"`
}

func (this CAMetricsRetriever) Machine() (statistics MachineStatistics, err error) {
	var r *http.Response
	var body []byte

	if r, err = request(machineEndpoint); err != nil {
		return
	}

	defer r.Body.Close()

	if body, err = ioutil.ReadAll(r.Body); err != nil {
		return
	}

	err = json.Unmarshal(body, &statistics)

	return
}
