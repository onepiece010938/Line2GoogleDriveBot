package drive

import (
	"io"
	"log"
	"os"
)

// Save content to ./tmp
func SaveContent(content io.ReadCloser) (*os.File, error) {
	defer content.Close()
	file, err := os.CreateTemp("./tmp", "")
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
