package winboard

import (
	"fmt"
	"syscall"
	"unsafe"
)

type iClipboardHistoryItemsResultVTable struct {
	iInspectableVTable
	GetStatus uintptr
	GetItems  uintptr
}

type iClipboardHistoryItemsResult struct {
	VTable *iClipboardHistoryItemsResultVTable
}

func (items *iClipboardHistoryItemsResult) GetStatus() (int32, error) {
	var status int32

	ret, _, _ := syscall.SyscallN(
		items.VTable.GetStatus,
		uintptr(unsafe.Pointer(items)),
		uintptr(unsafe.Pointer(&status)),
	)

	if int32(ret) < 0 {
		return 0, fmt.Errorf("failed to get status: 0x%X", ret)
	}

	return status, nil
}

func (items *iClipboardHistoryItemsResult) GetItems() (*iVectorView, error) {
	var list *iVectorView

	ret, _, _ := syscall.SyscallN(
		items.VTable.GetItems,
		uintptr(unsafe.Pointer(items)),
		uintptr(unsafe.Pointer(&list)),
	)

	if int32(ret) < 0 {
		return nil, fmt.Errorf("failed to get items: 0x%X", ret)
	}

	return list, nil
}
