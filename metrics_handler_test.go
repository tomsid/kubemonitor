package main

import (
	"errors"
	"github.com/tomsid/kubemonitor/castats"
	"testing"
)

type MockMetricsRetriever struct {
	ContainersMock  castats.ContainersStatistics
	ContainersError error
	MachineMock     castats.MachineStatistics
	MachineError    error
}

func (this MockMetricsRetriever) Containers() (castats.ContainersStatistics, error) {
	return this.ContainersMock, this.ContainersError
}

func (this MockMetricsRetriever) Machine() (castats.MachineStatistics, error) {
	return this.MachineMock, this.MachineError
}

func (this MockMetricsRetriever) Log(message string) {
	lastLogEntry = message
}

func TestGetMetrics(t *testing.T) {
	cases := []struct {
		Mock                 MockMetricsRetriever
		ExpectedResponse     string
		ExpectedLastLogEntry string
	}{
		{
			Mock: MockMetricsRetriever{
				ContainersMock: castats.ContainersStatistics{
					Name: "Some Name",
					Stats: []castats.ContainersStat{
						castats.ContainersStat{
							Timestamp: 143423752938,
							Cpu: castats.ContainersStatCpu{
								Usage: castats.ContainersStatCpuUsage{
									Total: 12432512,
								},
							},
							Memory: castats.ContainersStatMemory{
								Usage: 7287965184,
							},
							Filesystem: []castats.ContainersStatFilesystem{
								castats.ContainersStatFilesystem{
									Device:   "sdb4",
									Capacity: 1212314123,
									Usage:    5324123,
								},
							},
						},
					},
				},
				ContainersError: nil,
				MachineMock: castats.MachineStatistics{
					NumCores:       2,
					MemoryCapacity: 8287965184,
					SystemUUID:     "SIAD-ASDSA-DASDA12-ASDAS",
					Filesystems: []castats.MachineFilesystem{
						castats.MachineFilesystem{
							Device:   "sda1",
							Capacity: 123141423,
							Type:     "vfs",
							Inodes:   123123144,
						},
					},
					NetworkDevices: []castats.MachineNetworkDevice{
						castats.MachineNetworkDevice{
							Name: "vnet02",
						},
					},
				},
				MachineError: nil,
			},
			ExpectedResponse:     "{\"system_uuid\":\"SIAD-ASDSA-DASDA12-ASDAS\",\"num_cores\":2,\"cpu_usage\":12432512,\"memory_capacity\":8287965184,\"memory_free\":1000000000,\"filesystems\":[{\"device\":\"sdb4\",\"capacity\":1212314123,\"used\":5324123}],\"network_devices\":[{\"name\":\"vnet02\",\"enabled\":true}]}",
			ExpectedLastLogEntry: "",
		},
		{
			Mock: MockMetricsRetriever{
				ContainersMock: castats.ContainersStatistics{
					Name: "Test case 2",
					Stats: []castats.ContainersStat{
						castats.ContainersStat{
							Timestamp: 143423752938,
							Cpu: castats.ContainersStatCpu{
								Usage: castats.ContainersStatCpuUsage{
									Total: 12432512,
								},
							},
							Memory: castats.ContainersStatMemory{
								Usage: 7287965184,
							},
							Filesystem: []castats.ContainersStatFilesystem{
								castats.ContainersStatFilesystem{
									Device:   "sdb4",
									Capacity: 1212314123,
									Usage:    5324123,
								},
							},
						},
					},
				},
				ContainersError: nil,
				MachineMock:     castats.MachineStatistics{},
				MachineError:    errors.New("Some error"),
			},
			ExpectedResponse:     "Error while getting machine info. Error: Some error",
			ExpectedLastLogEntry: "Error while getting machine info. Error: Some error",
		},
		{
			Mock: MockMetricsRetriever{
				ContainersMock:  castats.ContainersStatistics{},
				ContainersError: errors.New("A really really bad error"),
				MachineMock: castats.MachineStatistics{
					NumCores:       2,
					MemoryCapacity: 123,
					SystemUUID:     "SIAD-ASDSA-DASDA12-ASDAS",
					Filesystems: []castats.MachineFilesystem{
						castats.MachineFilesystem{
							Device:   "sda1",
							Capacity: 123141423,
							Type:     "vfs",
							Inodes:   123123144,
						},
					},
					NetworkDevices: []castats.MachineNetworkDevice{
						castats.MachineNetworkDevice{
							Name: "vnet02",
						},
					},
				},
				MachineError: nil,
			},
			ExpectedResponse:     "Error while getting containers info. Error: A really really bad error",
			ExpectedLastLogEntry: "Error while getting containers info. Error: A really really bad error",
		},
		{
			Mock: MockMetricsRetriever{
				ContainersMock:  castats.ContainersStatistics{},
				ContainersError: errors.New("Some connection error."),
				MachineMock:     castats.MachineStatistics{},
				MachineError:    nil,
			},
			ExpectedResponse:     "Error while getting containers info. Error: Some connection error.",
			ExpectedLastLogEntry: "Error while getting containers info. Error: Some connection error.",
		},
	}

	for _, c := range cases {
		lastLogEntry = ""
		response := GetMetrics(c.Mock)

		if response != c.ExpectedResponse {
			t.Logf("Response changed! \n Want:\n\t%s,\nGot:\n\t%s", c.ExpectedResponse, response)
			t.Fail()
		}

		if lastLogEntry != c.ExpectedLastLogEntry {
			t.Log("The last long entry is incorrect")
			t.Logf("Expect: %s, Got: %s", c.ExpectedLastLogEntry, lastLogEntry)
			t.Fail()
		}
	}
}
