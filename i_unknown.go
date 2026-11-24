package winboard

import (
	"fmt"
	"syscall"
	"unsafe"
)

type iUnknownVTable struct {
	QueryInterface uintptr
	AddRef         uintptr
	Release        uintptr
}

type iUnknown struct {
	VTable *iUnknownVTable
}

func (unknown *iUnknown) QueryInterface(guid *syscall.GUID, output unsafe.Pointer) error {
	ret, _, _ := syscall.SyscallN(
		unknown.VTable.QueryInterface,
		uintptr(unsafe.Pointer(unknown)),
		uintptr(unsafe.Pointer(guid)),
		uintptr(output),
	)

	if int32(ret) < 0 {
		return fmt.Errorf("failed to query: 0x%X", ret)
	}

	return nil
}

func (unknown *iUnknown) Release() uint32 {
	ret, _, _ := syscall.SyscallN(unknown.VTable.Release, uintptr(unsafe.Pointer(unknown)))
	return uint32(ret)
}
