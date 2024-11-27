package internal

import (
	"bufio"
	"errors"
	"io"
	"mime"
	"mime/multipart"
	"os"
	"path/filepath"
)

// CreateDirectoryIfNotExists creates a directory if it does not exist
func CreateDirectoryIfNotExists(root, name string) error {
	// Create directory if it does not exist
	_, err := os.Stat(filepath.Join(root, name))
	if errors.Is(err, os.ErrNotExist) {
		err := os.Mkdir(filepath.Join(root, name), 0777)
		if err != nil {
			return err
		}
	}

	return nil
}

// DeleteDirectory removes a directory if it exists
func DeleteDirectory(root, name string) error {
	// Remove directory if it exists with all its contents
	err := os.RemoveAll(filepath.Join(root, name))
	if err != nil {
		return err
	}

	return nil
}

// DeleteFile removes a file if it exists
func DeleteFile(root, dir, name string) error {
	// Remove file if it exists
	err := os.Remove(filepath.Join(root, dir, name))
	if err != nil {
		return err
	}

	return nil
}

// SaveBytesToFile creates a file if it does not exist
func SaveBytesToFile(root, dir, name string, data []byte) error {
	// Create file if it does not exist
	_, err := os.Stat(filepath.Join(root, dir, name))
	if errors.Is(err, os.ErrNotExist) {
		err := os.WriteFile(filepath.Join(root, dir, name), data, 0644)
		if err != nil {
			return err
		}
	}

	return nil
}

// ReadBytesFromFile reads a file and returns its content
func ReadBytesFromFile(root, dir, name string) ([]byte, string, error) {
	// Read file content
	data, err := os.ReadFile(filepath.Join(root, dir, name))
	if err != nil {
		return nil, "", err
	}

	contentType := mime.TypeByExtension(filepath.Ext(name))
	if contentType == "" {
		contentType = "application/octet-stream"
	}

	return data, contentType, nil
}

// SaveMultipartToFile saves a multipart file to disk
func SaveMultipartToFile(root, dir, name string, fileHeader *multipart.FileHeader) error {
	file, err := fileHeader.Open()
	if err != nil {
		return err
	}
	defer file.Close()

	outfile, err := os.Create(filepath.Join(root, dir, name))
	if err != nil {
		return err
	}
	defer outfile.Close()

	bufferedWriter := bufio.NewWriter(outfile)
	defer bufferedWriter.Flush()

	_, err = io.Copy(bufferedWriter, file)
	if err != nil {
		return err
	}

	return nil
}
