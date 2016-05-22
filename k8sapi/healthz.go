package k8sapi

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

const K8sResourceHealthz string = "healthz"

type HealthzStatsRetriever struct{}

func (stats HealthzStatsRetriever) Stats() (r string, err error) {
	var response *http.Response
	if response, err = request(K8sResourceHealthz); err != nil {
		stats.Log(fmt.Sprintf("Error: %s", err))
		return
	}

	defer response.Body.Close()

	var rawResponse []byte
	if rawResponse, err = ioutil.ReadAll(response.Body); err != nil {
		stats.Log(fmt.Sprintf("Error: %s", err))
		return
	}

	return string(rawResponse), err
}

func (stats HealthzStatsRetriever) Log(message string) {
	log.Println(message)
}
