package winboard

import (
	"fmt"
	"syscall"
	"unsafe"
)

type iVectorViewVTable struct {
	iInspectableVTable
	GetAt   uintptr
	GetSize uintptr
	IndexOf uintptr
	GetMany uintptr
}

type iVectorView struct {
	VTable *iVectorViewVTable
}

func (view *iVectorView) GetSize() (uint32, error) {
	var size uint32

	ret, _, _ := syscall.SyscallN(
		view.VTable.GetSize,
		uintptr(unsafe.Pointer(view)),
		uintptr(unsafe.Pointer(&size)),
	)

	if int32(ret) < 0 {
		return 0, fmt.Errorf("failed to get size: 0x%X", ret)
	}

	return size, nil
}

func iVectorViewGetAt[T any](view *iVectorView, index uint32) (*T, error) {
	var result *T

	ret, _, _ := syscall.SyscallN(
		view.VTable.GetAt,
		uintptr(unsafe.Pointer(view)),
		uintptr(index),
		uintptr(unsafe.Pointer(&result)),
	)

	if int32(ret) < 0 {
		return nil, fmt.Errorf("failed to get item: 0x%X", ret)
	}

	return result, nil
}
