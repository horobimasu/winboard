package winboard_test

import (
	"testing"

	"github.com/horobimasu/winboard"
)

func TestGetPinnedItems(t *testing.T) {
	pinned, err := winboard.GetPinnedItems()
	if err != nil {
		t.Fatalf("failed to get pinned items - %v", err)
	}

	t.Logf("items (%d): %v", len(pinned), pinned)
}
