package drive

import (
	"io"
	"log"
	"os"
)

// Save content to  "./tmp" or "/tmp"
func SaveContent(content io.ReadCloser) (*os.File, error) {
	defer content.Close()

	// Check if it is in lambda
	tmpDir := "/tmp"
	if _, err := os.Stat(tmpDir); os.IsNotExist(err) {
		// If not in lambda save to ./tmp
		tmpDir = "./tmp"
		if _, err := os.Stat(tmpDir); os.IsNotExist(err) {
			// If ./tmp is not exist, create folder
			if err := os.Mkdir(tmpDir, os.ModePerm); err != nil {
				return nil, err
			}
		}
	}

	file, err := os.CreateTemp(tmpDir, "")
	if err != nil {
		return nil, err
	}
	_, err = io.Copy(file, content)
	if err != nil {
		return nil, err
	}
	_, err = file.Seek(0, 0)
	if err != nil {
		return nil, err
	}
	log.Printf("Saved %s", file.Name())
	return file, nil
}
