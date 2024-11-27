package vault

import (
	"datavault/cmd/internal"
	"fmt"
	"mime/multipart"
	"path/filepath"
)

// GetGroups returns list of all groups in the vault
func GetGroups() []string {
	records := VaultConfig.Index.SearchAll([]string{"groupId"})
	groups := make(map[string]struct{}, 0)
	for _, record := range records {
		groups[record.Attributes["groupId"]] = struct{}{}
	}

	groupsList := make([]string, 0)
	for group := range groups {
		groupsList = append(groupsList, group)
	}

	return groupsList
}

// PutGroup uploads records into the vault
func PutGroup(groupId string, files []*multipart.FileHeader) error {
	metadata, err := internal.ProcessMultipartFiles(files, groupId, VaultConfig.Root)
	if err != nil {
		return err
	}

	//create records
	for _, meta := range metadata {
		record := internal.Record{
			Id: meta.FileId,
			Attributes: map[string]string{
				"fileId":        meta.FileId,
				"fileName":      meta.FileName,
				"fileExtension": meta.FileExtension,
				"fileType":      meta.FileType,
				"fileSize":      meta.FileSize,
				"receivedTime":  meta.ReceivedTime,
				"groupId":       groupId,
			},
		}
		VaultConfig.Index.Add(record)
	}

	return nil
}

// FilterByGroup returns list of all records in the vault
func FilterByGroup(groupId string) []internal.Record {
	records := VaultConfig.Index.SearchAny(map[string]string{"groupId": groupId})
	return records
}

// DeleteGroup deletes a group and all its records from the vault
func DeleteGroup(groupId string) {
	records := VaultConfig.Index.SearchAny(map[string]string{"groupId": groupId})
	for _, record := range records {
		VaultConfig.Index.Remove(internal.Record{
			Id:         record.Id,
			Attributes: record.Attributes,
		})
	}
}

// FilterByGroupElement returns a record from a group-element pair
func FilterByGroupElement(groupId, elementId string) (internal.Record, error) {
	records := VaultConfig.Index.SearchAny(map[string]string{"groupId": groupId, "fileId": elementId})
	if len(records) == 0 {
		return internal.Record{}, fmt.Errorf("record not found")
	}
	return records[0], nil
}

// GetElement return a file associated with a record
func GetElement(recordId string) (string, error) {

	attributes := VaultConfig.Index.GetAttributes(recordId)

	dirId := attributes["groupId"]
	fileId := attributes["fileId"] + attributes["fileExtension"]

	path, err := filepath.Abs(filepath.Join(VaultConfig.Root, dirId, fileId))
	if err != nil {
		return "", err
	}

	return path, nil
}

// DeleteElement deletes a record from the vault
func DeleteElement(recordId string) error {
	record := VaultConfig.Index.Get(recordId)
	if record.Id == "" {
		return fmt.Errorf("record not found")
	}

	VaultConfig.Index.Remove(record)

	dirId := record.Attributes["groupId"]
	fileId := record.Attributes["fileId"] + record.Attributes["fileExtension"]
	metafileId := record.Attributes["fileId"] + "._meta"

	err := internal.DeleteFile(VaultConfig.Root, dirId, fileId)
	if err != nil {
		return err
	}

	err = internal.DeleteFile(VaultConfig.Root, dirId, metafileId)
	if err != nil {
		return err
	}

	return nil
}
