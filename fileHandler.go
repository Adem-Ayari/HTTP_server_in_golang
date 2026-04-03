package main

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"log"
	"os"
	"path"
)

func fileExistence(dir string, fileRequest string) bool {
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

func fileExistenceType(extension string) string {
	mimeTypes := map[string]string{
		".html": "text/html",
		".css":  "text/css",
		".js":   "application/javascript",
		".json": "application/json",
		".txt":  "text/plain",
		".xml":  "application/xml",
		".jpg":  "image/jpeg",
		".jpeg": "image/jpeg",
		".png":  "image/png",
		".gif":  "image/gif",
		".pdf":  "application/pdf",
		".mp4":  "video/*",
	}

	if mime, ok := mimeTypes[extension]; ok {
		return mime
	}
	return "text/plain"
}

func pathExistence(dir string, pathRequest string) bool {
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

func encodedGzip(n int, buffer []byte) (int, []byte) {

	var b bytes.Buffer
	gz := gzip.NewWriter(&b)
	_, err := gz.Write(buffer[:n])

	if err != nil {
		fmt.Printf("error%s", err)
	}
	if err := gz.Flush(); err != nil {
		fmt.Printf("error%s", err)
	}
	if err := gz.Close(); err != nil {
		fmt.Printf("error%s", err)
	}

	return b.Len(), b.Bytes()
}
