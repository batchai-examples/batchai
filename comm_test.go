package comm

import (
	"bytes"
	"strings"
	"testing"

	"github.com/pkg/errors"
)

func TestReadBytesP(t *testing.T) {
	// Test case: Happy path with valid input
	reader := strings.NewReader("Hello, World!")
	result := ReadBytesP(reader)
	expected := []byte("Hello, World!")
	if !bytes.Equal(result, expected) {
		t.Errorf("Expected %v, got %v", expected, result)
	}

	// Test case: Negative case with error from io.ReadAll
	reader = &errorReader{}
	result = ReadBytesP(reader)
	if result != nil {
		t.Errorf("Expected nil, got %v", result)
	}
}

func TestReadTextP(t *testing.T) {
	// Test case: Happy path with valid input
	reader := strings.NewReader("Hello, World!")
	result := ReadTextP(reader)
	expected := "Hello, World!"
	if result != expected {
		t.Errorf("Expected %s, got %s", expected, result)
	}

	// Test case: Negative case with error from io.ReadAll
	reader = &errorReader{}
	result = ReadTextP(reader)
	if result != "" {
		t.Errorf("Expected empty string, got %s", result)
	}
}

func TestReadLines(t *testing.T) {
	// Test case: Happy path with valid input
	reader := strings.NewReader("Line1\nLine2\nLine3")
	result := ReadLines(reader)
	expected := []string{"Line1", "Line2", "Line3"}
	if !equalSlices(result, expected) {
		t.Errorf("Expected %v, got %v", expected, result)
	}

	// Test case: Negative case with error from io.ReadAll
	reader = &errorReader{}
	result = ReadLines(reader)
	expected = []string{}
	if !equalSlices(result, expected) {
		t.Errorf("Expected %v, got %v", expected, result)
	}
}

func TestSplitBufferByLines(t *testing.T) {
	// Test case: Happy path with valid input
	buffer := new(string)
	*buffer = "Line1\nLine2"
	result := SplitBufferByLines(buffer)
	expected := []string{"Line1"}
	if !equalSlices(result, expected) {
		t.Errorf("Expected %v, got %v", expected, result)
	}
	// Check if the buffer is updated correctly
	if *buffer != "Line2" {
		t.Errorf("Expected buffer to be 'Line2', got '%s'", *buffer)
	}

	// Test case: Happy path with incomplete line at the end
	buffer = new(string)
	*buffer = "Line1\nLine2"
	result = SplitBufferByLines(buffer)
	expected = []string{"Line1"}
	if !equalSlices(result, expected) {
		t.Errorf("Expected %v, got %v", expected, result)
	}
	// Check if the buffer is updated correctly
	if *buffer != "Line2" {
		t.Errorf("Expected buffer to be 'Line2', got '%s'", *buffer)
	}

	// Test case: Negative case with empty input
	buffer = new(string)
	*buffer = ""
	result = SplitBufferByLines(buffer)
	expected = []string{}
	if !equalSlices(result, expected) {
		t.Errorf("Expected %v, got %v", expected, result)
	}
	// Check if the buffer is updated correctly
	if *buffer != "" {
		t.Errorf("Expected buffer to be empty, got '%s'", *buffer)
	}

	// Test case: Negative case with error from strings.Split
	buffer = new(string)
	*buffer = "Line1\nLine2"
	reader := &errorReader{}
	result = SplitBufferByLines(buffer)
	expected = []string{}
	if !equalSlices(result, expected) {
		t.Errorf("Expected %v, got %v", expected, result)
	}
	// Check if the buffer is updated correctly
	if *buffer != "Line1\nLine2" {
		t.Errorf("Expected buffer to be 'Line1\nLine2', got '%s'", *buffer)
	}
}

func equalSlices(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

type errorReader struct{}

func (r *errorReader) Read(p []byte) (n int, err error) {
	return 0, errors.New("read error")
}
