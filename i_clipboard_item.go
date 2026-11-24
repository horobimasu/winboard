package winboard

import (
	"fmt"
	"syscall"
	"unsafe"
)

type iClipboardHistoryItemVTable struct {
	iInspectableVTable
	GetId        uintptr
	GetTimestamp uintptr
	GetContent   uintptr
}

type iClipboardHistoryItem struct {
	VTable *iClipboardHistoryItemVTable
}

func (item *iClipboardHistoryItem) GetContent() (*iDataPackageView, error) {
	var data *iDataPackageView

	ret, _, _ := syscall.SyscallN(
		item.VTable.GetContent,
		uintptr(unsafe.Pointer(item)),
		uintptr(unsafe.Pointer(&data)),
	)

	if int32(ret) < 0 {
		return nil, fmt.Errorf("failed to get text content: 0x%X", ret)
	}

	return data, nil
}

func iClipboardHistoryItemGetText(item *iClipboardHistoryItem, format uintptr) (string, bool, error) {
	data, err := item.GetContent()
	if err != nil {
		return "", false, err
	}

	defer release(data)

	contains, err := data.Contains(format)
	if err != nil {
		return "", false, err
	}

	if !contains {
		return "", false, nil
	}

	operation, err := data.GetTextAsync()
	if err != nil {
		return "", false, err
	}

	defer release(operation)

	err = awaitAsync(operation)
	if err != nil {
		return "", false, err
	}

	var rtString uintptr

	err = operation.GetResults(unsafe.Pointer(&rtString))
	if err != nil {
		return "", false, err
	}

	text := rtStringToGoString(rtString)
	deleteRtString(rtString)

	return text, true, nil
}
