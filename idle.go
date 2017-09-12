package gouv

// #cgo pkg-config: libuv
/*
#include <uv.h>
#include <stdlib.h>
uv_idle_t* mallocIdleT() {
	return (uv_idle_t*)malloc(sizeof(uv_idle_t));
}
*/
import "C"
import "unsafe"

// UvIdle idle handles will run the given callback once per loop iteration, right before the uv_prepare_t handles.
type UvIdle struct {
	i *C.uv_idle_t
	l *C.uv_loop_t
	Handle
}

// UvIdleInit initialize the prepare handle
func UvIdleInit(loop *UvLoop, data interface{}) (*UvIdle, error) {
	t := C.mallocIdleT()

	if loop == nil {
		loop = UvLoopDefault()
	}

	if r := C.uv_idle_init(loop.GetNativeLoop(), t); r != 0 {
		return nil, ParseUvErr(r)
	}

	t.data = unsafe.Pointer(&callbackInfo{data: data})
	return &UvIdle{t, loop.GetNativeLoop(), Handle{(*C.uv_handle_t)(unsafe.Pointer(t)), t.data}}, nil
}

// Start (uv_prepare_start) start the timer. timeout and repeat are in milliseconds.
func (t *UvIdle) Start(cb func(*Handle, int)) (err error) {
	cbi := (*callbackInfo)(t.i.data)
	cbi.idle_cb = cb

	if r := uv_idle_start(t.i); r != 0 {
		err = ParseUvErr(r)
		return
	}

	return
}

// Stop (uv_idle_stop) the timer, the callback will not be called anymore.
func (t *UvIdle) Stop() (err error) {
	if r := C.uv_idle_stop(t.i); r != 0 {
		err = ParseUvErr(r)
		return
	}

	return
}

// Freemem freemem of prepare
func (t *UvIdle) Freemem() {
	C.free(unsafe.Pointer(t.i))
}
