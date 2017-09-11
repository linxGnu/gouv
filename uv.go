package binduv

// #cgo pkg-config: libuv
/*
#include <uv.h>
#include <stdio.h>
*/
import "C"
import "errors"

// GetLibuvVersion get libuv version in hex
func GetLibuvVersion() int {
	return int(C.uv_version())
}

// GetLibuvVersionString get libuv version in string
func GetLibuvVersionString() string {
	return C.GoString(C.uv_version_string())
}

// ParseUvErr parsing uv error
func ParseUvErr(r C.int) error {
	return errors.New(C.GoString(C.uv_strerror(r)))
}

// UvErrName returns the error name for the given error code. Leaks a few bytes of memory when you call it with an unknown error code.
func UvErrName(r C.int) string {
  return C.GoString(C.uv_err_name(r))
}

// TranslateSysError (uv_translate_sys_error) Returns the libuv error code equivalent to the given platform dependent error code: POSIX error codes on Unix (the ones stored in errno), and Win32 error codes on Windows (those returned by GetLastError() or WSAGetLastError()).
// If sys_errno is already a libuv error, it is simply returned.
func TranslateSysError(r C.int) C.int {
  return C.uv_translate_sys_error(r)
}
