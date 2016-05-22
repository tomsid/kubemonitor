package main

import (
	"encoding/json"
	"fmt"
	"github.com/tomsid/kubemonitor/k8sapi"
	"net/http"
)

type MetricsComputeHealthMetadata struct {
	Uid             string `json:"uid"`
	ResourceVersion string `json:"resourceVersion"`
}

type MetricsComputeHealthStatus struct {
	Healthy bool   `json:"healthy"`
	Message string `json:"message"`
}

type MetricsComputeHealth struct {
	Kind      string                       `json:"kind"`
	Name      string                       `json:"name"`
	Timestamp uint32                       `json:"timestamp"`
	Metadata  MetricsComputeHealthMetadata `json:"metadata"`
	Status    MetricsComputeHealthStatus   `json:"status"`
}

type MetricsComputeHealthStats struct {
	Items []MetricsComputeHealth `json:"items"`
}

const NodeConditionReady string = "Ready"
const ComponentConditionTypeHealthy string = "Healthy"

type ComputeHealthMetricsRetriever interface {
	Nodes() (c k8sapi.NodesStats, err error)
	Componentstatues() (c k8sapi.ComponentStatuses, err error)
	Log(string)
	Timestamp() int
}

func ComputeHealth(m ComputeHealthMetricsRetriever) string {
	var response []byte
	var err error
	var statsNodes []MetricsComputeHealth

	if nodesStats, err := m.Nodes(); err != nil {
		m.Log("Unable to get nodes status. Skipping...")
	} else {
		for _, nodeStats := range nodesStats.Items {
			var nodeHealthStatusCondition k8sapi.NodeStatusCondition

			for _, condition := range nodeStats.Status.Conditions {
				if condition.Type == NodeConditionReady {
					nodeHealthStatusCondition = condition
				}
			}

			statsNodes = append(statsNodes, MetricsComputeHealth{
				Kind:      "k8s",
				Name:      "Node " + nodeStats.Metadata.Name,
				Timestamp: uint32(m.Timestamp()),
				Metadata: MetricsComputeHealthMetadata{
					Uid:             nodeStats.Metadata.Uid,
					ResourceVersion: nodesStats.ApiVersion,
				},
				Status: MetricsComputeHealthStatus{
					Healthy: nodeHealthStatusCondition.Status == "True",
					Message: nodeHealthStatusCondition.Message,
				},
			})
		}
	}

	if componentsStats, err := m.Componentstatues(); err != nil {
		m.Log("Unable to get components status. Skipping...")
	} else {
		for _, componentStats := range componentsStats.Items {
			var componentCondition k8sapi.ComponentStatusesItemCondition

			for _, condition := range componentStats.Conditions {
				if condition.Type == ComponentConditionTypeHealthy {
					componentCondition = condition
				}
			}

			statsNodes = append(statsNodes, MetricsComputeHealth{
				Kind:      "k8s",
				Name:      componentStats.Metadata.Name,
				Timestamp: uint32(uint32(m.Timestamp())),
				Metadata: MetricsComputeHealthMetadata{
					Uid:             "",
					ResourceVersion: componentsStats.ApiVersion,
				},
				Status: MetricsComputeHealthStatus{
					Healthy: componentCondition.Status == "True",
					Message: componentCondition.Message,
				},
			})
		}
	}

	if len(statsNodes) == 0 {
		m.Log("Can't get any stats")
		return "[]"
	}

	if response, err = json.Marshal(statsNodes); err != nil {
		return "Unable to marshal struct to json"
	}

	return string(response)
}

func metricsComputeHealthHandler(w http.ResponseWriter, req *http.Request) {
	fmt.Fprint(w, ComputeHealth(k8sapi.K8sStatsRetriever{}))
}
