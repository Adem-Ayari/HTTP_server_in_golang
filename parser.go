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

	headers := requestArray[1 : pos-1]
	body := requestArray[pos]

	for i := range headers {
		tmp := strings.Split(headers[i], ":")
		requestParse[tmp[0]] = strings.TrimSpace(tmp[1])
	}

	requestParse["request body"] = body

	return requestParse
}
