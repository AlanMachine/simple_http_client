package main

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
)

func main() {
	data, err := httpClient("https://google.com/")
	if err != nil {
		log.Fatal(err)
	}

	re, err := regexp.Compile(`(?m)<title>(.+?)</title>`)
	if err != nil {
		log.Fatal(err)
	}

	match := re.FindStringSubmatch(data)
	if len(match) != 0 {
		fmt.Println("Title: ", match[1])
	}
}

func httpClient(url string) (string, error) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", nil
	}

	req.Header.Add("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b2")
	req.Header.Add("Accept-Encoding", "gzip")
	req.Header.Add("Accept-Language", "ru-RU,ru;q=0.9,en-US;q=0.8,en;q=0.7")
	req.Header.Add("User-Agent", "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/71.0.3578.98 Safari/537.36")

	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}

	body, err := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil {
		return "", err
	}
	if !resp.Uncompressed {
		gz, err := gzip.NewReader(bytes.NewReader(body))
		if err != nil {
			return "", nil
		}
		body, err = ioutil.ReadAll(gz)
		gz.Close()
		if err != nil {
			return "", nil
		}
	}

	return string(body), nil
}
