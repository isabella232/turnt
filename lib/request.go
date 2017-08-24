package lib

import (
	"net/http"
	"strings"
	"io/ioutil"
	"bytes"
	"encoding/json"
)

// Generate and send the request to the Turnstile server
func GenerateRequest(method string, url string, payload string, headers map[string]string) (error, bytes.Buffer) {
	client := &http.Client{}
	emptyBuffer := *bytes.NewBuffer(make([]byte, 0))

	req, err := http.NewRequest(method, url, strings.NewReader(payload))
	if err != nil {
		return err, emptyBuffer
	}
	for header, value := range headers {
		req.Header.Add(header, value)
	}

	resp, err := client.Do(req)
	defer resp.Body.Close()
	if err != nil {
		return err, emptyBuffer
	}

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err, emptyBuffer
	}

	var out bytes.Buffer
	if err := json.Indent(&out, data, "", "  "); err != nil {
		return err, emptyBuffer
	}

	return nil, out
}