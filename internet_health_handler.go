package main

import (
	"encoding/json"
	"fmt"
	"github.com/miekg/dns"
	"log"
	"net/http"
	"net/url"
	"sync"
	"time"
)

const DNSServer string = "8.8.8.8"
const InternetHealthStatusHealthy string = "healthy"
const InternetHealthStatusUnhealthy string = "unhealthy"

var urlsToCheck map[string]map[string]string = map[string]map[string]string{
	"google": map[string]string{
		"url":          "https://www.google.com/",
		"service_name": "Google",
	},
	"dockerhub": map[string]string{
		"url":          "https://hub.docker.com/",
		"service_name": "DockerHub",
	},
	"facebook": map[string]string{
		"url":          "https://www.facebook.com/",
		"service_name": "Facebook",
	},
}

type InternetHealthMetadata struct {
	Url string `json:"url"`
}

type InternetHealthStatus struct {
	Healthy bool   `json:"healthy"`
	Message string `json:"message"`
}

type InternetHealthItem struct {
	Kind      string                 `json:"kind"`
	Name      string                 `json:"name"`
	Timestamp int64                  `json:"timestamp"`
	Metadata  InternetHealthMetadata `json:"metadata"`
	Status    InternetHealthStatus   `json:"status"`
}

type InternetHealth struct {
	Items []InternetHealthItem `json:"items"`
}

type InternetChecker interface {
	Check(url string) (isHealthy bool)
}

type DnsInternetChecker struct{}

func (this DnsInternetChecker) Check(host string) (healthy bool) {
	client := dns.Client{}
	message := dns.Msg{}

	message.SetQuestion(host+".", dns.TypeA)
	res, _, _ := client.Exchange(&message, DNSServer+":53")

	return res.Rcode == dns.RcodeSuccess
}

func checkInternet(checker InternetChecker) string {
	var wg sync.WaitGroup
	var healthItems InternetHealth

	wg.Add(len(urlsToCheck))

	var items = make([]InternetHealthItem, 0, len(urlsToCheck))

	for serviceSlug, info := range urlsToCheck {
		go func(serviceSlug string, info map[string]string) {
			defer wg.Done()

			u, err := url.Parse(info["url"])

			if err != nil {
				log.Printf("Unable to parse url: %s", info["url"])
				return
			}

			healthStatus := checker.Check(u.Host)

			var healthStatusMessage string

			if healthStatus {
				healthStatusMessage = InternetHealthStatusHealthy
			} else {
				healthStatusMessage = InternetHealthStatusUnhealthy
			}

			items = append(items, InternetHealthItem{
				Kind:      "internet",
				Name:      serviceSlug,
				Timestamp: time.Now().Unix(),
				Metadata: InternetHealthMetadata{
					Url: info["url"],
				},
				Status: InternetHealthStatus{
					Healthy: healthStatus,
					Message: fmt.Sprintf("Connection to %s is %s", info["service_name"], healthStatusMessage),
				},
			})
		}(serviceSlug, info)
	}

	wg.Wait()

	healthItems.Items = items

	response, _ := json.Marshal(healthItems)

	return string(response)
}

func checkInternetHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, checkInternet(DnsInternetChecker{}))
}
