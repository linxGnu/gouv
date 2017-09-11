package binduv

// #cgo pkg-config: libuv
/*
#include <uv.h>
#include <stdlib.h>
uv_check_t* mallocCheckT() {
	return (uv_check_t*)malloc(sizeof(uv_check_t));
}
*/
import "C"
import "unsafe"

// UvCheck check handles will run the given callback once per loop iteration, right after polling for i/o.
type UvCheck struct {
	c *C.uv_check_t
	l *C.uv_loop_t
	Handle
}

// UvCheckInit initialize the prepare handle
func UvCheckInit(loop *UvLoop, data interface{}) (timer *UvCheck, err error) {
	t := C.mallocCheckT()

	if loop == nil {
		loop = UvLoopDefault()
	}

	if r := C.uv_check_init(loop.GetNativeLoop(), t); r != 0 {
		return nil, ParseUvErr(r)
	}

	t.data = unsafe.Pointer(&callback_info{data: data})
	return &UvCheck{t, loop.GetNativeLoop(), Handle{(*C.uv_handle_t)(unsafe.Pointer(t)), t.data}}, nil
}

// Start (uv_prepare_start) start the timer. timeout and repeat are in milliseconds.
func (t *UvCheck) Start(cb func(*Handle)) (err error) {
	cbi := (*callback_info)(t.c.data)
	cbi.check_cb = cb

	if r := uv_check_start(t.c); r != 0 {
		err = ParseUvErr(r)
		return
	}

	return nil
}

// Stop (uv_prepare_stop) the timer, the callback will not be called anymore.
func (t *UvCheck) Stop() (err error) {
	if r := C.uv_check_stop(t.c); r != 0 {
		err = ParseUvErr(r)
		return
	}

	return nil
}

// Freemem freemem of prepare
func (t *UvCheck) Freemem() {
	C.free(unsafe.Pointer(t.c))
}
