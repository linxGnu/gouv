package gouv

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
func UvCheckInit(loop *UvLoop, data interface{}) (*UvCheck, error) {
	t := C.mallocCheckT()

	if loop == nil {
		loop = UvLoopDefault()
	}

	res := &UvCheck{}
	t.data = unsafe.Pointer(&callbackInfo{data: data, ptr: res})
	res.c, res.l, res.Handle = t, loop.GetNativeLoop(), Handle{(*C.uv_handle_t)(unsafe.Pointer(t)), t.data, res}
	if r := C.uv_check_init(loop.GetNativeLoop(), t); r != 0 {
		C.free(unsafe.Pointer(t))
		return nil, ParseUvErr(r)
	}

	return res, nil
}

// Start (uv_prepare_start) start the timer. timeout and repeat are in milliseconds.
func (c *UvCheck) Start(cb func(*Handle)) C.int {
	cbi := (*callbackInfo)(c.c.data)
	cbi.check_cb = cb

	return uv_check_start(c.c)
}

// Stop (uv_prepare_stop) the timer, the callback will not be called anymore.
func (c *UvCheck) Stop() C.int {
	return C.uv_check_stop(c.c)
}

// Freemem freemem handle
func (c *UvCheck) Freemem() {
	C.free(unsafe.Pointer(c.c))
}

// GetAsyncHandle get handle
func (c *UvCheck) GetCheckHandle() *C.uv_check_t {
	return c.c
}
