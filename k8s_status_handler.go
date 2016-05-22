package main

import (
	"encoding/json"
	"fmt"
	"github.com/tomsid/kubemonitor/k8sapi"
	"net/http"
)

type HealthStatus struct {
	Healthy bool   `json:"healthy"`
	Message string `json:"message"`
}

type Status struct {
	HealthStatus `json:"status"`
}

type StatsRetriever interface {
	Stats() (response string, err error)
	Log(string)
}

func GetStats(r StatsRetriever) string {
	var isHealthy bool
	var isHealthyMessage string

	response, err := r.Stats()

	if err != nil {
		r.Log(fmt.Sprintf("Error: %s", err))
		isHealthy = false
		isHealthyMessage = "unhealthy"
	} else {

		if string(response) == "ok" {
			isHealthy = true
			isHealthyMessage = "healthy"
		} else {
			isHealthy = false
			isHealthyMessage = "unhealthy"
		}
	}

	responseBody, _ := json.Marshal(
		Status{
			HealthStatus{
				Healthy: isHealthy,
				Message: "Kubelet service is " + isHealthyMessage,
			},
		},
	)

	return string(responseBody)
}

func statusHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, GetStats(k8sapi.HealthzStatsRetriever{}))
}
