package main

import (
	"bytes"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
)

type HostPort struct {
	IP   string
	Port int
}

// Retrieves healthy servers from openrvs-registry over HTTP.
func getServersFromRegistry() ([]HostPort, error) {
	var hostports = make([]HostPort, 0)
	resp, err := http.Get(RegistryURL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	lines := bytes.Split(bytes.TrimSuffix(b, []byte{'\n'}), []byte{'\n'})
	for i := 1; i < len(lines); i++ {
		fields := bytes.Split(lines[i], []byte{','})
		host := string(fields[1])
		portBytes := fields[2]
		port, err := strconv.Atoi(string(portBytes))
		if err != nil {
			log.Println("atoi error:", err)
			continue
		}
		hostports = append(hostports, HostPort{IP: host, Port: port})
	}
	return hostports, nil
}
