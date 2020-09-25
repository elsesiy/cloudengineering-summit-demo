package main

import (
	"io/ioutil"
	"net/http"
)

// GetCurrentIP resolves the clients current public IP for the purpose of whitelisting it on the Kubernetes API server
func GetCurrentIP() string {
	resp, err := http.Get("https://ifconfig.me")
	if err != nil {
		return ""
	}

	defer func() { _ = resp.Body.Close() }()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return ""
	}

	return string(body)
}
