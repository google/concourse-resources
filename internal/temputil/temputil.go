package temputil

import (
	"io/ioutil"
	"log"
	"os"
)

type TempFileManager struct {
	TempDir string
	tempFiles []string
}

func (m *TempFileManager) Create(prefix, fileContents string) (path string, err error) {
	f, err := ioutil.TempFile(m.TempDir, prefix)
	if err != nil {
		return
	}
	defer f.Close()
	path = f.Name()

	m.tempFiles = append(m.tempFiles, path)

	_, err = f.WriteString(fileContents)
	return
}

func (m *TempFileManager) Cleanup() {
	for _, tempFile := range m.tempFiles {
		err := os.Remove(tempFile)
		if err != nil {
			log.Printf("error removing temp file %s: %v", tempFile, err)
		}
	}
	m.tempFiles = []string{}
}
