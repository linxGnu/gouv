package gouv

// #cgo pkg-config: libuv
/*
#include <uv.h>
#include <stdlib.h>
uv_async_t* mallocAsyncT() {
	return (uv_async_t*)malloc(sizeof(uv_async_t));
}
*/
import "C"
import "unsafe"

// UvAsync handles allow the user to “wakeup” the event loop and get a callback called from another thread.
type UvAsync struct {
	Async *C.uv_async_t
	Loop  *C.uv_loop_t
	Handle
}

// UvAsyncInit initialize the prepare handle
func UvAsyncInit(loop *UvLoop, data interface{}, cb func(*Handle)) (*UvAsync, error) {
	t := C.mallocAsyncT()

	if loop == nil {
		loop = UvLoopDefault()
	}

	res := &UvAsync{}
	t.data = unsafe.Pointer(&callbackInfo{data: data, async_cb: cb, ptr: res})
	res.Async, res.Loop, res.Handle = t, loop.GetNativeLoop(), Handle{(*C.uv_handle_t)(unsafe.Pointer(t)), t.data, res}
	if r := uv_async_init(loop.GetNativeLoop(), t); r != 0 {
		C.free(unsafe.Pointer(t))
		return nil, ParseUvErr(r)
	}

	return res, nil
}

// Send (uv_async_send) wake up the event loop and call the async handle’s callback.
func (a *UvAsync) Send() C.int {
	return C.uv_async_send(a.Async)
}

// Freemem freemem handle
func (a *UvAsync) Freemem() {
	C.free(unsafe.Pointer(a.Async))
}
