package main

import (
	"github.com/tomsid/kubemonitor/k8sapi"
	"net/http"
	"os"
)

const (
	EnvKubernetesHost string = "K8SMONITOR_HOST"
	EnvKubernetesPort string = "K8SMONITOR_PORT"
	ApplicationPort   string = "9091"
)

func init() {
	k8sapi.K8sHost = "localhost"
	k8sapi.K8sPort = "8080"

	if os.Getenv(EnvKubernetesHost) != "" {
		k8sapi.K8sHost = os.Getenv(EnvKubernetesHost)
	}

	if os.Getenv(EnvKubernetesPort) != "" {
		k8sapi.K8sPort = os.Getenv(EnvKubernetesPort)
	}
}

func initHandlers() {
	http.HandleFunc("/v1/status", statusHandler)
	http.HandleFunc("/v1/metrics/internet/health/", checkInternetHandler)
	http.HandleFunc("/v1/metrics/", metricsHandler)
	http.HandleFunc("/v1/metrics/compute/health", metricsComputeHealthHandler)
}

func main() {
	initHandlers()

	if err := http.ListenAndServe(":"+ApplicationPort, nil); err != nil {
		panic(err)
	}
}
