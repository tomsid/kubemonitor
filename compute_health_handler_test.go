package main

import (
	"errors"
	"github.com/tomsid/kubemonitor/k8sapi"
	"testing"
)

type MockComputeHealthMetricsRetriever struct {
	NodesMock                  k8sapi.NodesStats
	NodesMockError             error
	ComponentstatusesMock      k8sapi.ComponentStatuses
	ComponentstatusesMockError error
}

func (this MockComputeHealthMetricsRetriever) Nodes() (c k8sapi.NodesStats, err error) {
	return this.NodesMock, this.NodesMockError
}

func (this MockComputeHealthMetricsRetriever) Componentstatues() (c k8sapi.ComponentStatuses, err error) {
	return this.ComponentstatusesMock, this.ComponentstatusesMockError
}

func (this MockComputeHealthMetricsRetriever) Log(message string) {
	lastLogEntry = message
}

func (this MockComputeHealthMetricsRetriever) Timestamp() int {
	return 100000000
}

func TestComputeHealth(t *testing.T) {
	cases := []struct {
		Mock                 MockComputeHealthMetricsRetriever
		ExpectedResponse     string
		ExpectedLastLogEntry string
	}{
		{
			Mock: MockComputeHealthMetricsRetriever{
				NodesMock: k8sapi.NodesStats{
					ApiVersion: "v1",
					Items: []k8sapi.Node{
						k8sapi.Node{
							Metadata: k8sapi.NodeMetadata{
								Name: "Some Name",
								Uid:  "UID",
							},
							Status: k8sapi.NodeStatus{
								Conditions: []k8sapi.NodeStatusCondition{
									k8sapi.NodeStatusCondition{
										Type:    "Some type",
										Status:  "Healthy",
										Reason:  "Component works well",
										Message: "Component status is healthy",
									},
								},
							},
						},
					},
				},
				NodesMockError: nil,
				ComponentstatusesMock: k8sapi.ComponentStatuses{
					ApiVersion: "v1",
					Items: []k8sapi.ComponentStatusesItem{
						k8sapi.ComponentStatusesItem{
							Metadata: k8sapi.ComponentStatusesItemMetadata{
								Name: "Metadata Name",
							},
							Conditions: []k8sapi.ComponentStatusesItemCondition{
								k8sapi.ComponentStatusesItemCondition{
									Type:    "Item Condition Type",
									Status:  "Test Status",
									Message: "OK",
								},
							},
						},
					},
				},
				ComponentstatusesMockError: nil,
			},
			ExpectedResponse:     "[{\"kind\":\"k8s\",\"name\":\"Node Some Name\",\"timestamp\":100000000,\"metadata\":{\"uid\":\"UID\",\"resourceVersion\":\"v1\"},\"status\":{\"healthy\":false,\"message\":\"\"}},{\"kind\":\"k8s\",\"name\":\"Metadata Name\",\"timestamp\":100000000,\"metadata\":{\"uid\":\"\",\"resourceVersion\":\"v1\"},\"status\":{\"healthy\":false,\"message\":\"\"}}]",
			ExpectedLastLogEntry: "",
		},
		{
			Mock: MockComputeHealthMetricsRetriever{
				NodesMock:      k8sapi.NodesStats{},
				NodesMockError: errors.New("Some error"),
				ComponentstatusesMock: k8sapi.ComponentStatuses{
					ApiVersion: "v1",
					Items: []k8sapi.ComponentStatusesItem{
						k8sapi.ComponentStatusesItem{
							Metadata: k8sapi.ComponentStatusesItemMetadata{
								Name: "Metadata Name",
							},
							Conditions: []k8sapi.ComponentStatusesItemCondition{
								k8sapi.ComponentStatusesItemCondition{
									Type:    "Item Condition Type",
									Status:  "Test Status",
									Message: "OK",
								},
							},
						},
					},
				},
				ComponentstatusesMockError: nil,
			},
			ExpectedResponse:     "[{\"kind\":\"k8s\",\"name\":\"Metadata Name\",\"timestamp\":100000000,\"metadata\":{\"uid\":\"\",\"resourceVersion\":\"v1\"},\"status\":{\"healthy\":false,\"message\":\"\"}}]",
			ExpectedLastLogEntry: "Unable to get nodes status. Skipping...",
		},
		{
			Mock: MockComputeHealthMetricsRetriever{
				NodesMock: k8sapi.NodesStats{
					ApiVersion: "v1",
					Items: []k8sapi.Node{
						k8sapi.Node{
							Metadata: k8sapi.NodeMetadata{
								Name: "Some Name",
								Uid:  "UID",
							},
							Status: k8sapi.NodeStatus{
								Conditions: []k8sapi.NodeStatusCondition{
									k8sapi.NodeStatusCondition{
										Type:    "Some type",
										Status:  "Healthy",
										Reason:  "Component works well",
										Message: "Component status is healthy",
									},
								},
							},
						},
					},
				},
				NodesMockError:             nil,
				ComponentstatusesMock:      k8sapi.ComponentStatuses{},
				ComponentstatusesMockError: errors.New("Error while getting components statuses"),
			},
			ExpectedResponse:     "[{\"kind\":\"k8s\",\"name\":\"Node Some Name\",\"timestamp\":100000000,\"metadata\":{\"uid\":\"UID\",\"resourceVersion\":\"v1\"},\"status\":{\"healthy\":false,\"message\":\"\"}}]",
			ExpectedLastLogEntry: "Unable to get components status. Skipping...",
		},
		{
			Mock: MockComputeHealthMetricsRetriever{
				NodesMock:                  k8sapi.NodesStats{},
				NodesMockError:             errors.New("Error"),
				ComponentstatusesMock:      k8sapi.ComponentStatuses{},
				ComponentstatusesMockError: errors.New("An error too"),
			},
			ExpectedResponse:     "[]",
			ExpectedLastLogEntry: "Can't get any stats",
		},
	}

	for _, c := range cases {
		lastLogEntry = ""
		response := ComputeHealth(c.Mock)

		if response != c.ExpectedResponse {
			t.Logf("Response changed!")
			t.Logf("Expect: %s, Got: %s", c.ExpectedResponse, response)
			t.Fail()
		}

		if lastLogEntry != c.ExpectedLastLogEntry {
			t.Log("The last long entry is incorrect")
			t.Logf("Expect: %s, Got: %s", c.ExpectedLastLogEntry, lastLogEntry)
			t.Fail()
		}
	}
}
