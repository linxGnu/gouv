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
	a *C.uv_async_t
	l *C.uv_loop_t
	Handle
}

// UvAsyncInit initialize the prepare handle
func UvAsyncInit(loop *UvLoop, data interface{}, cb func(*Handle)) (timer *UvAsync, err error) {
	t := C.mallocAsyncT()
	t.data = unsafe.Pointer(&callback_info{data: data, async_cb: cb})

	if loop == nil {
		loop = UvLoopDefault()
	}

	if r := uv_async_init(loop.GetNativeLoop(), t); r != 0 {
		return nil, ParseUvErr(r)
	}

	return &UvAsync{t, loop.GetNativeLoop(), Handle{(*C.uv_handle_t)(unsafe.Pointer(t)), t.data}}, nil
}

// Send (uv_async_send) wake up the event loop and call the async handle’s callback.
func (t *UvAsync) Send() error {
	if r := C.uv_async_send(t.a); r != 0 {
		return ParseUvErr(r)
	}

	return nil
}

// Freemem freemem of prepare
func (t *UvAsync) Freemem() {
	C.free(unsafe.Pointer(t.a))
}
