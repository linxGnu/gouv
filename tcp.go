package gouv

// #cgo pkg-config: libuv
/*
#include <uv.h>
#include <stdlib.h>
uv_tcp_t* mallocTCPT() {
	return (uv_tcp_t*)malloc(sizeof(uv_tcp_t));
}

char* testWriteTCP(uv_stream_t *client, ssize_t nread, uv_buf_t* buf) {
	char* tmp;
	tmp = malloc(nread + 1);
	memcpy(tmp, buf->base, nread);
	tmp[nread] = '\0';
	uv_write_t* write_req;

	uv_buf_t bufw = uv_buf_init(tmp, strlen(tmp) + 1);

	write_req = malloc(sizeof(uv_write_t));
	uv_write(write_req, client, &bufw, 1, NULL);

	return tmp;
}
*/
import "C"
import (
	"fmt"
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
			return nil, ParseUvErr(r)
		}
	} else {
		if r := C.uv_tcp_init_ex(loop.GetNativeLoop(), t, C.uint(*flags)); r != 0 {
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
func (t *UvTCP) Open(sock C.uv_os_sock_t) (err error) {
	if r := C.uv_tcp_open(t.t, sock); r != 0 {
		err = ParseUvErr(r)
		return
	}

	return
}

// NoDelay (uv_tcp_nodelay) enable TCP_NODELAY, which disables Nagle’s algorithm.
func (t *UvTCP) NoDelay(enable int) (err error) {
	if r := C.uv_tcp_nodelay(t.t, C.int(enable)); r != 0 {
		err = ParseUvErr(r)
		return
	}

	return
}

// KeepAlive (uv_tcp_keepalive) enable/disable TCP keep-alive. delay is the initial delay in seconds, ignored when enable is zero.
func (t *UvTCP) KeepAlive(enable int, delay uint) (err error) {
	if r := C.uv_tcp_keepalive(t.t, C.int(enable), C.uint(delay)); r != 0 {
		err = ParseUvErr(r)
		return
	}

	return
}

// SimultaneousAccepts (uv_tcp_simultaneous_accepts) enable/disable simultaneous asynchronous accept requests that are queued by the operating system when listening for new TCP connections.
// This setting is used to tune a TCP server for the desired performance. Having simultaneous accepts can significantly improve the rate of accepting connections (which is why it is enabled by default) but may lead to uneven load distribution in multi-process setups.
func (t *UvTCP) SimultaneousAccepts(enable int) (err error) {
	if r := C.uv_tcp_simultaneous_accepts(t.t, C.int(enable)); r != 0 {
		err = ParseUvErr(r)
		return
	}

	return
}

// Bind (uv_tcp_bind) bind the handle to an address and port. addr should point to an initialized struct sockaddr_in or struct sockaddr_in6.
// When the port is already taken, you can expect to see an UV_EADDRINUSE error from either uv_tcp_bind(), uv_listen() or uv_tcp_connect().
// That is, a successful call to this function does not guarantee that the call to uv_listen() or uv_tcp_connect() will succeed as well.
// flags can contain UV_TCP_IPV6ONLY, in which case dual-stack support is disabled and only IPv6 is used.
func (t *UvTCP) Bind(sockAddr SockaddrIn, flags uint) (err error) {
	if r := C.uv_tcp_bind(t.t, sockAddr.GetSockAddr(), C.uint(flags)); r != 0 {
		err = ParseUvErr(r)
		return
	}

	return
}

// Connect (uv_tcp_connect) establish an IPv4 or IPv6 TCP connection. Provide an initialized TCP handle and an uninitialized uv_connect_t. addr should point to an initialized struct sockaddr_in or struct sockaddr_in6.
// The callback is made when the connection has been established or when a connection error happened.
func (t *UvTCP) Connect(req *C.uv_connect_t, sockAddr SockaddrIn, cb func(*Request, int)) (err error) {
	cbi := (*callbackInfo)(req.data)
	cbi.connect_cb = cb

	if r := uv_tcp_connect(req, t.t, sockAddr.GetSockAddr()); r != 0 {
		err = ParseUvErr(r)
		return
	}

	return
}

func sampleTCPReadHandling(h *Handle, buf *C.uv_buf_t, nRead C.ssize_t) {
	fmt.Println(C.GoString(C.testWriteTCP(h.ptr.(*UvTCP).s, nRead, buf)))
}

// TODO: uv_tcp_getsockname
// TODO: uv_tcp_getpeername
