package k8sapi

import (
	"log"
	"net/http"
	"time"
)

var K8sHost string = "localhost"
var K8sPort string = "8080"

type K8sStatsRetriever struct{}

func request(path string) (resp *http.Response, err error) {
	resp, err = http.Get("http://" + K8sHost + ":" + K8sPort + "/" + path)
	return
}

func (this K8sStatsRetriever) Log(message string) {
	log.Println(message)
}

func (this K8sStatsRetriever) Timestamp() int {
	return int(time.Now().Unix())
}
