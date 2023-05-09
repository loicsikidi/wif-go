package attribute

import (
	"fmt"
	"testing"
)

func TestGetAttributeInputVar(t *testing.T) {
	tests := []struct {
		val      map[string]any
		expected string
	}{
		{
			val:      map[string]any{"google.subject": "123456789"},
			expected: `{"assertion":{},"attribute":{},"google":{"subject":"123456789"}}`,
		},
		{
			val:      map[string]any{"google.subject": "123456789", "attribute.sub": "123456789"},
			expected: `{"assertion":{},"attribute":{"sub":"123456789"},"google":{"subject":"123456789"}}`,
		},
		{
			val:      map[string]any{"google.subject": "123456789", "attribute.sub": "123456789", "assertion": map[string]any{"sub": "123456789"}},
			expected: `{"assertion":{"sub":"123456789"},"attribute":{"sub":"123456789"},"google":{"subject":"123456789"}}`,
		},
	}

	for i, tst := range tests {
		tc := tst
		t.Run(fmt.Sprintf("[%d]", i), func(t *testing.T) {
			if out, _ := GetAttributeInputVar(tc.val); out != tc.expected {
				t.Fatalf("GetAttributeInputVar(%s) = %s, expected %s", tc.val, out, tc.expected)
			}
		})
	}
}
