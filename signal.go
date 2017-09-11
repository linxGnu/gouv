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

	if r := C.uv_signal_init(loop.GetNativeLoop(), t); r != 0 {
		return nil, ParseUvErr(r)
	}

	t.data = unsafe.Pointer(&callback_info{data: data})
	return &UvSignal{t, loop.GetNativeLoop(), Handle{(*C.uv_handle_t)(unsafe.Pointer(t)), t.data}}, nil
}

// Start (uv_signal_start) start the handle with the given callback, watching for the given signal.
func (s *UvSignal) Start(cb func(*Handle, C.int), sigNum int) (err error) {
	cbi := (*callback_info)(s.s.data)
	cbi.signal_cb = cb

	if r := uv_signal_start(s.s, sigNum); r != 0 {
		err = ParseUvErr(r)
		return
	}

	return
}

// StartOneShot (uv_signal_start_oneshot) same functionality as uv_signal_start() but the signal handler is reset the moment the signal is received.
func (s *UvSignal) StartOneShot(cb func(*Handle, C.int), sigNum int) (err error) {
	cbi := (*callback_info)(s.s.data)
	cbi.signal_cb = cb

	if r := uv_signal_start_oneshot(s.s, sigNum); r != 0 {
		err = ParseUvErr(r)
		return
	}

	return
}

// Stop (uv_signal_stop) stop the handle, the callback will no longer be called.
func (s *UvSignal) Stop() error {
	if r := C.uv_signal_stop(s.s); r != 0 {
		return ParseUvErr(r)
	}

	return nil
}

// Freemem freemem of poll handle
func (s *UvSignal) Freemem() {
	C.free(unsafe.Pointer(s.s))
}
