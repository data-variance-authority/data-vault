package gatekeeper

import (
	"encoding/json"
	"net/http"
	"net/http/httputil"
	"net/url"
	"sync"
	"time"
)

// HandlerPing is a simple health check endpoint
func HandlerPing(w http.ResponseWriter, r *http.Request) {
	response := map[string]string{
		"id":       "gatekeeper",
		"instance": "gatekeeper",
		"extended": "",
	}
	// Ping every Vault in the configuration
	vaultsNumber := len(KeeperConfig.Vaults)
	vaultsOnline := 0
	vaultsFailed := make([]string, 0)

	// Broadcast the ping request to all vaults
	responses := BroadcastGETRequest("http://", "/ping", KeeperConfig.Vaults)

	// Check the responses
	for i, resp := range responses {
		if resp == nil { // Skip failed requests
			vaultsFailed = append(vaultsFailed, KeeperConfig.Vaults[i])
			continue
		}

		if resp.StatusCode == http.StatusOK {
			vaultsOnline++
		} else {
			vaultsFailed = append(vaultsFailed, KeeperConfig.Vaults[i])
		}
	}

	results := PingResults{
		VaultsNumber: vaultsNumber,
		VaultsOnline: vaultsOnline,
		VaultsFailed: vaultsFailed,
	}

	tbytes, err := json.Marshal(results)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	response["extended"] = string(tbytes)

	// Write the response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// HandlerGroups returns a list of groups from all vaults
func HandlerGroups(w http.ResponseWriter, r *http.Request) {
	// Broadcast the request to all vaults
	responses := BroadcastGETRequest("http://", "/groups", KeeperConfig.Vaults)

	// Check the responses
	groups := make([]string, 0)
	for _, resp := range responses {
		if resp == nil { // Skip failed requests
			continue
		}
		if resp.StatusCode == http.StatusOK {
			var _groups []string
			err := json.NewDecoder(resp.Body).Decode(&_groups)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			groups = append(groups, _groups...)
		}
	}

	// Write the response
	w.Header().Set("Content-Type", "application/json")
	err := json.NewEncoder(w).Encode(groups)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// HandlerGroup returns a list of records in a group
func HandlerGroup(w http.ResponseWriter, r *http.Request) {
	YxorpRequest(w, r, r.URL.Query().Get("groupId"))
}

// HandlerGroupUpload uploads files to a group
func HandlerGroupUpload(w http.ResponseWriter, r *http.Request) {
	YxorpRequest(w, r, r.URL.Query().Get("groupId"))
}

// HandlerGroupDelete deletes a group
func HandlerGroupDelete(w http.ResponseWriter, r *http.Request) {
	YxorpRequest(w, r, r.URL.Query().Get("groupId"))
}

// HandlerElementGet returns a record from a group
func HandlerElementGet(w http.ResponseWriter, r *http.Request) {
	YxorpRequest(w, r, r.URL.Query().Get("groupId"))
}

// HandlerElementUpload uploads a record to a group
func HandleElementDelete(w http.ResponseWriter, r *http.Request) {
	YxorpRequest(w, r, r.URL.Query().Get("groupId"))
}

// PingResults is a struct to store the results of the ping request
type PingResults struct {
	VaultsNumber int      `json:"vaults_number"`
	VaultsOnline int      `json:"vaults_online"`
	VaultsFailed []string `json:"vaults_failed"`
}

// BroadcastGETRequest sends a GET request to multiple addresses
func BroadcastGETRequest(protocol, path string, addresses []string) (responses []*http.Response) {

	responses = make([]*http.Response, len(addresses))

	// Broadcast the request to all vaults using goroutines
	wg := sync.WaitGroup{}
	for i, address := range addresses {
		url := protocol + address + path
		wg.Add(1)
		go func() {
			defer wg.Done()

			client := &http.Client{
				Timeout: time.Duration(KeeperConfig.BroadcastTimeout) * time.Second,
			}
			resp, err := client.Get(url)
			if err != nil {
				return
			}
			responses[i] = resp
		}()
	}

	wg.Wait()

	return responses
}

// YxorpRequest forwards the request to a specific vault
func YxorpRequest(w http.ResponseWriter, r *http.Request, ringNode string) {

	address, ok := KeeperConfig.Ring.GetNode(ringNode)
	if !ok {
		http.Error(w, "group cannot be assigned to a vault", http.StatusBadRequest)
		return
	}

	// Create reverse proxy
	proxy := httputil.NewSingleHostReverseProxy(&url.URL{
		Scheme: "http",
		Host:   address,
	})
	// Serve the request via the reverse proxy
	proxy.ServeHTTP(w, r)
}
