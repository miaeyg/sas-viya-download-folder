package core

import (
	"crypto/tls"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"
)

// CallRest is the default function for calling SAS Viya REST APIs
func CallRest(baseURL string, endpoint string, headers map[string][]string, method string, data url.Values, query url.Values) []byte {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	restURL := baseURL + endpoint
	req, errreq := http.NewRequest(method, restURL, strings.NewReader(data.Encode()))
	if errreq != nil {
		log.Println(errreq)
	}
	req.Header = headers
	req.URL.RawQuery = query.Encode()
	client := &http.Client{Transport: tr}
	resp, err := client.Do(req)
	if err != nil {
		log.Println(err)
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
	}
	return []byte(body)
}
