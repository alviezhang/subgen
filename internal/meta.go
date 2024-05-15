package internal

import (
	"strconv"
	"strings"
)

type Config struct {
	ClientConfig ClientConfig `yaml:"client_config"`
	Nodes        []NodeConfig `yaml:"nodes"`
}

type ClientConfig struct {
	UUID string `yaml:"uuid"`
}

type NodeConfig struct {
	Name     string `yaml:"name"`
	Endpoint string `yaml:"endpoint"`
	Region   string `yaml:"region"`
}

type Node struct {
	Name    string `json:"ps"`
	Addr    string `json:"addr"`
	Port    uint16 `json:"port"`
	Region  string `json:"region"`
	IsRelay bool   `json:"is_relay"`
}

func parseEndpoint(endpoint string) (string, uint16) {
	parts := strings.Split(endpoint, ":")
	port, err := strconv.Atoi(parts[1])
	if err != nil {
		panic("Invalid port: " + parts[1])
	}
	return parts[0], uint16(port)
}

func Convert(nodes []NodeConfig) []Node {
	converted := make([]Node, len(nodes))
	for index, node := range nodes {
		addr, port := parseEndpoint(node.Endpoint)
		converted[index] = Node{
			Name:    node.Name,
			Addr:    addr,
			Port:    port,
			Region:  node.Region,
			IsRelay: false,
		}
	}
	return converted
}
