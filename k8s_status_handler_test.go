package main

import (
	"errors"
	"fmt"
	"testing"
)

type MockStatsRetriever struct {
	ResponseBody string
	Error        error
}

var lastLogEntry string

func (sr MockStatsRetriever) Stats() (string, error) {
	return sr.ResponseBody, sr.Error
}

func (sr MockStatsRetriever) Log(message string) {
	lastLogEntry = message
}

func TestGetStats(t *testing.T) {
	cases := []struct {
		sr               MockStatsRetriever
		expectedResponse string
		expectedError    error
	}{
		{
			sr: MockStatsRetriever{
				ResponseBody: "ok",
				Error:        nil,
			},
			expectedResponse: "{\"status\":{\"healthy\":true,\"message\":\"Kubelet service is healthy\"}}",
		},
		{
			sr: MockStatsRetriever{
				ResponseBody: "notOkOrOtherResponseBody",
				Error:        nil,
			},
			expectedResponse: "{\"status\":{\"healthy\":false,\"message\":\"Kubelet service is unhealthy\"}}",
		},
		{
			sr: MockStatsRetriever{
				ResponseBody: "ok",
				Error:        errors.New("Error! The response body doesn't matter in this case"),
			},
			expectedResponse: "{\"status\":{\"healthy\":false,\"message\":\"Kubelet service is unhealthy\"}}",
		},
		{
			sr: MockStatsRetriever{
				ResponseBody: "",
				Error:        nil,
			},
			expectedResponse: "{\"status\":{\"healthy\":false,\"message\":\"Kubelet service is unhealthy\"}}",
		},
		{
			sr: MockStatsRetriever{
				ResponseBody: "",
				Error:        errors.New("Kubernetes didn't return response"),
			},
			expectedResponse: "{\"status\":{\"healthy\":false,\"message\":\"Kubelet service is unhealthy\"}}",
		},
	}

	for _, c := range cases {
		response := GetStats(c.sr)

		if response != c.expectedResponse {
			t.Logf("Expected and actual responses are not equal")
			t.Logf("\nWant:\n\t%s\nGot:\n\t%s\n", c.expectedResponse, response)
			t.Fail()
		}

		//Check if errors are logged
		if c.sr.Error != nil {
			if fmt.Sprintf("Error: %s", c.sr.Error) != lastLogEntry {
				t.Log("Last log entry and error message are not equal")
				t.Logf("\nWant:\n\t%s\nGot:\n\t%s\n", fmt.Sprintf("Error: %s", c.sr.Error), lastLogEntry)
				t.Fail()
			}
		}
	}
}
