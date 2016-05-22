package k8sapi

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
)

const k8sNodesPath string = "/api/v1/nodes"

type NodeMetadata struct {
	Name string `json:"name"`
	Uid  string `json:"uid"`
}

type NodeStatusCondition struct {
	Type    string `json:"type"`
	Status  string `json:"status"`
	Reason  string `json:"reason"`
	Message string `json:"message"`
}

type NodeStatus struct {
	Conditions []NodeStatusCondition `json:"conditions"`
}

type Node struct {
	Metadata NodeMetadata `json:"metadata"`
	Status   NodeStatus   `json:"status"`
}

type NodesStats struct {
	ApiVersion string `json:"apiVersion"`
	Items      []Node `json:"items"`
}

func (this K8sStatsRetriever) Nodes() (c NodesStats, err error) {
	var r *http.Response
	var body []byte

	if r, err = request(k8sNodesPath); err != nil {
		return
	}
	defer r.Body.Close()

	if body, err = ioutil.ReadAll(r.Body); err != nil {
		return
	}

	err = json.Unmarshal(body, &c)

	return
}
