package internal

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

type NNR struct {
	token string
}

type NNRResponse struct {
	Status int         `json:"status"`
	Data   interface{} `json:"data"`
}

type AddRuleParam struct {
	Sid    string `json:"sid"`
	Name   string `json:"name"`
	Remote string `json:"remote"`
	Rport  uint16 `json:"rport"`
	Type   string `json:"type"`
}

type EditRuleParam struct {
	Rid     string      `json:"rid"`
	Name    string      `json:"name"`
	Remote  string      `json:"remote"`
	Rport   uint16      `json:"rport"`
	Type    string      `json:"type"`
	Setting interface{} `json:"setting"`
}

// type ServerResponse struct {
// 	Status int `json:"status"`
// 	Data   []struct {
// 		Sid    string   `json:"sid"`
// 		Type   string   `json:"type"`
// 		Name   string   `json:"name"`
// 		Host   string   `json:"host"`
// 		Min    int      `json:"min"`
// 		Max    int      `json:"max"`
// 		Mf     int      `json:"mf"` // 倍率
// 		Level  int      `json:"level"`
// 		Top    int      `json:"top"` // 排序
// 		Status int      `json:"status"`
// 		Detail string   `json:"detail"`
// 		Types  []string `json:"types"` // "tcp", "udp", "tcp+udp"
// 	} `json:"data"`
// }

// type RuleResponse struct {
// 	Status int `json:"status"`
// 	Data   []struct {
// 		Rid           string `json:"rid"`
// 		Uid           string `json:"uid"`
// 		Sid           string `json:"sid"`
// 		Host          string `json:"host"`
// 		Port          int    `json:"port"`
// 		Remote        string `json:"remote"`
// 		Rport         int    `json:"rport"`
// 		Type          string `json:"type"`
// 		Status        int    `json:"status"`
// 		Name          string `json:"name"`
// 		Traffic       int    `json:"traffic"`
// 		Data          string `json:"data"`
// 		Date          int    `json:"date"`
// 		LoadbalanceID string `json:"loadbalanceId"`
// 		Setting       struct {
// 			ProxyProtocol  int           `json:"proxyProtocol"`
// 			LoadbalanceMod string        `json:"loadbalanceMode"`
// 			Mix0rtt        bool          `json:"mix0rtt"`
// 			Src            interface{}   `json:"src"`
// 			Cfips          []interface{} `json:"cfips"`
// 		} `json:"setting"`
// 	} `json:"data"`
// }

func listServers(token string) (*NNRResponse, error) {
	var data NNRResponse
	client := &http.Client{}
	req, err := http.NewRequest("POST", "https://nnr.moe/api/servers", nil)

	if err != nil {
		return nil, err
	}

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Token", token)

	resp, err := client.Do(req)

	if err != nil {
		return nil, err
	}

	body, _ := io.ReadAll(resp.Body)

	//convert body to interface
	err = json.Unmarshal(body, &data)
	return &data, err
}

func listRules(token string) (*NNRResponse, error) {
	var data NNRResponse
	client := &http.Client{}
	req, err := http.NewRequest("POST", "https://nnr.moe/api/rules", nil)

	if err != nil {
		return nil, err
	}

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Token", token)

	resp, err := client.Do(req)

	if err != nil {
		return nil, err
	}

	body, _ := io.ReadAll(resp.Body)

	//convert body to interface
	err = json.Unmarshal(body, &data)
	return &data, err
}

func addRule(token string, param AddRuleParam) (*NNRResponse, error) {
	var data NNRResponse
	client := &http.Client{}

	paramJson, err := json.Marshal(param)
	if err != nil {
		return nil, err
	}
	paramReader := io.Reader(bytes.NewReader(paramJson))

	req, err := http.NewRequest("POST", "https://nnr.moe/api/rules/add", paramReader)

	if err != nil {
		return nil, err
	}

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Token", token)

	resp, err := client.Do(req)

	if err != nil {
		return nil, err
	}

	body, _ := io.ReadAll(resp.Body)

	//convert body to interface
	err = json.Unmarshal(body, &data)
	return &data, err
}

func editRule(token string, param EditRuleParam) (*NNRResponse, error) {
	var data NNRResponse
	client := &http.Client{}

	paramJson, err := json.Marshal(param)
	if err != nil {
		return nil, err
	}
	paramReader := io.Reader(bytes.NewReader(paramJson))

	req, err := http.NewRequest("POST", "https://nnr.moe/api/rules/edit", paramReader)

	if err != nil {
		return nil, err
	}

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Token", token)

	resp, err := client.Do(req)

	if err != nil {
		return nil, err
	}

	body, _ := io.ReadAll(resp.Body)

	//convert body to interface
	err = json.Unmarshal(body, &data)
	return &data, err
}

