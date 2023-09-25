package main

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func TestGetPageContent(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Test HTTP body content"))
	}))
	defer server.Close()

	response, err := GetPageContent(server.URL)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		t.Errorf("Expected status code 200, got %v", response.Status)
	}

	body, err := io.ReadAll(response.Body)
	if err != nil {
		t.Fatalf("Error reading response body: %v", err)
	}
	expectedContent := "Test HTTP body content"
	if string(body) != expectedContent {
		t.Errorf("Expected response body %q, got %q", expectedContent, string(body))
	}
}

func TestGetImageHistory(t *testing.T) {
	tmpFile, err := os.CreateTemp("", "testfile-*.txt")
	if err != nil {
		t.Fatalf("Failed to create a temporary file: %v", err)
	}
	defer os.Remove(tmpFile.Name())
	defer tmpFile.Close()

	data := "1 2 3 4 5"

	_, err = tmpFile.WriteString(data)
	if err != nil {
		t.Fatalf("Error writing to the file: %v", err)
	}

	imageHistory, err := GetImageHistory(tmpFile.Name())
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	expectedImageHistory := map[int]bool{1: true, 2: true, 3: true, 4: true, 5: true}

	for key := range expectedImageHistory {
		if _, ok := imageHistory[key]; !ok {
			t.Errorf("Expected key %d in image history, but it was not found", key)
		}
	}

	for key := range imageHistory {
		if _, ok := expectedImageHistory[key]; !ok {
			t.Errorf("Unexpected key %d found in image history", key)
		}
	}
}
func TestGetImageUrl(t *testing.T) {
	testCases := []struct {
		input    string
		sub      string
		expected string
		errMsg   string
	}{
		{
			input:    `<img src="https://example.com/image.jpg" alt="Example Image">`,
			sub:      `src="`,
			expected: "https://example.com/image.jpg",
			errMsg:   "",
		},
		{
			input:    `<img alt="Example Image">`,
			sub:      `src="`,
			expected: "",
			errMsg:   "image URL not found",
		},
	}

	for _, testCase := range testCases {
		t.Run(fmt.Sprintf("sub=%s", testCase.sub), func(t *testing.T) {
			imageUrl, err := GetImageUrl(testCase.input, testCase.sub)

			if err != nil && err.Error() != testCase.errMsg {
				t.Errorf("Expected error message '%s', got '%s'", testCase.errMsg, err.Error())
			}

			if imageUrl != testCase.expected {
				t.Errorf("Expected URL '%s', got '%s'", testCase.expected, imageUrl)
			}
		})
	}
}

func TestCalculateMD5(t *testing.T) {
	data := bytes.NewBufferString("Hello, World!")

	md5Hash, err := CalculateMD5(*data)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	expectedHash := "65a8e27d8879283831b664bd8b7f0ad4"

	if md5Hash != expectedHash {
		t.Errorf("Expected MD5 hash %s, got %s", expectedHash, md5Hash)
	}
}

func TestWriteBufferToFile(t *testing.T) {
	content := "This is a test."

	tempFile, err := os.CreateTemp("", "testfile_*.txt")
	if err != nil {
		t.Errorf("Error creating temporary file: %v", err)
		return
	}
	defer tempFile.Close()
	defer os.Remove(tempFile.Name())

	buffer := bytes.NewBufferString(content)

	err = WriteBufferToFile(tempFile.Name(), "", buffer)
	if err != nil {
		t.Errorf("Error writing buffer to file: %v", err)
	}

	fileContent, err := os.ReadFile(tempFile.Name())
	if err != nil {
		t.Errorf("Error reading file: %v", err)
		return
	}

	expectedContent := []byte(content)
	if !bytes.Equal(fileContent, expectedContent) {
		t.Errorf("File content does not match expected content.")
	}
}
