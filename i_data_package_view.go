package winboard

import (
	"fmt"
	"syscall"
	"unsafe"
)

type iDataPackageViewVTable struct {
	iInspectableVTable
	GetProperties            uintptr
	GetRequestedOperation    uintptr
	ReportOperationCompleted uintptr
	GetAvailableFormats      uintptr
	Contains                 uintptr
	GetDataAsync             uintptr
	GetTextAsync             uintptr
	GetCustomTextAsync       uintptr
	GetUriAsync              uintptr
	GetHtmlFormatAsync       uintptr
	GetResourceMapAsync      uintptr
	GetRtfAsync              uintptr
	GetBitmapAsync           uintptr
	GetStorageItemsAsync     uintptr
}

type iDataPackageView struct {
	VTable *iDataPackageViewVTable
}

func (data *iDataPackageView) Contains(format uintptr) (bool, error) {
	var contains uint8

	ret, _, _ := syscall.SyscallN(
		data.VTable.Contains,
		uintptr(unsafe.Pointer(data)),
		format,
		uintptr(unsafe.Pointer(&contains)),
	)

	if int32(ret) < 0 {
		return false, fmt.Errorf("failed to check contais: 0x%X", ret)
	}

	return contains != 0, nil
}

func (data *iDataPackageView) GetTextAsync() (*iAsyncOperation, error) {
	var operation *iAsyncOperation

	ret, _, _ := syscall.SyscallN(
		data.VTable.GetTextAsync,
		uintptr(unsafe.Pointer(data)),
		uintptr(unsafe.Pointer(&operation)),
	)

	if int32(ret) < 0 {
		return nil, fmt.Errorf("failed to get text: 0x%X", ret)
	}

	return operation, nil
}