func deleteRule(token string, rid string) (*NNRResponse, error) {
	var data NNRResponse
	client := &http.Client{}

	paramJson, err := json.Marshal(struct {
		Rid string `json:"rid"`
	}{
		Rid: rid,
	})
	if err != nil {
		return nil, err
	}
	paramReader := io.Reader(bytes.NewReader(paramJson))

	req, err := http.NewRequest("POST", "https://nnr.moe/api/rules/del", paramReader)

	if err != nil {
		return nil, err
	}

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Token", token)

	resp, err := client.Do(req)

	if err != nil {
		return nil, err
	}

	body, _ := io.ReadAll(resp.Body)

	//convert body to interface
	err = json.Unmarshal(body, &data)
	return &data, err
}

func (nnr *NNR) Generate(nodeConfigList []NodeConfig) ([]Node, error) {
	nodes := Convert(nodeConfigList)

	// get all the servers
	servers, err := listServers(nnr.token)
	if err != nil {
		return nil, err
	}

	// get all the rules
	rules, err := listRules(nnr.token)
	if err != nil {
		return nil, err
	}

	// build a map for rules
	rulesMap := make(map[string]interface{})
	for _, rule := range rules.Data.([]interface{}) {
		sid := rule.(map[string]interface{})["sid"].(string)
		remote := rule.(map[string]interface{})["remote"].(string)
		key := fmt.Sprintf("%s-%s", sid, remote)
		if _, ok := rulesMap[key]; !ok {
			rulesMap[key] = rule
		} else {
			rid := rule.(map[string]interface{})["rid"].(string)
			deleteRule(nnr.token, rid)
			fmt.Printf("rule: %s duplicated deleted success\n", rid)
		}
	}

	// iterate the nodes
	for _, node := range nodes {
		for _, server := range servers.Data.([]interface{}) {
			_server := server.(map[string]interface{})

			// build rule name
			ruleName := fmt.Sprintf("%s-%s-%s", node.Region, node.Name, _server["name"].(string))
			key := fmt.Sprintf("%s-%s", _server["sid"].(string), node.Addr)
			if _, ok := rulesMap[key]; !ok {
				// if the rule is not in the rules, add it
				_, err := addRule(nnr.token, AddRuleParam{
					Sid:    server.(map[string]interface{})["sid"].(string),
					Remote: node.Addr,
					Rport:  uint16(node.Port),
					Type:   "tcp+udp",
					Name:   ruleName,
				})
				if err != nil {
					return nil, err
				}
				fmt.Printf("rule: %s added success\n", ruleName)
			} else {
				// fmt.Printf("rule: %s already exists\n", ruleName)
				onlineRule := rulesMap[key].(map[string]interface{})
				if onlineRule["name"].(string) != ruleName {
					_, err := editRule(nnr.token, EditRuleParam{
						Rid:     onlineRule["rid"].(string),
						Name:    ruleName,
						Remote:  onlineRule["remote"].(string),
						Rport:   uint16(onlineRule["rport"].(float64)),
						Type:    onlineRule["type"].(string),
						Setting: onlineRule["setting"],
					})
					if err != nil {
						return nil, err
					} else {
						fmt.Printf("rule: %s update success\n", ruleName)
					}
				}
				delete(rulesMap, key)
			}
		}
	}

	// get all the rules
	allRules, err := listRules(nnr.token)
	if err != nil {
		return nil, err
	}

	for _, rule := range rulesMap {
		deleteRule(nnr.token, rule.(map[string]interface{})["rid"].(string))
	}

	var outputList []Node

	for _, node := range nodes {
		outputList = append(outputList, Node{
			Name:    fmt.Sprintf("%s-%s-direct", node.Region, node.Name),
			Addr:    node.Addr,
			Port:    node.Port,
			Region:  node.Region,
			IsRelay: false,
		})
	}

	for _, rule := range allRules.Data.([]interface{}) {
		_rule := rule.(map[string]interface{})
		addrs := strings.Split(_rule["host"].(string), ",")
		for index, addr := range addrs {
			var name = _rule["name"].(string)
			if len(addrs) != 1 {
				name = fmt.Sprintf("%s %c", name, 'A'+index)
			}
			outputList = append(outputList, Node{
				Name:    name,
				Addr:    addr,
				Port:    uint16(_rule["port"].(float64)),
				Region:  strings.Split(name, "-")[0],
				IsRelay: true,
			})
		}
	}

	return outputList, nil
}
