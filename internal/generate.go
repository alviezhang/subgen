package internal

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/patrickmn/go-cache"
	"gopkg.in/yaml.v3"
)

var _cache = cache.New(5*time.Minute, 10*time.Minute)
var _lock sync.Mutex

func loadConfig(path string) (*Config, error) {
	yamlFile, err := os.ReadFile(path)
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

func generate(config Config, subtype string) string {
	_lock.Lock()
	// 确保在函数退出时解锁
	defer _lock.Unlock()

	if data, found := _cache.Get(subtype); found {
		fmt.Println("Use cache")
		return data.(string)
	}
	nnr := NNR{token: "7b38f174-97a5-4e36-a50b-c82931b8a977"}
	nodes, err := nnr.Generate(config.Nodes)

	if err != nil {
		panic("Failed generate nodes")
	}
	var generated_config string
	switch subtype {
	case "v2ray":
		generated_config = generateV2Ray(config.ClientConfig, nodes)
	case "quantumult":
		generated_config = generateQuantumult(config.ClientConfig, nodes)
	case "clash":
		generated_config = generateClash(config.ClientConfig, nodes)
	default:
		panic("Unexpected subtype: " + subtype)
	}
	_cache.Set(subtype, generated_config, cache.DefaultExpiration)
	return generated_config
}

func generateV2Ray(clientConfig ClientConfig, nodes []Node) string {
	var encodedNodes []string
	for _, node := range nodes {
		config := map[string]string{
			"v":    "2",
			"ps":   node.Name,
			"add":  node.Addr,
			"port": fmt.Sprint(node.Port),
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
		line := fmt.Sprintf("vmess=%s:%d, method=chacha20-ietf-poly1305, password=%s, fast-open=false, udp-relay=true, tag=%s",
			node.Addr, node.Port, clientConfig.UUID, node.Name)
		lines = append(lines, line)
	}
	return strings.Join(lines, "\n")
}

func generateClash(clientConfig ClientConfig, nodes []Node) string {
	var lines []string
	for _, node := range nodes {
		proxy := map[string]string{
			"name":    node.Name,
			"type":    "vmess",
			"server":  node.Addr,
			"port":    fmt.Sprint(node.Port),
			"uuid":    clientConfig.UUID,
			"alterId": "0",
			"cipher":  "auto",
		}
		data, _ := json.Marshal(proxy)
		lines = append(lines, "  - "+string(data))
	}
	return "proxies:\n" + strings.Join(lines, "\n")
}

func Generate(filename string, subtype string) (string, error) {
	if subtype != "v2ray" && subtype != "quantumult" && subtype != "clash" {
		return "", fmt.Errorf("invalid subtype: " + subtype)
	}
	if filename == "" {
		return "", fmt.Errorf("invalid filename: " + filename)
	}

	config, err := loadConfig(filename)
	if err != nil {
		return "", fmt.Errorf("read file %s\n error: %s", filename, err)
	}

	return generate(*config, subtype), nil
}
