package main

import (
	"strings"
)

func parser(request string) map[string]string {
	requestParse := make(map[string]string)

	requestArray := strings.Split(request, "\r\n")

	requestParse["request line"] = requestArray[0]

	var pos int

	for i := range requestArray {
		if strings.Compare(requestArray[i], "") == 0 {
			pos = i
		}
	}

	headers := requestArray[1:pos]
	var body string
	if pos+1 < len(requestArray) {
		body = requestArray[pos+1]
	}

	for i := range headers {
		tmpIndex := strings.Index(headers[i], ":")

		if tmpIndex == -1 {
			continue
		}

		name := strings.TrimSpace(headers[i][0:tmpIndex])
		if tmpIndex < len(headers[i])-1 {
			value := strings.TrimSpace(headers[i][tmpIndex+1:])
			requestParse[name] = value
		} else {
			requestParse[name] = ""
		}
	}

	requestParse["request body"] = body

	return requestParse
}
