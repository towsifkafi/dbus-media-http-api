package main

import (
	"fmt"
	"net/url"
	"os"

	"encoding/base64"
)

func fetchBase64EncodedArt(artUrl string) (string, error) {
	if artUrl == "" {
		return "", fmt.Errorf("art URL is empty")
	}

	parsedUrl, err := url.Parse(artUrl)
	if err != nil {
		return "", err
	}

	if parsedUrl.Scheme != "file" {
		return "", fmt.Errorf("unsupported URL scheme")
	}

	filePath := parsedUrl.Path
	fileData, err := os.ReadFile(filePath)
	if err != nil {
		return "", err
	}

	return base64.StdEncoding.EncodeToString(fileData), nil
}
