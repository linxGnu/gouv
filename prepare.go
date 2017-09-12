package gouv

// #cgo pkg-config: libuv
/*
#include <uv.h>
#include <stdlib.h>
uv_prepare_t* mallocPrepareT() {
	return (uv_prepare_t*)malloc(sizeof(uv_prepare_t));
}
*/
import "C"
import "unsafe"

// UvPrepare handles will run the given callback once per loop iteration, right before polling for i/o.
type UvPrepare struct {
	p *C.uv_prepare_t
	l *C.uv_loop_t
	Handle
}

// UvPrepareInit initialize the prepare handle
func UvPrepareInit(loop *UvLoop, data interface{}) (*UvPrepare, error) {
	t := C.mallocPrepareT()

	if loop == nil {
		loop = UvLoopDefault()
	}

	res := &UvPrepare{}
	t.data = unsafe.Pointer(&callbackInfo{data: data, ptr: res})
	res.p, res.l, res.Handle = t, loop.GetNativeLoop(), Handle{(*C.uv_handle_t)(unsafe.Pointer(t)), t.data, res}
	if r := C.uv_prepare_init(loop.GetNativeLoop(), t); r != 0 {
		return nil, ParseUvErr(r)
	}

	return res, nil
}

// Start (uv_prepare_start) start the timer. timeout and repeat are in milliseconds.
func (t *UvPrepare) Start(cb func(*Handle)) (err error) {
	cbi := (*callbackInfo)(t.p.data)
	cbi.prepare_cb = cb

	if r := uv_prepare_start(t.p); r != 0 {
		err = ParseUvErr(r)
		return
	}

	return
}

// Stop (uv_prepare_stop) the timer, the callback will not be called anymore.
func (t *UvPrepare) Stop() (err error) {
	if r := C.uv_prepare_stop(t.p); r != 0 {
		err = ParseUvErr(r)
		return
	}

	return
}

// Freemem freemem of prepare
func (t *UvPrepare) Freemem() {
	C.free(unsafe.Pointer(t.p))
}
