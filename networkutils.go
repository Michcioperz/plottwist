package main

import (
	"net/http"
	"io/ioutil"
)

const TwistRoot = "https://twist.moe"

func FetchPageContents(url string) (string, error) {
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	return string(body), nil
}

func UrlPseudoJoin(path string) string {
	return TwistRoot + path
}
