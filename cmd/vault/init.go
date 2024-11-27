package vault

import (
	"datavault/cmd/internal"
	"datavault/configs"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/google/uuid"
)

// Config represents the configuration for the vault server
type Config struct {
	Id   string `json:"id"`   // Unique identifier for the vault
	Root string `json:"root"` // Root folder for the vault
	Port string `json:"port"` // Port for the vault server

	IN_MEMORY_UPLOAD_SIZE int64 `json:"in_memory_upload_size"` // Maximum size of in-memory upload
	MAX_UPLOAD_SIZE       int64 `json:"max_upload_size"`       // Maximum size of upload

	Index internal.Index // Inverted index for the vault
}

var VaultConfig Config

// Init initializes the vault server
func Init() {
	//Read configuration
	configBytes := configs.Instance.ConfigFileData
	//Parse configuration
	err := json.Unmarshal(configBytes, &VaultConfig)
	if err != nil {
		log.Fatalf("Error parsing vault configuration: %v\n", err)
	}

	//Setup root folder for vault
	if VaultConfig.Root == "" {
		log.Fatalf("Root folder for vault is not set\n")
	}

	//Check if root folder exists, if not create it
	if _, err := os.Stat(VaultConfig.Root); errors.Is(err, os.ErrNotExist) {
		log.Println("Root folder for vault does not exist, creating it...")
		err := os.MkdirAll(VaultConfig.Root, 0777)
		if err != nil {
			log.Fatalf("Error creating root folder for vault: %v\n", err)
		}
	}

	//Initialize inverted index
	VaultConfig.Index, err = generateVaultIndex(VaultConfig.Root)
	if err != nil {
		log.Fatalf("Error reconstructing index: %v\n", err)
	}
}

// generateVaultIndex reconstructs the inverted index from the vault root folder
func generateVaultIndex(root string) (internal.Index, error) {
	index := internal.NewIndex()

	dirs, err := os.ReadDir(root)
	if err != nil {
		return index, err
	}

	for _, dir := range dirs {
		if !dir.IsDir() {
			continue
		}

		groupPath := filepath.Join(root, dir.Name())
		groupDirs, err := os.ReadDir(groupPath)
		if err != nil {
			return index, err
		}

		filesList := make([]string, 0)
		for _, groupFile := range groupDirs {
			if groupFile.IsDir() {
				return index, fmt.Errorf("unexpected directory in group folder: %s", groupFile.Name())
			}

			recordPath := filepath.Join(groupPath, groupFile.Name())
			filesList = append(filesList, recordPath)
		}

		matchedFiles := make(map[string]string, 0)

		for _, recordPath := range filesList {
			if filepath.Ext(recordPath) != "._meta" {
				metaPath := recordPath[:len(recordPath)-len(filepath.Ext(recordPath))] + "._meta"
				if _, err := os.Stat(metaPath); errors.Is(err, os.ErrNotExist) {
					return index, fmt.Errorf("missing meta file for record: %s", recordPath)
				}
				matchedFiles[recordPath] = metaPath
			}
		}

		for _, metaPath := range matchedFiles {
			//create record from reading recordPath filename with extension, and content from metaPath file
			record := internal.Record{}
			metaFile, err := os.ReadFile(metaPath)
			if err != nil {
				return index, err
			}

			var attributes map[string]string
			err = json.Unmarshal(metaFile, &attributes)
			if err != nil {
				return index, err
			}

			record.Id = uuid.New().String()
			record.Attributes = attributes

			index.Add(record)
		}
	}
	return index, err
}
