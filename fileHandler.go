package main

import (
	"log"
	"os"
	"path"
)

func fileExistance(dir string, fileRequest string) bool {
	filePath := path.Clean(fileRequest)
	_, err := os.Stat(dir + filePath)

	if err != nil {
		if err == os.ErrNotExist {
			return false
		}
		return false
	}
	return true
}

func pathExistance(dir string, pathRequest string) bool {
	filePath := path.Dir(pathRequest)
	_, err := os.Stat(dir + filePath)

	if err != nil {
		if err == os.ErrPermission {
			log.Fatal(err)
		}
		return false
	}

	return true
}

func fileHandlerGET(dir string, fileRequest string) (int, []byte) {
	filePath := path.Clean(fileRequest)
	file, err := os.Open(dir + filePath)

	if err != nil {
		os.Exit(1)
	}
	defer file.Close()

	buffer := make([]byte, 1048576)

	n, err := file.Read(buffer)

	if err != nil {
		log.Fatal(err)
	}

	return n, buffer[:n]
}

func fileHandlerPOST(dir string, fileRequest string, data string) bool {
	filePath := path.Clean(fileRequest)
	file, err := os.Create(dir + filePath)

	if err != nil {
		return false
	}
	defer file.Close()

	n, err := file.Write([]byte(data))

	if err != nil {
		return false
	} else if n != len(data) {
		os.Remove(dir + filePath)
		return false
	}

	return true
}
