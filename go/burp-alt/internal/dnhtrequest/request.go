package dnhtrequest

import (
	"bufio"
	"net/http"
	"os"
	"strings"
)

func ParseRequestFile(filePath string) (method, url string, headers http.Header, body string, err error) {
	file, err := os.Open(filePath)
	if err != nil {
		return "", "", nil, "", err
	}
	defer file.Close()

	reader := bufio.NewReader(file)

	// Read the request line
	requestLine, _ := reader.ReadString('\n')
	method, url, _ = parseRequestLine(requestLine)

	// Read headers
	headers = make(http.Header)
	for {
		line, _ := reader.ReadString('\n')
		line = strings.TrimSpace(line)
		if line == "" {
			break
		}
		key, value := parseHeader(line)
		headers.Add(key, value)
	}

	// Read body
	bodyBuilder := new(strings.Builder)
	_, _ = reader.WriteTo(bodyBuilder)
	body = bodyBuilder.String()

	return method, url, headers, body, nil
}

func parseRequestLine(line string) (method, url, version string) {
	parts := strings.Split(line, " ")
	if len(parts) >= 3 {
		method = parts[0]
		url = parts[1]
		version = parts[2]
	}
	return
}

func parseHeader(line string) (key, value string) {
	parts := strings.SplitN(line, ":", 2)
	if len(parts) == 2 {
		key = strings.TrimSpace(parts[0])
		value = strings.TrimSpace(parts[1])
	}
	return
}
