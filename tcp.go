package gouv

// #cgo pkg-config: libuv
/*
#include <uv.h>
#include <stdlib.h>
uv_tcp_t* mallocTCPT() {
	return (uv_tcp_t*)malloc(sizeof(uv_tcp_t));
}
*/
import "C"
import (
	"unsafe"
)

// UvTCP handles are used to represent both TCP streams and servers.
type UvTCP struct {
	t *C.uv_tcp_t
	l *C.uv_loop_t
	UvStream
}

// TCPInit (uv_tcp_init) initialize the handle. No socket is created as of yet.
func TCPInit(loop *UvLoop, flags *uint, data interface{}) (*UvTCP, error) {
	t := C.mallocTCPT()

	if loop == nil {
		loop = UvLoopDefault()
	}

	if flags == nil {
		if r := C.uv_tcp_init(loop.GetNativeLoop(), t); r != 0 {
			C.free(unsafe.Pointer(t))
			return nil, ParseUvErr(r)
		}
	} else {
		if r := C.uv_tcp_init_ex(loop.GetNativeLoop(), t, C.uint(*flags)); r != 0 {
			C.free(unsafe.Pointer(t))
			return nil, ParseUvErr(r)
		}
	}

	res := &UvTCP{}
	t.data = unsafe.Pointer(&callbackInfo{ptr: res, data: data})
	res.s, res.l, res.t, res.Handle = (*C.uv_stream_t)(unsafe.Pointer(t)), loop.GetNativeLoop(), t, Handle{(*C.uv_handle_t)(unsafe.Pointer(t)), t.data, res}

	return res, nil
}

// Open (uv_tcp_open) open an existing file descriptor or SOCKET as a TCP handle.
// Note: the passed file descriptor or SOCKET is not checked for its type, but it’s required that it represents a valid stream socket.
func (t *UvTCP) Open(sock C.uv_os_sock_t) C.int {
	return C.uv_tcp_open(t.t, sock)
}

// NoDelay (uv_tcp_nodelay) enable TCP_NODELAY, which disables Nagle’s algorithm.
func (t *UvTCP) NoDelay(enable int) C.int {
	return C.uv_tcp_nodelay(t.t, C.int(enable))
}

// KeepAlive (uv_tcp_keepalive) enable/disable TCP keep-alive. delay is the initial delay in seconds, ignored when enable is zero.
func (t *UvTCP) KeepAlive(enable int, delay uint) C.int {
	return C.uv_tcp_keepalive(t.t, C.int(enable), C.uint(delay))
}

// SimultaneousAccepts (uv_tcp_simultaneous_accepts) enable/disable simultaneous asynchronous accept requests that are queued by the operating system when listening for new TCP connections.
// This setting is used to tune a TCP server for the desired performance. Having simultaneous accepts can significantly improve the rate of accepting connections (which is why it is enabled by default) but may lead to uneven load distribution in multi-process setups.
func (t *UvTCP) SimultaneousAccepts(enable int) C.int {
	return C.uv_tcp_simultaneous_accepts(t.t, C.int(enable))
}

// Bind (uv_tcp_bind) bind the handle to an address and port. addr should point to an initialized struct sockaddr_in or struct sockaddr_in6.
// When the port is already taken, you can expect to see an UV_EADDRINUSE error from either uv_tcp_bind(), uv_listen() or uv_tcp_connect().
// That is, a successful call to this function does not guarantee that the call to uv_listen() or uv_tcp_connect() will succeed as well.
// flags can contain UV_TCP_IPV6ONLY, in which case dual-stack support is disabled and only IPv6 is used.
func (t *UvTCP) Bind(sockAddr SockaddrIn, flags uint) C.int {
	return C.uv_tcp_bind(t.t, sockAddr.GetSockAddr(), C.uint(flags))
}

// Connect (uv_tcp_connect) establish an IPv4 or IPv6 TCP connection. Provide an initialized TCP handle and an uninitialized uv_connect_t. addr should point to an initialized struct sockaddr_in or struct sockaddr_in6.
// The callback is made when the connection has been established or when a connection error happened.
func (t *UvTCP) Connect(req *UvConnect, sockAddr SockaddrIn, cb func(*Request, int)) C.int {
	cbi := (*callbackInfo)(req.c.data)
	cbi.connect_cb = cb
	cbi.ptr = t

	return uv_tcp_connect(req.c, t.t, sockAddr.GetSockAddr())
}

// GetTCPHandle get handle
func (t *UvTCP) GetTCPHandle() *C.uv_tcp_t {
	return t.t
}

// TODO: uv_tcp_getsockname
// TODO: uv_tcp_getpeername
