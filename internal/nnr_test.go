package internal

import (
	"fmt"
	"testing"
)

// func TestListServersSuccess(t *testing.T) {
// 	token := "7b38f174-97a5-4e36-a50b-c82931b8a977"
// 	resp, err := listServers(token)

// 	if err != nil {
// 		t.Errorf("Error: %v", err)
// 	}
// 	fmt.Println(resp)
// }

// func TestListRuleSuccess(t *testing.T) {
// 	token := "7b38f174-97a5-4e36-a50b-c82931b8a977"
// 	_, err := listRules(token)

// 	if err != nil {
// 		t.Errorf("Error: %v", err)
// 	}
// }

// func TestAddRuleSuccess(t *testing.T) {
// 	token := "7b38f174-97a5-4e36-a50b-c82931b8a977"
// 	resp, err := addRule(token, RuleParam{
// 		Sid:    "6cf4a639-f0dd-42d3-81d7-f47e426264f6",
// 		Remote: "cot.cloud.i.alvie.net",
// 		Rport:  10010,
// 		Type:   "tcp+udp",
// 		Name:   "ttttt",
// 	})

// 	if err != nil {
// 		t.Errorf("Error: %v", err)
// 	}
// 	fmt.Println(resp)
// }

func TestGenerate(t *testing.T) {
	nnr := NNR{token: "7b38f174-97a5-4e36-a50b-c82931b8a977"}

	origin_nodes := []NodeConfig{
		{
			Name:     "ch",
			Endpoint: "ch.cloud.i.alvie.net:10010",
			Region:   "HKG",
		},
		{
			Name:     "coe0",
			Endpoint: "coe.cloud.i.alvie.net:10010",
			Region:   "JPN",
		},
		{
			Name:     "coe1",
			Endpoint: "cot.cloud.i.alvie.net:10010",
			Region:   "JPN",
		},
	}

	generated_nodes, err := nnr.Generate(origin_nodes)
	fmt.Println("generated_nodes:", generated_nodes)
	fmt.Println("err:", err)
}
