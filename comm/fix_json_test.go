package comm

import (
	"strings"
	"testing"
)

func TestFixJson(t *testing.T) {
	testCases := []struct {
		input      string
		isGolang   bool
		expected   string
	}{
		{
			input:      `{"key":"value"}`,
			isGolang:   false,
			expected:   `{"key":"value"}`,
		},
		{
			input:      `{"key":"value\n"}`,
			isGolang:   false,
			expected:   `{"key":"value\n"}`,
		},
		{
			input:      `{"key":"value\t"}`,
			isGolang:   false,
			expected:   `{"key":"value\t"}`,
		},
		{
			input:      `{"key":"value\r"}`,
			isGolang:   false,
			expected:   `{"key":"value\r"}`,
		},
		{
			input:      `{"key":"value\"}`,
			isGolang:   false,
			expected:   `{"key":"value\"}`,
		},
		{
			input:      `import ("fmt")`,
			isGolang:   true,
			expected:   `import ("fmt")`,
		},
		{
			input:      `import ("fmt" "os")`,
			isGolang:   true,
			expected:   `import ("fmt" "os")`,
		},
		{
			input:      `{}`,
			isGolang:   false,
			expected:   `{}`,
		},
		{
			input:      `[1,2,3]`,
			isGolang:   false,
			expected:   `[1,2,3]`,
		},
		{
			input:      `{"key":"value\\n"}`,
			isGolang:   false,
			expected:   `{"key":"value\n"}`,
		},
	}

	for _, tc := range testCases {
		t.Run(fmt.Sprintf("Input=%s, IsGolang=%v", tc.input, tc.isGolang), func(t *testing.T) {
			result := FixJson(tc.input, tc.isGolang)
			if result != tc.expected {
				t.Errorf("FixJson(%q, %v) = %q; want %q", tc.input, tc.isGolang, result, tc.expected)
			}
		})
	}
}
