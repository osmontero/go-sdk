package go_sdk

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
	"crypto/tls"
)

// DoReq sends an HTTP request and processes the response.
//
// This function sends an HTTP request to the specified URL with the given
// method, data, and headers. It returns the response body unmarshalled into
// the specified response type, the HTTP status code, and an error if any
// occurred during the process.
//
// Type Parameters:
//   - response: The type into which the response body will be unmarshalled.
//
// Parameters:
//   - url: The URL to which the request is sent.
//   - data: The request payload as a byte slice.
//   - method: The HTTP method to use for the request (e.g., "GET", "POST").
//   - headers: A map of headers to include in the request.
//
// Returns:
//   - response: The response body unmarshalled into the specified type.
//   - int: The HTTP status code of the response.
//   - error: An error if any occurred during the request or response
//     processing, otherwise nil.
func DoReq[response any](url string,
	data []byte, method string,
	headers map[string]string) (response, int, error) {

	var result response

	// Add request size limit
	const maxRequestSize = 1048576 // 1MB
	if len(data) > maxRequestSize {
		return result, http.StatusRequestEntityTooLarge, fmt.Errorf("request too large")
	}

	req, err := http.NewRequest(method, url, bytes.NewBuffer(data))
	if err != nil {
		return result,
			http.StatusInternalServerError,
			fmt.Errorf("error creating request: %s", err.Error())
	}

	for k, v := range headers {
		req.Header.Add(k, v)
	}

	// Configure HTTP client with security settings
	client := &http.Client{
		Timeout: 30 * time.Second,
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				MinVersion: tls.VersionTLS12,
			},
			DisableCompression: true,
		},
	}

	resp, err := client.Do(req)
	if err != nil {
		return result,
			http.StatusInternalServerError,
			fmt.Errorf("error sending request: %s", err.Error())
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return result,
			http.StatusInternalServerError,
			fmt.Errorf("error reading response: %s", err.Error())
	}

	if resp.StatusCode >= 400 {
		return result,
			resp.StatusCode,
			fmt.Errorf("received status code %d", resp.StatusCode)
	}

	if resp.StatusCode == http.StatusNoContent {
		return result, resp.StatusCode, nil
	}

	err = json.Unmarshal(body, &result)
	if err != nil {
		return result,
			http.StatusInternalServerError,
			fmt.Errorf("error unmarshalling response: %s", err.Error())
	}

	return result, resp.StatusCode, nil
}

// Download downloads the content from the specified URL and saves it to the specified file.
// It returns an error if any error occurs during the process.
//
// Parameters:
//   - url: The URL from which to download the content.
//   - file: The path to the file where the content should be saved.
//
// Returns:
//   - error: An error object if an error occurs, otherwise nil.
func Download(url, file string) error {
	out, err := os.Create(file)
	if err != nil {
		return fmt.Errorf("could not create file: %s", err.Error())
	}

	defer out.Close()

	// Add secure HTTP client configuration
	client := &http.Client{
		Timeout: 30 * time.Second,
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				MinVersion: tls.VersionTLS12,
			},
		},
	}

	resp, err := client.Get(url)
	if err != nil {
		return fmt.Errorf("could not do request to the URL: %s", err.Error())
	}

	defer resp.Body.Close()

	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return fmt.Errorf("could not save data to file: %s", err.Error())
	}

	return nil
}
