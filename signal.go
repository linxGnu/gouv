package gouv

// #cgo pkg-config: libuv
/*
#include <uv.h>
#include <stdlib.h>
uv_signal_t* mallocSignalT() {
	return (uv_signal_t*)malloc(sizeof(uv_signal_t));
}
*/
import "C"
import "unsafe"

// UvSignal handles implement Unix style signal handling on a per-event loop bases.
type UvSignal struct {
	s *C.uv_signal_t
	l *C.uv_loop_t
	Handle
}

// UvSignalInit initialize the handle.
func UvSignalInit(loop *UvLoop, data interface{}) (*UvSignal, error) {
	t := C.mallocSignalT()

	if loop == nil {
		loop = UvLoopDefault()
	}

	res := &UvSignal{}
	t.data = unsafe.Pointer(&callbackInfo{data: data, ptr: res})
	res.s, res.l, res.Handle = t, loop.GetNativeLoop(), Handle{(*C.uv_handle_t)(unsafe.Pointer(t)), t.data, res}
	if r := C.uv_signal_init(loop.GetNativeLoop(), t); r != 0 {
		return nil, ParseUvErr(r)
	}

	return res, nil
}

// Start (uv_signal_start) start the handle with the given callback, watching for the given signal.
func (s *UvSignal) Start(cb func(*Handle, int), sigNum int) C.int {
	cbi := (*callbackInfo)(s.s.data)
	cbi.signal_cb = cb

	return uv_signal_start(s.s, sigNum)
}

// StartOneShot (uv_signal_start_oneshot) same functionality as uv_signal_start() but the signal handler is reset the moment the signal is received.
func (s *UvSignal) StartOneShot(cb func(*Handle, int), sigNum int) C.int {
	cbi := (*callbackInfo)(s.s.data)
	cbi.signal_cb = cb

	return uv_signal_start_oneshot(s.s, sigNum)
}

// Stop (uv_signal_stop) stop the handle, the callback will no longer be called.
func (s *UvSignal) Stop() C.int {
	return C.uv_signal_stop(s.s)
}

// GetSignalHandle get handle
func (s *UvSignal) GetSignalHandle() *C.uv_signal_t {
	return s.s
}
