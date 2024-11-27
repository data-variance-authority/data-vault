package vault

import (
	"datavault/cmd/internal"
	"encoding/json"
	"net/http"
	"regexp"
)

// HandlerPing is a simple health check endpoint
func HandlerPing(w http.ResponseWriter, r *http.Request) {
	response := map[string]string{
		"id":       VaultConfig.Id,
		"instance": "vault",
		"extended": "",
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	err := json.NewEncoder(w).Encode(response)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// HandlerGroups returns a list of groups
func HandlerGroups(w http.ResponseWriter, r *http.Request) {
	groups := GetGroups()
	w.Header().Set("Content-Type", "application/json")
	err := json.NewEncoder(w).Encode(groups)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// HandlerGroup returns a list of records in a group
func HandlerGroup(w http.ResponseWriter, r *http.Request) {
	groupId := r.URL.Query().Get("groupId")
	if !validateString(groupId) {
		http.Error(w, "Invalid Group ID", http.StatusBadRequest)
		return
	}

	records := FilterByGroup(groupId)
	w.Header().Set("Content-Type", "application/json")
	err := json.NewEncoder(w).Encode(records)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// HandlerGroupUpload uploads files to a group
func HandlerGroupUpload(w http.ResponseWriter, r *http.Request) {
	groupId := r.URL.Query().Get("groupId")
	if !validateString(groupId) {
		http.Error(w, "Invalid Group ID", http.StatusBadRequest)
		return
	}

	// Create group directory in the vault
	err := internal.CreateDirectoryIfNotExists(VaultConfig.Root, groupId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	r.Body = http.MaxBytesReader(w, r.Body, VaultConfig.MAX_UPLOAD_SIZE)
	defer r.Body.Close()
	if err := r.ParseMultipartForm(VaultConfig.IN_MEMORY_UPLOAD_SIZE); err != nil {
		http.Error(w, "files too large", http.StatusBadRequest)
		r.MultipartForm.RemoveAll()
		return
	}
	defer r.MultipartForm.RemoveAll()

	files := r.MultipartForm.File["files"]
	if len(files) == 0 {
		http.Error(w, "No files uploaded", http.StatusBadRequest)
		return
	}

	err = PutGroup(groupId, files)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

// HandlerGroupDelete deletes a group from the vault
func HandlerGroupDelete(w http.ResponseWriter, r *http.Request) {
	groupId := r.URL.Query().Get("groupId")
	if !validateString(groupId) {
		http.Error(w, "Invalid Group ID", http.StatusBadRequest)
		return
	}

	// Delete group directory from the vault
	err := internal.DeleteDirectory(VaultConfig.Root, groupId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Delete group from the vault index
	DeleteGroup(groupId)
	w.WriteHeader(http.StatusOK)
}

// HandleElementGet retrieves an element from the vault
func HandlerElementGet(w http.ResponseWriter, r *http.Request) {
	//this is for API consistency between gatekeeper and vault
	groupId := r.URL.Query().Get("groupId")
	if !validateString(groupId) {
		http.Error(w, "Invalid Group ID", http.StatusBadRequest)
		return
	}

	recordId := r.URL.Query().Get("elementId")
	if !validateString(recordId) {
		http.Error(w, "Invalid Element ID", http.StatusBadRequest)
		return
	}

	path, err := GetElement(recordId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	http.ServeFile(w, r, path)
}

// HandleElementDelete deletes an element from the vault
func HandleElementDelete(w http.ResponseWriter, r *http.Request) {
	//this is for API consistency between gatekeeper and vault
	groupId := r.URL.Query().Get("groupId")
	if !validateString(groupId) {
		http.Error(w, "Invalid Group ID", http.StatusBadRequest)
		return
	}

	recordId := r.URL.Query().Get("elementId")
	if !validateString(recordId) {
		http.Error(w, "Invalid Element ID", http.StatusBadRequest)
		return
	}

	err := DeleteElement(recordId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

// validateString checks if the string is alphanumeric, underscore and hyphen
func validateString(x string) bool {
	re := regexp.MustCompile(`^[a-zA-Z0-9_\-]+$`) // only allow alphanumeric, underscore and hyphen
	return re.MatchString(x)
}
