package main

import (
	"os"
	"path"
)

func fileExistance(fileRequest string) bool {
	filePath := path.Clean(fileRequest)
	_, err := os.Stat(filePath)

	if err != nil {
		if err == os.ErrNotExist {
			return false
		}
		return false
	}
	return true
}

func fileHandler(fileRequest string) (int, []byte) {
	filePath := path.Clean(fileRequest)
	file, err := os.Open(filePath)

	if err != nil {
		os.Exit(1)
	}

	buffer := make([]byte, 2048)

	n, err := file.Read(buffer)

	if err != nil {
		os.Exit(1)
	}

	return n, buffer[:n]
}
