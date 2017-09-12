package gouv

// #cgo pkg-config: libuv
/*
#include <uv.h>
#include <stdlib.h>
uv_shutdown_t* mallocShutdownT() {
	return (uv_shutdown_t*)malloc(sizeof(uv_shutdown_t));
}
*/
import "C"
import (
	"unsafe"
)

//UvShutdown shutdown request type
type UvShutdown struct {
	s *C.uv_shutdown_t
	Handle
}

// NewUvShutdown new uv_shutdown req
func NewUvShutdown(data interface{}) *UvShutdown {
	t := C.mallocShutdownT()

	res := &UvShutdown{}
	t.data = unsafe.Pointer(&callbackInfo{data: data, ptr: res})
	res.s, res.Handle = t, Handle{(*C.uv_handle_t)(unsafe.Pointer(t)), t.data, res}

	return res
}

// UvStream handles provide an abstraction of a duplex communication channel. uv_stream_t is an abstract type, libuv provides 3 stream implementations in the form of uv_tcp_t, uv_pipe_t and uv_tty_t.
type UvStream struct {
	s *C.uv_stream_t
	Handle
}

// Shutdown (uv_shutdown) shutdown the outgoing (write) side of a duplex stream. It waits for pending write requests to complete. The handle should refer to a initialized stream. req should be an uninitialized shutdown request struct. The cb is called after shutdown is complete.
func (s *UvStream) Shutdown(req *C.uv_shutdown_t, cb func(*Request, int)) (err error) {
	cbi := (*callbackInfo)(req.data)
	cbi.shutdown_cb = cb

	if r := uv_shutdown(req, s.s); r != 0 {
		err = ParseUvErr(r)
		return
	}

	return
}

// Listen start listening for incoming connections. backlog indicates the number of connections the kernel might queue, same as listen(2). When a new incoming connection is received the uv_connection_cb callback is called.
func (s *UvStream) Listen(backlog int, cb func(*Handle, int)) (err error) {
	cbi := (*callbackInfo)(s.s.data)
	cbi.connection_cb = cb

	if r := uv_listen(s.s, backlog); r != 0 {
		err = ParseUvErr(r)
		return
	}

	return
}

// ServerAccept this call is used in conjunction with uv_listen() to accept incoming connections.
// Call this function after receiving a uv_connection_cb to accept the connection.
// Before calling this function the client handle must be initialized. < 0 return value indicates an error.
func (s *UvStream) ServerAccept(c *C.uv_stream_t) (err error) {
	if r := uv_accept(s.s, c); r != 0 {
		err = ParseUvErr(r)
		return
	}

	return nil
}

// ClientAccept this call is used in conjunction with uv_listen() to accept incoming connections.
// Call this function after receiving a uv_connection_cb to accept the connection.
// Before calling this function the server handle must be initialized. < 0 return value indicates an error.
func (s *UvStream) ClientAccept(c *C.uv_stream_t) (err error) {
	if r := uv_accept(c, s.s); r != 0 {
		err = ParseUvErr(r)
		return
	}

	return nil
}

// ReadStart (uv_read_start) read data from an incoming stream. The uv_read_cb callback will be made several times until there is no more data to read or uv_read_stop() is called.
func (s *UvStream) ReadStart(cb func(*Handle, *C.uv_buf_t, C.ssize_t)) (err error) {
	cbi := (*callbackInfo)(s.s.data)
	cbi.read_cb = cb

	if r := uv_read_start(s.s); r != 0 {
		err = ParseUvErr(r)
		return
	}

	return
}

// ReadStop (uv_read_stop) stop reading data from the stream. The uv_read_cb callback will no longer be called.
// This function is idempotent and may be safely called on a stopped stream.
func (s *UvStream) ReadStop() (err error) {
	if r := uv_read_stop(s.s); r != 0 {
		err = ParseUvErr(r)
		return
	}

	return
}

// Write (uv_write) write data to stream. Buffers are written in order.
func (s *UvStream) Write(req *C.uv_write_t, buf *C.uv_buf_t, bufcnt int, cb func(*Request, int)) C.int {
	cbi := (*callbackInfo)(req.data)
	cbi.write_cb = cb

	return uv_write(req, s.s, buf, bufcnt)
}

// Write2 (uv_write2) Extended write function for sending handles over a pipe. The pipe must be initialized with ipc == 1.
func (s *UvStream) Write2(req *C.uv_write_t, stream *C.uv_stream_t, buf *C.uv_buf_t, bufcnt int, sendHandle *C.uv_stream_t, cb func(*Request, int)) C.int {
	cbi := (*callbackInfo)(req.data)
	cbi.write_cb = cb

	return uv_write2(req, s.s, buf, bufcnt, sendHandle)
}

// TryWrite (uv_try_write) same as uv_write(), but won’t queue a write request if it can’t be completed immediately.
func (s *UvStream) TryWrite(buf *C.uv_buf_t, bufcnt int) C.int {
	return uv_try_write(s.s, buf, bufcnt)
}

// IsReadable (uv_is_readable) returns if the stream is readable.
func (s *UvStream) IsReadable() bool {
	return uv_is_readable(s.s)
}

// IsWritable (uv_is_writable) returns if the stream is writable.
func (s *UvStream) IsWritable() bool {
	return uv_is_writable(s.s)
}

// SetBlocking (uv_stream_set_blocking) enable or disable blocking mode for a stream.
// When blocking mode is enabled all writes complete synchronously.
// The interface remains unchanged otherwise, e.g. completion or failure of the operation will still be reported through a callback which is made asynchronously.
func (s *UvStream) SetBlocking(blocking int) (err error) {
	if r := uv_stream_set_blocking(s.s, blocking); r != 0 {
		err = ParseUvErr(r)
		return
	}

	return
}

// Freemem freemem of base stream
func (s *UvStream) Freemem() {
	C.free(unsafe.Pointer(s.s))
}
