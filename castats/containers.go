package castats

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
)

const containersEndpoint string = "containers"

type ContainersStatistics struct {
	Name  string           `json:"name"`
	Stats []ContainersStat `json:"stats"`
}

type ContainersStatCpuUsage struct {
	Total uint64 `json:"total"`
}

type ContainersStatCpu struct {
	Usage ContainersStatCpuUsage `json:"usage"`
}

type ContainersStatMemory struct {
	Usage uint64 `json:"usage"`
}

type ContainersStatFilesystem struct {
	Device   string `json:"device"`
	Capacity uint64 `json:"capacity"`
	Usage    uint64 `json:"usage"`
}

type ContainersStat struct {
	Timestamp  uint64                     `json:"timestamp"`
	Cpu        ContainersStatCpu          `json:"cpu"`
	Memory     ContainersStatMemory       `json:"memory"`
	Filesystem []ContainersStatFilesystem `json:"filesystem"`
}

func (this CAMetricsRetriever) Containers() (statistics ContainersStatistics, err error) {
	var r *http.Response
	var body []byte

	if r, err = request(containersEndpoint); err != nil {
		return
	}

	defer r.Body.Close()

	if body, err = ioutil.ReadAll(r.Body); err != nil {
		return
	}

	err = json.Unmarshal(body, &statistics)

	return
}
