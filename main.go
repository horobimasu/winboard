package winboard

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"unsafe"
)

var pinnedItemsPath = filepath.Join(os.Getenv("LOCALAPPDATA"), "Microsoft", "Windows", "Clipboard", "Pinned")

type iInspectableVTable struct {
	iUnknownVTable
	GetIids             uintptr
	GetRuntimeClassName uintptr
	GetTrustLevel       uintptr
}

// this will not clear pinned items
func ClearClipboardHistory() (bool, error) {
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()

	err := initWinRt()
	if err != nil {
		return false, err
	}

	defer unInitWinRt()

	statics, err := getStatics()
	if err != nil {
		return false, err
	}

	defer release(statics)

	cleared, err := statics.ClearHistory()
	if err != nil {
		return false, err
	}

	return cleared, nil
}

func GetClipboardHistory() ([]string, error) {
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()

	err := initWinRt()
	if err != nil {
		return nil, err
	}

	defer unInitWinRt()

	statics, err := getStatics()
	if err != nil {
		return nil, err
	}

	defer release(statics)

	operation, err := statics.GetHistoryItemsAsync()
	if err != nil {
		return nil, err
	}

	defer release(operation)

	err = awaitAsync(operation)
	if err != nil {
		return nil, err
	}

	var res *iClipboardHistoryItemsResult

	err = operation.GetResults(unsafe.Pointer(&res))
	if err != nil {
		return nil, err
	}

	defer release(res)

	status, err := res.GetStatus()
	if err != nil {
		return nil, err
	}

	if status != 0 {
		return nil, fmt.Errorf("failed with status 0x%X", status)
	}

	items, err := res.GetItems()
	if err != nil {
		return nil, err
	}

	defer release(items)

	size, err := items.GetSize()
	if err != nil {
		return nil, err
	}

	format, err := createRtString("Text")
	if err != nil {
		return nil, err
	}

	defer deleteRtString(format)

	history := make([]string, 0, size)

	for i := range size {
		item, err := iVectorViewGetAt[iClipboardHistoryItem](items, i)
		if err != nil {
			return nil, err
		}

		str, found, err := iClipboardHistoryItemGetText(item, format)
		if err != nil {
			release(item)
			return nil, err
		}

		release(item)

		if found {
			history = append(history, str)
		}
	}

	return history, nil
}

func GetPinnedItems() ([]string, error) {
	files, err := filepath.Glob(filepath.Join(pinnedItemsPath, "*", "*", "VGV4dA==")) // not sure why they encode this
	if err != nil {
		return nil, err
	}

	var pinned []string

	for _, file := range files {
		encrypted, err := os.ReadFile(file)
		if err != nil {
			return nil, err
		}

		decrypted, err := unprotectSecret(encrypted)
		if err != nil {
			return nil, err
		}

		pinned = append(pinned, string(decrypted))
	}

	return pinned, nil
}
