package goNest

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"time"
)

type SimpleResponse struct {
	StatusCode   int
	ResponseBody string
}

// CheckInternet attempts to connect to Google's public DNS server
// to verify if there is an active internet connection.
func CheckInternet() bool {
	conn, err := net.DialTimeout("tcp", "8.8.8.8:53", 2*time.Second)
	if err != nil {
		return false
	}
	defer conn.Close()
	return true
}

// SendGetRequest sends a GET request to the specified URL with optional headers and a 2-second timeout.
// It returns the response body as a string or an error message.
func SendGetRequest(url string, headers map[string]string) SimpleResponse {
	client := &http.Client{
		Timeout: 2 * time.Second, // Set a 2-second timeout
	}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return SimpleResponse{
			StatusCode: 500, ResponseBody: fmt.Sprintf("sendGetRequest : %v", err),
		}
	}

	// Set headers if provided
	for key, value := range headers {
		req.Header.Set(key, value)
	}

	resp, err := client.Do(req)
	if err != nil {
		return SimpleResponse{
			StatusCode: 500, ResponseBody: fmt.Sprintf("sendGetRequest : %v", err),
		}
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return SimpleResponse{
			StatusCode: 500, ResponseBody: fmt.Sprintf("sendGetRequest : %v", err),
		}
	}

	strBody := BytesToString(body, false)
	return SimpleResponse{
		StatusCode: resp.StatusCode, ResponseBody: strBody,
	}
}

// SendPostRequest sends a POST request to the specified endpoint with headers and a JSON request body.
// It returns the response body as a string or an error message if any issue occurs.
func SendPostRequest(endpoint string, headers map[string]string, reqBody map[string]interface{}) SimpleResponse {
	client := &http.Client{
		Timeout: 10 * time.Second, // 10-second timeout
	}

	// Marshal request body into JSON
	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return SimpleResponse{
			StatusCode: 500, ResponseBody: fmt.Sprintf("sendPostRequest : %v", err),
		}
	}

	req, err := http.NewRequest("POST", endpoint, bytes.NewBuffer(jsonData))
	if err != nil {
		return SimpleResponse{
			StatusCode: 500, ResponseBody: fmt.Sprintf("sendPostRequest : %v", err),
		}
	}

	// Set JSON headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	// Add custom headers if provided
	for key, value := range headers {
		req.Header.Set(key, value)
	}

	resp, err := client.Do(req)
	if err != nil {
		return SimpleResponse{
			StatusCode: 500, ResponseBody: fmt.Sprintf("sendPostRequest : %v", err),
		}
	}
	defer resp.Body.Close()

	// Read the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return SimpleResponse{
			StatusCode: 500, ResponseBody: fmt.Sprintf("sendPostRequest : %v", err),
		}
	}

	return SimpleResponse{
		StatusCode: resp.StatusCode, ResponseBody: BytesToString(body, false),
	}
}
