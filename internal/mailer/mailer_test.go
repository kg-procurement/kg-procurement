
package mailer

import (
	"testing"
)

func TestEmailStatusEnum_String(t *testing.T) {
	tests := []struct {
		status   EmailStatusEnum
		expected string
	}{
		{Success, "success"},
		{Failed, "failed"},
		{InProgress, "in_progress"},
		{Completed, "completed"},
		{EmailStatusEnum(999), "unknown"},
	}

	for _, test := range tests {
		result := test.status.String()
		if result != test.expected {
			t.Errorf("Expected %s, but got %s", test.expected, result)
		}
	}
}

func TestParseEmailStatusEnum(t *testing.T) {
	tests := []struct {
		status   string
		expected EmailStatusEnum
		err      bool
	}{
		{"success", Success, false},
		{"failed", Failed, false},
		{"in_progress", InProgress, false},
		{"completed", Completed, false},
		{"invalid", -1, true},
	}

	for _, test := range tests {
		result, err := ParseEmailStatusEnum(test.status)
		if (err != nil) != test.err {
			t.Errorf("Expected error: %v, but got error: %v", test.err, err)
		}
		if result != test.expected {
			t.Errorf("Expected %d, but got %d", test.expected, result)
		}
	}
}