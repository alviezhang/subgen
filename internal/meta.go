package internal

type Config struct {
	ClientConfig ClientConfig `yaml:"client_config"`
	Nodes        []NodeConfig `yaml:"nodes"`
}

type ClientConfig struct {
	UUID string `yaml:"uuid"`
}

type NodeConfig struct {
	Name   string `yaml:"name"`
	Host   string `yaml:"host"`
	Port   uint16 `yaml:"port"`
	Region string `yaml:"region"`
}

type Node struct {
	Name    string `json:"ps"`
	Addr    string `json:"addr"`
	Port    uint16 `json:"port"`
	Region  string `json:"region"`
	IsRelay bool   `json:"is_relay"`
}

func Convert(nodes []NodeConfig) []Node {
	converted := make([]Node, len(nodes))
	for index, node := range nodes {
		converted[index] = Node{
			Name:    node.Name,
			Addr:    node.Host,
			Port:    node.Port,
			Region:  node.Region,
			IsRelay: false,
		}
	}
	return converted
}
