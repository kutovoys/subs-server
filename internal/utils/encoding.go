package utils

import (
	"bufio"
	"bytes"
	"encoding/base64"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

var (
	SupportedProtocols = []string{"vless://", "vmess://", "ss://", "trojan://", "socks://"}

	HTTPProtocols = []string{"http://", "https://"}
)

func ProcessContent(content []byte) ([]byte, error) {
	scanner := bufio.NewScanner(bytes.NewReader(content))
	var resultLines []string

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}

		if hasProtocol(line, SupportedProtocols) {
			resultLines = append(resultLines, line)
			continue
		}

		if hasProtocol(line, HTTPProtocols) {
			fetchedLines, err := fetchFromURL(line)
			if err != nil {
				fmt.Printf("Warning: failed to fetch from %s: %v\n", line, err)
				continue
			}
			resultLines = append(resultLines, fetchedLines...)
			continue
		}

		fmt.Printf("Warning: skipping line with unsupported protocol: %s\n", line)
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error scanning content: %w", err)
	}

	resultContent := strings.Join(resultLines, "\n")
	return []byte(base64.StdEncoding.EncodeToString([]byte(resultContent))), nil
}

func hasProtocol(line string, protocols []string) bool {
	for _, protocol := range protocols {
		if strings.HasPrefix(line, protocol) {
			return true
		}
	}
	return false
}

func fetchFromURL(url string) ([]string, error) {
	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	resp, err := client.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch URL: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("bad status code: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	decoded, err := base64.StdEncoding.DecodeString(string(body))
	var content string
	if err != nil {
		content = string(body)
	} else {
		content = string(decoded)
	}

	var lines []string
	scanner := bufio.NewScanner(strings.NewReader(content))
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line != "" {
			lines = append(lines, line)
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error scanning fetched content: %w", err)
	}

	return lines, nil
}

func EncodeToBase64WithoutEmptyLines(content []byte) string {
	processed, err := ProcessContent(content)
	if err != nil {
		lines := strings.Split(string(content), "\n")
		var nonEmptyLines []string
		for _, line := range lines {
			if strings.TrimSpace(line) != "" {
				nonEmptyLines = append(nonEmptyLines, line)
			}
		}
		filteredContent := strings.Join(nonEmptyLines, "\n")
		return base64.StdEncoding.EncodeToString([]byte(filteredContent))
	}
	return string(processed)
}
