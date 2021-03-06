package gouv

// #cgo pkg-config: libuv
/*
#include <uv.h>
#include <stdlib.h>
uv_poll_t* mallocPollT() {
	return (uv_poll_t*)malloc(sizeof(uv_poll_t));
}
*/
import "C"
import "unsafe"

// UvPoll handles are used to watch file descriptors for readability, writability and disconnection similar to the purpose of poll(2).
// The purpose of poll handles is to enable integrating external libraries that rely on the event loop to signal it about the socket status changes, like c-ares or libssh2.
// Using uv_poll_t for any other purpose is not recommended; uv_tcp_t, uv_udp_t, etc. provide an implementation that is faster and more scalable than what can be achieved with uv_poll_t, especially on Windows.
type UvPoll struct {
	Poll *C.uv_poll_t
	Loop *C.uv_loop_t
	Handle
}

// UvPollInit (uv_poll_init) initialize the handle using a file descriptor.
func UvPollInit(loop *UvLoop, fd int, data interface{}) (*UvPoll, error) {
	t := C.mallocPollT()

	if loop == nil {
		loop = UvLoopDefault()
	}

	res := &UvPoll{}
	t.data = unsafe.Pointer(&callbackInfo{data: data, ptr: res})
	res.Poll, res.Loop, res.Handle = t, loop.GetNativeLoop(), Handle{(*C.uv_handle_t)(unsafe.Pointer(t)), t.data, res}
	if r := C.uv_poll_init(loop.GetNativeLoop(), t, C.int(fd)); r != 0 {
		C.free(unsafe.Pointer(t))
		return nil, ParseUvErr(r)
	}

	return res, nil
}

// UvPollInitSocket (uv_poll_init_socket) initialize the handle using a socket descriptor. On Unix this is identical to uv_poll_init().
// On windows it takes a SOCKET handle.
func UvPollInitSocket(loop *UvLoop, socket C.uv_os_sock_t, data interface{}) (*UvPoll, error) {
	t := C.mallocPollT()

	if loop == nil {
		loop = UvLoopDefault()
	}

	res := &UvPoll{}
	t.data = unsafe.Pointer(&callbackInfo{data: data, ptr: res})
	res.Poll, res.Loop, res.Handle = t, loop.GetNativeLoop(), Handle{(*C.uv_handle_t)(unsafe.Pointer(t)), t.data, res}
	if r := C.uv_poll_init_socket(loop.GetNativeLoop(), t, socket); r != 0 {
		C.free(unsafe.Pointer(t))
		return nil, ParseUvErr(r)
	}

	return res, nil
}

// Start (uv_poll_start) starts polling the file descriptor. events is a bitmask made up of UV_READABLE, UV_WRITABLE, UV_PRIORITIZED and UV_DISCONNECT.
// As soon as an event is detected the callback will be called with status set to 0, and the detected events set on the events field.
func (p *UvPoll) Start(event int, cb func(*Handle, int, int)) C.int {
	cbi := (*callbackInfo)(p.Poll.data)
	cbi.poll_cb = cb

	return uv_poll_start(p.Poll, event)
}

// Stop (uv_poll_stop) stop polling the file descriptor, the callback will no longer be called.
func (p *UvPoll) Stop() C.int {
	return C.uv_poll_stop(p.Poll)
}

// Freemem freemem handle
func (p *UvPoll) Freemem() {
	C.free(unsafe.Pointer(p.Poll))
}
