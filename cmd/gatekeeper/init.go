package gatekeeper

import (
	"datavault/configs"
	"encoding/json"
	"log"

	"github.com/serialx/hashring"
)

// Config is the configuration for the gatekeeper server
type Config struct {
	Port string `json:"port"` // Port for the vault server

	Vaults           []string `json:"vaults"`            // List of vaults addresses (host:port)
	BroadcastTimeout int      `json:"broadcast_timeout"` // Timeout for broadcast requests in seconds

	Ring *hashring.HashRing // Consistency hash ring
}

var KeeperConfig Config

// Init initializes the gatekeepr server
func Init() {
	//Read configuration
	configBytes := configs.Instance.ConfigFileData
	//Parse configuration
	err := json.Unmarshal(configBytes, &KeeperConfig)
	if err != nil {
		log.Fatalf("Error parsing gatekeeper configuration: %v\n", err)
	}

	//Initialize the hash ring
	nodes := make(map[string]int)
	for _, vault := range KeeperConfig.Vaults {
		nodes[vault] = 1
	}

	KeeperConfig.Ring = hashring.NewWithWeights(nodes)
}
