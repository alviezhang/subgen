package internal

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"strings"

	"gopkg.in/yaml.v3"
)

type Node struct {
	Name    string `json:"ps"`
	Addr    string `json:"add"`
	Port    string `json:"port"`
	Region  string `json:"-"`
	IsRelay bool   `json:"-"`
}

type ClientConfig struct {
	UUID string `yaml:"uuid"`
}

type NodeConfig struct {
	Name     string `yaml:"name"`
	Endpoint string `yaml:"endpoint"`
	Region   string `yaml:"region"`
	Relays   []struct {
		Name     string `yaml:"name"`
		Endpoint string `yaml:"endpoint"`
	} `yaml:"relays"`
}

type Config struct {
	ClientConfig ClientConfig `yaml:"client_config"`
	Nodes        []NodeConfig `yaml:"nodes"`
}

func loadConfig(path string, subType string) (*Config, error) {
	yamlFile, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var config Config
	err = yaml.Unmarshal(yamlFile, &config)

	if err != nil {
		return nil, err
	}

	return &config, nil
}

func parseEndpoint(endpoint string) (string, string) {
	parts := strings.Split(endpoint, ":")
	return parts[0], parts[1]
}

func generate(config Config, subtype string) string {
	var nodeList []Node
	for _, nodeConfig := range config.Nodes {
		nodeName := nodeConfig.Name
		addr, port := parseEndpoint(nodeConfig.Endpoint)
		nodeList = append(nodeList, Node{
			Name:    nodeName,
			Addr:    addr,
			Port:    port,
			Region:  nodeConfig.Region,
			IsRelay: false,
		})

		for _, relayConfig := range nodeConfig.Relays {
			addr, port := parseEndpoint(relayConfig.Endpoint)
			nodeList = append(nodeList, Node{
				Name:    fmt.Sprintf("%s-%s", nodeName, relayConfig.Name),
				Addr:    addr,
				Port:    port,
				Region:  "",
				IsRelay: true,
			})
		}
	}
	switch subtype {
	case "v2ray":
		return generateV2Ray(config.ClientConfig, nodeList)
	case "quantumult":
		return generateQuantumult(config.ClientConfig, nodeList)
	case "clash":
		return generateClash(config.ClientConfig, nodeList)
	}
	panic("Unexpected subtype: " + subtype)
}

func generateV2Ray(clientConfig ClientConfig, nodes []Node) string {
	var encodedNodes []string
	for _, node := range nodes {
		config := map[string]string{
			"v":    "2",
			"ps":   node.Name,
			"add":  node.Addr,
			"port": node.Port,
			"id":   clientConfig.UUID,
			// "aid":  "0",
			// "scy":  "auto",
			"net":  "tcp",
			"type": "none",
		}
		data, _ := json.Marshal(config)
		encoded := base64.StdEncoding.EncodeToString(data)
		encodedNodes = append(encodedNodes, "vmess://"+encoded)
	}
	return strings.Join(encodedNodes, "\n")
}

func generateQuantumult(clientConfig ClientConfig, nodes []Node) string {
	var lines []string
	for _, node := range nodes {
		line := fmt.Sprintf("vmess=%s:%s, method=chacha20-ietf-poly1305, password=%s, fast-open=false, udp-relay=true, tag=%s",
			node.Addr, node.Port, clientConfig.UUID, node.Name)
		lines = append(lines, line)
	}
	return strings.Join(lines, "\n")
}

func generateClash(clientConfig ClientConfig, nodes []Node) string {
	var lines []string
	for _, node := range nodes {
		proxy := map[string]string{
			"name":   node.Name,
			"type":   "vmess",
			"server": node.Addr,
			"port":   node.Port,
			"uuid":   clientConfig.UUID,
			// "alterId": "0",
			// "cipher": "auto",
		}
		data, _ := json.Marshal(proxy)
		lines = append(lines, "  - "+string(data))
	}
	return "proxies:\n" + strings.Join(lines, "\n")
}

func Generate(filename string, subtype string) (string, error) {
	if subtype != "v2ray" && subtype != "quantumult" && subtype != "clash" {
		return "", errors.New("Invalid subtype: " + subtype)
	}
	if filename == "" {
		return "", errors.New("Invalid filename: " + filename)
	}

	config, err := loadConfig(filename, subtype)
	if err != nil {
		return "", errors.New(fmt.Sprintf("Read file %s\n error: %s", filename))
	}

	return generate(*config, subtype), nil
}
