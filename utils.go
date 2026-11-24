package winboard

import (
	"errors"
	"fmt"
	"syscall"
	"unicode/utf16"
	"unsafe"
)

const _NCRYPT_SILENT_FLAG = 0x00000040

const _RO_INIT_SINGLETHREADED = 0 // struct RO_INIT_TYPE

var (
	combase = syscall.NewLazyDLL("combase.dll")
	ncrypt  = syscall.NewLazyDLL("ncrypt.dll")

	roInitialize              = combase.NewProc("RoInitialize")
	roUninitialize            = combase.NewProc("RoUninitialize")
	roGetActivationFactory    = combase.NewProc("RoGetActivationFactory")
	windowsCreateString       = combase.NewProc("WindowsCreateString")
	windowsDeleteString       = combase.NewProc("WindowsDeleteString")
	windowsGetStringRawBuffer = combase.NewProc("WindowsGetStringRawBuffer")

	nCryptUnprotectSecret = ncrypt.NewProc("NCryptUnprotectSecret")
)

func initWinRt() error {
	ret, _, _ := syscall.SyscallN(roInitialize.Addr(), _RO_INIT_SINGLETHREADED)

	// we do not care about the higher bits
	lRet := uint32(ret)
	if int32(lRet) < 0 {
		return fmt.Errorf("failed to initialize the runtime: 0x%X", lRet)
	}

	return nil
}

func unInitWinRt() {
	syscall.SyscallN(roUninitialize.Addr())
}

func release[T any](ptr *T) {
	if ptr == nil {
		return
	}

	(*iUnknown)(unsafe.Pointer(ptr)).Release()
}

func createRtString(str string) (uintptr, error) {
	wide, err := syscall.UTF16FromString(str)
	if err != nil {
		return 0, err
	}

	var hString uintptr

	ret, _, _ := syscall.SyscallN(
		windowsCreateString.Addr(),
		uintptr(unsafe.Pointer(&wide[0])),
		uintptr(len(wide)-1),
		uintptr(unsafe.Pointer(&hString)),
	)

	if int32(ret) < 0 {
		return 0, fmt.Errorf("failed to create runtime string: 0x%X", ret)
	}

	return hString, nil
}

func deleteRtString(ptr uintptr) {
	if ptr == 0 {
		return
	}

	syscall.SyscallN(windowsDeleteString.Addr(), ptr)
}

func rtStringToGoString(ptr uintptr) string {
	if ptr == 0 {
		return ""
	}

	var length uint32

	ret, _, _ := syscall.SyscallN(windowsGetStringRawBuffer.Addr(), ptr, uintptr(unsafe.Pointer(&length)))
	if ret == 0 || length == 0 {
		return ""
	}

	decoded := utf16.Decode(unsafe.Slice((*uint16)(unsafe.Pointer(ret)), length))
	return string(decoded)
}

func unprotectSecret(data []byte) ([]byte, error) {
	var buffer *byte
	var bufferLen uint32

	ret, _, _ := syscall.SyscallN(
		nCryptUnprotectSecret.Addr(),
		0,
		_NCRYPT_SILENT_FLAG,
		uintptr(unsafe.Pointer(&data[0])),
		uintptr(len(data)),
		0,
		0,
		uintptr(unsafe.Pointer(&buffer)),
		uintptr(unsafe.Pointer(&bufferLen)),
	)

	if uint32(ret) != 0 {
		return nil, fmt.Errorf("failed to un-protect secret: 0x%X", ret)
	}

	if buffer == nil {
		return nil, errors.New("failed to un-protect secret: buffer empty")
	}

	defer syscall.LocalFree(syscall.Handle(uintptr(unsafe.Pointer(buffer))))

	unprotected := make([]byte, bufferLen)
	copy(unprotected, unsafe.Slice(buffer, bufferLen))

	return unprotected, nil
}
