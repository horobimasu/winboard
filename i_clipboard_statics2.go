package winboard

import (
	"fmt"
	"syscall"
	"unsafe"
)

var clipboardStaticsGuid = syscall.GUID{
	Data1: 0xD2AC1B6A,
	Data2: 0xD29F,
	Data3: 0x554B,
	Data4: [8]byte{0xB3, 0x03, 0xF0, 0x45, 0x23, 0x45, 0xFE, 0x02},
}

type iClipboardStatics2VTable struct {
	iInspectableVTable
	GetHistoryItemsAsync        uintptr
	ClearHistory                uintptr
	DeleteItemFromHistory       uintptr
	SetHistoryItemAsContent     uintptr
	IsHistoryEnabled            uintptr
	IsRoamingEnabled            uintptr
	SetContentWithOptions       uintptr
	AddHistoryChanged           uintptr
	RemoveHistoryChanged        uintptr
	AddRoamingEnabledChanged    uintptr
	RemoveRoamingEnabledChanged uintptr
	AddHistoryEnabledChanged    uintptr
	RemoveHistoryEnabledChanged uintptr
}

type iClipboardStatics2 struct {
	VTable *iClipboardStatics2VTable
}

func (statics *iClipboardStatics2) GetHistoryItemsAsync() (*iAsyncOperation, error) {
	var operation *iAsyncOperation

	ret, _, _ := syscall.SyscallN(
		statics.VTable.GetHistoryItemsAsync,
		uintptr(unsafe.Pointer(statics)),
		uintptr(unsafe.Pointer(&operation)),
	)

	if int32(ret) < 0 {
		return nil, fmt.Errorf("failed to get history: 0x%X", ret)
	}

	return operation, nil
}

func (statics *iClipboardStatics2) ClearHistory() (bool, error) {
	var cleared uint8

	ret, _, _ := syscall.SyscallN(
		statics.VTable.ClearHistory,
		uintptr(unsafe.Pointer(statics)),
		uintptr(unsafe.Pointer(&cleared)),
	)

	if int32(ret) < 0 {
		return false, fmt.Errorf("failed to clear history: 0x%X", ret)
	}

	return cleared != 0, nil
}

func getStatics() (*iClipboardStatics2, error) {
	class, err := createRtString("Windows.ApplicationModel.DataTransfer.Clipboard")
	if err != nil {
		return nil, err
	}

	defer deleteRtString(class)

	var statics *iClipboardStatics2

	ret, _, _ := syscall.SyscallN(
		roGetActivationFactory.Addr(),
		class,
		uintptr(unsafe.Pointer(&clipboardStaticsGuid)),
		uintptr(unsafe.Pointer(&statics)),
	)

	if int32(ret) < 0 {
		return nil, fmt.Errorf("failed to get activation factory: 0x%X", ret)
	}

	return statics, nil
}
