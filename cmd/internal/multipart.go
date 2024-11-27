package internal

import (
	"encoding/json"
	"fmt"
	"mime/multipart"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/google/uuid"
)

// Meta is the metadata for a file
type Meta struct {
	FileId        string `json:"fileId"`
	FileType      string `json:"fileType"`
	FileName      string `json:"fileName"`
	FileExtension string `json:"fileExtension"`
	FileSize      string `json:"fileSize"`
	ReceivedTime  string `json:"receivedTime"`
	GroupId       string `json:"groupId"`
}

// ProcessMultipartFiles processes multiple files in parallel
func ProcessMultipartFiles(files []*multipart.FileHeader, groupId, root string) ([]Meta, error) {

	var wg sync.WaitGroup
	wg.Add(len(files))

	var metadata = make([]Meta, len(files))
	errs := make(chan error)
	for i, fileHeader := range files {
		go ProcessFile(fileHeader, groupId, root, &wg, errs, &metadata[i])
	}

	wg.Wait()
	close(errs)

	for err := range errs {
		return nil, err
	}

	return metadata, nil
}

// ProcessFile processes a single file
func ProcessFile(file *multipart.FileHeader, groupId, root string, wg *sync.WaitGroup, errs chan error, responseMeta *Meta) {
	defer wg.Done()

	fileId := strings.ReplaceAll(uuid.New().String(), "-", "")
	filetype := file.Header.Get("Content-Type")
	filename := file.Filename
	filesize := file.Size
	extension := filepath.Ext(filename)

	var metadata Meta
	metadata.FileId = fileId
	metadata.FileType = filetype
	metadata.FileName = filename
	metadata.FileExtension = extension
	metadata.FileSize = fmt.Sprintf("%d", filesize)
	metadata.ReceivedTime = fmt.Sprintf("%d", time.Now().UnixMilli())
	metadata.GroupId = groupId

	// Save file to disk
	metaBytes, err := json.Marshal(metadata)
	if err != nil {
		errs <- err
		return
	}

	err = SaveBytesToFile(root, groupId, fileId+"._meta", metaBytes)
	if err != nil {
		errs <- err
		return
	}

	err = SaveMultipartToFile(root, groupId, fileId+extension, file)
	if err != nil {
		errs <- err
		return
	}

	*responseMeta = metadata
}
