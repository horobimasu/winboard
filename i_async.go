package winboard

import (
	"errors"
	"fmt"
	"syscall"
	"time"
	"unsafe"
)

const _AWAIT_ASYNC_TIMEOUT = 10 * time.Second

const (
	_ASYNC_STATUS_STARTED   = 0 // AsyncStatus.Started
	_ASYNC_STATUS_COMPLETED = 1 // AsyncStatus.Completed
	_ASYNC_STATUS_CANCELED  = 2 // AsyncStatus.Canceled
	_ASYNC_STATUS_ERROR     = 3 // AsyncStatus.Error
)

var asyncInfoGuid = syscall.GUID{
	Data1: 0x00000036,
	Data2: 0x0000,
	Data3: 0x0000,
	Data4: [8]byte{0xC0, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x46},
}

type iAsyncInfoVTable struct {
	iInspectableVTable
	GetId        uintptr
	GetStatus    uintptr
	GetErrorCode uintptr
	Cancel       uintptr
	Close        uintptr
}

type iAsyncInfo struct {
	VTable *iAsyncInfoVTable
}

func (info *iAsyncInfo) GetStatus() (int32, error) {
	var status int32

	ret, _, _ := syscall.SyscallN(
		info.VTable.GetStatus,
		uintptr(unsafe.Pointer(info)),
		uintptr(unsafe.Pointer(&status)),
	)

	if int32(ret) < 0 {
		return 0, fmt.Errorf("failed to get status: 0x%X", ret)
	}

	return status, nil
}

func (info *iAsyncInfo) GetErrorCode() (int32, error) {
	var code int32

	ret, _, _ := syscall.SyscallN(
		info.VTable.GetErrorCode,
		uintptr(unsafe.Pointer(info)),
		uintptr(unsafe.Pointer(&code)),
	)

	if int32(ret) < 0 {
		return 0, fmt.Errorf("failed to get error code: 0x%X", ret)
	}

	return code, nil
}

type iAsyncOperationVTable struct {
	iInspectableVTable
	PutCompleted uintptr
	GetCompleted uintptr
	GetResults   uintptr
}

type iAsyncOperation struct {
	VTable *iAsyncOperationVTable
}

func (operation *iAsyncOperation) GetResults(output unsafe.Pointer) error {
	ret, _, _ := syscall.SyscallN(
		operation.VTable.GetResults,
		uintptr(unsafe.Pointer(operation)),
		uintptr(output),
	)

	if int32(ret) < 0 {
		return fmt.Errorf("failed to get results: 0x%X", ret)
	}

	return nil
}

func awaitAsync(operation *iAsyncOperation) error {
	info, err := queryInterface[iAsyncInfo](operation, &asyncInfoGuid)
	if err != nil {
		return err
	}

	defer release(info)

	timeout := time.Now().Add(_AWAIT_ASYNC_TIMEOUT)

	for {
		status, err := info.GetStatus()
		if err != nil {
			return err
		}

		switch status {
		case _ASYNC_STATUS_STARTED:
			if time.Now().After(timeout) {
				return errors.New("the operation timed out")
			}

			time.Sleep(2 * time.Millisecond)

		case _ASYNC_STATUS_COMPLETED:
			return nil

		case _ASYNC_STATUS_CANCELED:
			return errors.New("the operation has been canceled")

		case _ASYNC_STATUS_ERROR:
			code, err := info.GetErrorCode()
			if err != nil {
				return errors.New("the operation has failed")
			}

			return fmt.Errorf("the operation has failed: 0x%X", code)
		}
	}
}

func queryInterface[Out any, In any](obj *In, guid *syscall.GUID) (*Out, error) {
	var out *Out

	err := (*iUnknown)(unsafe.Pointer(obj)).QueryInterface(guid, unsafe.Pointer(&out))
	return out, err
}
