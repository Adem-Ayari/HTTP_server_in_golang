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
		tmp := strings.Split(headers[i], ":")
		if len(tmp) == 2 {
			requestParse[tmp[0]] = strings.TrimSpace(tmp[1])
		} else if len(tmp) == 1 {
			requestParse[tmp[0]] = ""
		}
	}

	requestParse["request body"] = body

	return requestParse
}
