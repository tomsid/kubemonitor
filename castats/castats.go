package castats

import (
	"log"
	"net/http"
)

const apiEndpoint string = "api/v1.3/"

var Host string = "localhost"
var Port string = "4194"

type CAMetricsRetriever struct{}

func request(path string) (r *http.Response, err error) {
	r, err = http.Get("http://" + Host + ":" + Port + "/" + apiEndpoint + path)
	return
}

func (this CAMetricsRetriever) Log(message string) {
	log.Println(message)
}
