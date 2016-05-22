package k8sapi

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
)

const k8sComponentStatusesPath string = "api/v1/componentstatuses"

type ComponentStatusesItemMetadata struct {
	Name string `json:"name"`
}

type ComponentStatusesItemCondition struct {
	Type    string `json:"type"`
	Status  string `json:"status"`
	Message string `json:"message"`
}

type ComponentStatusesItem struct {
	Metadata   ComponentStatusesItemMetadata    `json:"metadata"`
	Conditions []ComponentStatusesItemCondition `json:"conditions"`
}

type ComponentStatuses struct {
	ApiVersion string                  `json:"apiVersion"`
	Items      []ComponentStatusesItem `json:"items"`
}

func (this K8sStatsRetriever) Componentstatues() (c ComponentStatuses, err error) {
	var r *http.Response
	var body []byte

	if r, err = request(k8sComponentStatusesPath); err != nil {
		return
	}

	defer r.Body.Close()

	if body, err = ioutil.ReadAll(r.Body); err != nil {
		return
	}

	err = json.Unmarshal(body, &c)

	return
}
