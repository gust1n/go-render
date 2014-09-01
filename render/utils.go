package render

import (
	"errors"
	"io/ioutil"
	"os"
	"path/filepath"
)

func generateTemplateName(base, path string) string {
	return filepath.ToSlash(path[len(base)+1:])
}

func file_content(path string) (string, error) {
	// Read the file content of the template
	file, err := os.Open(path)
	if file != nil {
		defer file.Close()
	}
	if err != nil {
		return "", err
	}
	b, err := ioutil.ReadAll(file)
	if err != nil {
		return "", err
	}
	s := string(b)

	if len(s) < 1 {
		return "", errors.New("render: template file is empty")
	}

	return s, nil
}
