package winboard_test

import (
	"testing"

	"github.com/horobimasu/winboard"
)

func TestClearHistory(t *testing.T) {
	cleared, err := winboard.ClearClipboardHistory()
	if err != nil {
		t.Fatalf("failed to clear - %v", err)
	}

	t.Logf("call success, cleared: %t", cleared)
}
