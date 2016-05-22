package main

import (
	"encoding/json"
	"testing"
)

type MockInternetChecker struct {
	ResponseFromServices bool
}

func (this MockInternetChecker) Check(string) bool {
	return this.ResponseFromServices
}

func TestCheckInternet(t *testing.T) {
	var cases = []struct {
		Mock                   MockInternetChecker
		ExpectedHealthStatuses bool
	}{
		{
			Mock: MockInternetChecker{
				ResponseFromServices: true,
			},
			ExpectedHealthStatuses: true,
		},
		{
			Mock: MockInternetChecker{
				ResponseFromServices: false,
			},
			ExpectedHealthStatuses: false,
		},
	}

	for _, c := range cases {
		response := checkInternet(c.Mock)

		var responseObject InternetHealth

		json.Unmarshal([]byte(response), &responseObject)

		for _, item := range responseObject.Items {
			if item.Status.Healthy != c.ExpectedHealthStatuses {
				t.Log("Expected and actual health status mismatch")
				t.Fail()
			}
		}
	}
}
