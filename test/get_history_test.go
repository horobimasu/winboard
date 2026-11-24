package winboard_test

import (
	"testing"

	"github.com/horobimasu/winboard"
)

func TestGetHistory(t *testing.T) {
	history, err := winboard.GetClipboardHistory()
	if err != nil {
		t.Fatalf("failed to get history - %v", err)
	}

	t.Logf("history (%d): %v", len(history), history)
}
