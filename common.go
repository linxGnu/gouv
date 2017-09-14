package gouv

// #cgo pkg-config: libuv
/*
#include "common.h"
*/
import "C"
import (
	"unsafe"
)

// UV_HANDLE_TYPE the kind of the libuv handle.
type UV_HANDLE_TYPE int

const (
	UV_UNKNOWN_HANDLE  UV_HANDLE_TYPE = 0
	UV_ASYNC           UV_HANDLE_TYPE = 1
	UV_CHECK           UV_HANDLE_TYPE = 2
	UV_FS_EVENT        UV_HANDLE_TYPE = 3
	UV_FS_POLL         UV_HANDLE_TYPE = 4
	UV_HANDLE          UV_HANDLE_TYPE = 5
	UV_IDLE            UV_HANDLE_TYPE = 6
	UV_NAMED_PIPE      UV_HANDLE_TYPE = 7
	UV_POLL            UV_HANDLE_TYPE = 8
	UV_PREPARE         UV_HANDLE_TYPE = 9
	UV_PROCESS         UV_HANDLE_TYPE = 10
	UV_STREAM          UV_HANDLE_TYPE = 11
	UV_TCP             UV_HANDLE_TYPE = 12
	UV_TIMER           UV_HANDLE_TYPE = 13
	UV_TTY             UV_HANDLE_TYPE = 14
	UV_UDP             UV_HANDLE_TYPE = 15
	UV_SIGNAL          UV_HANDLE_TYPE = 16
	UV_FILE            UV_HANDLE_TYPE = 17
	UV_HANDLE_TYPE_MAX UV_HANDLE_TYPE = 18
)

type UV_REQ_TYPE int

const (
	UV_UNKNOWN_REQ      UV_REQ_TYPE = 0
	UV_REQ              UV_REQ_TYPE = 1
	UV_CONNECT          UV_REQ_TYPE = 2
	UV_WRITE            UV_REQ_TYPE = 3
	UV_SHUTDOWN         UV_REQ_TYPE = 4
	UV_UDP_SEND         UV_REQ_TYPE = 5
	UV_FS               UV_REQ_TYPE = 6
	UV_WORK             UV_REQ_TYPE = 7
	UV_GETADDRINFO      UV_REQ_TYPE = 8
	UV_GETNAMEINFO      UV_REQ_TYPE = 9
	UV_REQ_TYPE_PRIVATE UV_REQ_TYPE = 10
	UV_REQ_TYPE_MAX     UV_REQ_TYPE = 11
)

// UV_RUN_MODE mode used to run the loop with uv_run()
type UV_RUN_MODE int

const (
	// UV_RUN_DEFAULT runs the event loop until there are no more active and referenced handles or requests. Returns non-zero if uv_stop() was called and there are still active handles or requests. Returns zero in all other cases.
	UV_RUN_DEFAULT UV_RUN_MODE = 0

	// UV_RUN_ONCE poll for i/o once. Note that this function blocks if there are no pending callbacks. Returns zero when done (no active handles or requests left), or non-zero if more callbacks are expected (meaning you should run the event loop again sometime in the future).
	UV_RUN_ONCE UV_RUN_MODE = 1

	// UV_RUN_NOWAIT poll for i/o once but don’t block if there are no pending callbacks. Returns zero if done (no active handles or requests left), or non-zero if more callbacks are expected (meaning you should run the event loop again sometime in the future).
	UV_RUN_NOWAIT UV_RUN_MODE = 2
)

// UV_POLL_EVENT poll event types
type UV_POLL_EVENT int

const (
	// UV_READABLE fd is readable
	UV_READABLE UV_POLL_EVENT = 1

	// UV_WRITABLE fd is writable
	UV_WRITABLE UV_POLL_EVENT = 2

	// UV_DISCONNECT event is optional in the sense that it may not be reported and the user is free to ignore it, but it can help optimize the shutdown path because an extra read or write call might be avoided.
	UV_DISCONNECT UV_POLL_EVENT = 4

	// UV_PRIORITIZED event is used to watch for sysfs interrupts or TCP out-of-band messages.
	UV_PRIORITIZED UV_POLL_EVENT = 8
)

// UvStdioFlags flags specifying how a stdio should be transmitted to the child process.
type UV_STDIO_FLAGS int

const (
	UV_IGNORE         UV_STDIO_FLAGS = 0x00
	UV_CREATE_PIPE    UV_STDIO_FLAGS = 0x01
	UV_INHERIT_FD     UV_STDIO_FLAGS = 0x02
	UV_INHERIT_STREAM UV_STDIO_FLAGS = 0x04
	/*
	 * When UV_CREATE_PIPE is specified, UV_READABLE_PIPE and UV_WRITABLE_PIPE
	 * determine the direction of flow, from the child process' perspective. Both
	 * flags may be specified to create a duplex data stream.
	 */
	UV_READABLE_PIPE UV_STDIO_FLAGS = 0x10
	UV_WRITABLE_PIPE UV_STDIO_FLAGS = 0x20
)

// UV_PROCESS_FLAGS flags to be set on the flags field of uv_process_options_t.
type UV_PROCESS_FLAGS uint

const (
	/*
	 * Set the child process' user id.
	 */
	UV_PROCESS_SETUID UV_PROCESS_FLAGS = (1 << 0)
	/*
	 * Set the child process' group id.
	 */
	UV_PROCESS_SETGID UV_PROCESS_FLAGS = (1 << 1)
	/*
	 * Do not wrap any arguments in quotes, or perform any other escaping, when
	 * converting the argument list into a command line string. This option is
	 * only meaningful on Windows systems. On Unix it is silently ignored.
	 */
	UV_PROCESS_WINDOWS_VERBATIM_ARGUMENTS UV_PROCESS_FLAGS = (1 << 2)
	/*
	 * Spawn the child process in a detached state - this will make it a process
	 * group leader, and will effectively enable the child to keep running after
	 * the parent exits. Note that the child process will still keep the
	 * parent's event loop alive unless the parent process calls uv_unref() on
	 * the child's process handle.
	 */
	UV_PROCESS_DETACHED UV_PROCESS_FLAGS = (1 << 3)
	/*
	 * Hide the subprocess console window that would normally be created. This
	 * option is only meaningful on Windows systems. On Unix it is silently
	 * ignored.
	 */
	UV_PROCESS_WINDOWS_HIDE UV_PROCESS_FLAGS = (1 << 4)
)

// UV_TTY_MODE TTY mode type
type UV_TTY_MODE int

const (
	/* UV_TTY_MODE_NORMAL initial/normal terminal mode */
	UV_TTY_MODE_NORMAL UV_TTY_MODE = 0

	// UV_TTY_MODE_RAW raw input mode (On Windows, ENABLE_WINDOW_INPUT is also enabled)
	UV_TTY_MODE_RAW UV_TTY_MODE = 1

	// UV_TTY_MODE_IO binary-safe I/O mode for IPC (Unix-only)
	UV_TTY_MODE_IO UV_TTY_MODE = 2
)

// Request (uv_req_t) is the base type for all libuv request types.
type Request struct {
	r      *C.uv_req_t
	Handle *Handle
}

// Handle (uv_handle_t) is the base type for all libuv handle types.
type Handle struct {
	h    *C.uv_handle_t
	Data interface{}
	Ptr  interface{}
}

type callbackInfo struct {
	connection_cb func(*Handle, int)
	connect_cb    func(*Request, int)
	read_cb       func(*Handle, *C.uv_buf_t, C.ssize_t)
	udp_recv_cb   func(*Handle, []byte, *C.struct_sockaddr, uint)
	write_cb      func(*Request, int)
	udp_send_cb   func(*Request, int)
	close_cb      func(*Handle)
	prepare_cb    func(*Handle)
	poll_cb       func(*Handle, int, int)
	check_cb      func(*Handle)
	shutdown_cb   func(*Request, int)
	timer_cb      func(*Handle, int)
	signal_cb     func(*Handle, int)
	idle_cb       func(*Handle, int)
	exit_cb       func(*Handle, int, int)
	async_cb      func(*Handle)
	data          interface{}
	ptr           interface{}
}

func (handle *Handle) Close(cb func(*Handle)) {
	cbi := (*callbackInfo)(handle.h.data)
	cbi.close_cb = cb
	uv_close(handle.h)
}

// IsActive (uv_is_active) returns if the handle is active.
func (handle *Handle) IsActive() bool {
	return uv_is_active(handle.h)
}

// IsClosing (uv_is_closing) returns if the handle is closing or closed.
func (handle *Handle) IsClosing() bool {
	return uv_is_closing(handle.h)
}

// uv_tcp_bind (uv_tcp_bind) bind the handle to an address and port. addr should point to an initialized struct sockaddr_in or struct sockaddr_in6.
func uv_tcp_bind(tcp *C.uv_tcp_t, sa *C.struct_sockaddr, flags uint) C.int {
	return C.uv_tcp_bind(tcp, sa, C.uint(flags))
}

// uv_tcp_connect (uv_tcp_connect) establish an IPv4 or IPv6 TCP connection. Provide an initialized TCP handle and an uninitialized uv_connect_t. addr should point to an initialized struct sockaddr_in or struct sockaddr_in6.
func uv_tcp_connect(req *C.uv_connect_t, tcp *C.uv_tcp_t, sa *C.struct_sockaddr) C.int {
	return C._uv_tcp_connect(req, tcp, sa)
}

// uv_pipe_connect (uv_pipe_connect) connect to the Unix domain socket or the named pipe.
func uv_pipe_connect(req *C.uv_connect_t, pipe *C.uv_pipe_t, name string) {
	pname := C.CString(name)
	defer C.free(unsafe.Pointer(pname))

	C._uv_pipe_connect(req, pipe, pname)
}

// uv_is_active (uv_is_active) returns if the handle is active.
func uv_is_active(handle *C.uv_handle_t) bool {
	return C.uv_is_active(handle) != 0
}

// uv_is_closing (uv_is_closing) returns if the handle is closing or closed.
func uv_is_closing(handle *C.uv_handle_t) bool {
	return C.uv_is_closing(handle) != 0
}

// uv_close (uv_close) request handle to be closed. close_cb will be called asynchronously after this call. This MUST be called on each handle before memory is released.
func uv_close(handle *C.uv_handle_t) {
	C._uv_close(handle)
}

// uv_ref (uv_ref) reference the given handle. References are idempotent, that is, if a handle is already referenced calling this function again will have no effect.
func uv_ref(h *C.uv_handle_t) {
	C.uv_ref(h)
}

// uv_unref (uv_unref) unreference the given handle. References are idempotent, that is, if a handle is not referenced calling this function again will have no effect.
func uv_unref(h *C.uv_handle_t) {
	C.uv_unref(h)
}

// UvHauv_has_refsRef (uv_has_ref) returns if the handle referenced.
func uv_has_ref(h *C.uv_handle_t) bool {
	return C.uv_has_ref(h) != 0
}

// uv_handle_size (uv_handle_size) returns the size of the given handle type. Useful for FFI binding writers who don’t want to know the structure layout.
func uv_handle_size(t UV_HANDLE_TYPE) C.size_t {
	return C.uv_handle_size(C.uv_handle_type(t))
}

// uv_send_buffer_size (uv_send_buffer_size) gets or sets the size of the send buffer that the operating system uses for the socket.
// if *value == 0, it will return the current send buffer size, otherwise it will use *value to set the new send buffer size.
// This function works for TCP, pipe and UDP handles on Unix and for TCP and UDP handles on Windows.
func uv_send_buffer_size(h *C.uv_handle_t, value *C.int) C.int {
	return C.uv_send_buffer_size(h, value)
}

// uv_recv_buffer_size (uv_recv_buffer_size) gets or sets the size of the receive buffer that the operating system uses for the socket.
// If *value == 0, it will return the current receive buffer size, otherwise it will use *value to set the new receive buffer size.
// This function works for TCP, pipe and UDP handles on Unix and for TCP and UDP handles on Windows.
func uv_recv_buffer_size(h *C.uv_handle_t, value *C.int) C.int {
	return C.uv_recv_buffer_size(h, value)
}

// uv_fileno (uv_fileno) gets the platform dependent file descriptor equivalent.
// The following handles are supported: TCP, pipes, TTY, UDP and poll. Passing any other handle type will fail with UV_EINVAL.
// If a handle doesn’t have an attached file descriptor yet or the handle itself has been closed, this function will return UV_EBADF.
func uv_fileno(h *C.uv_handle_t, fd *C.uv_os_fd_t) C.int {
	return C.uv_fileno(h, fd)
}

// uv_cancel (uv_cancel) cancel a pending request. Fails if the request is executing or has finished executing.
// Returns 0 on success, or an error code < 0 on failure.
// Only cancellation of uv_fs_t, uv_getaddrinfo_t, uv_getnameinfo_t and uv_work_t requests is currently supported.
func uv_cancel(req *C.uv_req_t) C.int {
	return C.uv_cancel(req)
}

// uv_req_size (uv_req_size) returns the size of the given request type. Useful for FFI binding writers who don’t want to know the structure layout.
func uv_req_size(t UV_REQ_TYPE) C.size_t {
	return C.uv_req_size(C.uv_req_type(t))
}

// uv_shutdown (uv_shutdown) shutdown the outgoing (write) side of a duplex stream. It waits for pending write requests to complete. The handle should refer to a initialized stream. req should be an uninitialized shutdown request struct. The cb is called after shutdown is complete.
func uv_shutdown(req *C.uv_shutdown_t, stream *C.uv_stream_t) C.int {
	return C._uv_shutdown(req, stream)
}

// uv_listen (uv_listen) start listening for incoming connections. backlog indicates the number of connections the kernel might queue, same as listen(2). When a new incoming connection is received the uv_connection_cb callback is called.
func uv_listen(stream *C.uv_stream_t, backlog int) C.int {
	return C._uv_listen(stream, C.int(backlog))
}

// uv_accept (uv_accept) this call is used in conjunction with uv_listen() to accept incoming connections. Call this function after receiving a uv_connection_cb to accept the connection. Before calling this function the client handle must be initialized. < 0 return value indicates an error.
func uv_accept(stream *C.uv_stream_t, client *C.uv_stream_t) C.int {
	return C.uv_accept(stream, client)
}

// uv_read_start (uv_read_start) read data from an incoming stream. The uv_read_cb callback will be made several times until there is no more data to read or uv_read_stop() is called.
func uv_read_start(stream *C.uv_stream_t) C.int {
	return C._uv_read_start(stream)
}

// uv_read_stop (uv_read_stop) stop reading data from the stream. The uv_read_cb callback will no longer be called.
// This function is idempotent and may be safely called on a stopped stream.
func uv_read_stop(stream *C.uv_stream_t) C.int {
	return C.uv_read_stop(stream)
}

// uv_write (uv_write) write data to stream. Buffers are written in order.
// Note: the memory pointed to by the buffers must remain valid until the callback gets called. This also holds for uv_write2().
func uv_write(req *C.uv_write_t, stream *C.uv_stream_t, buf *C.uv_buf_t, bufcnt int) C.int {
	return C._uv_write(req, stream, buf, C.int(bufcnt))
}

// uv_write2 (uv_write2) Extended write function for sending handles over a pipe. The pipe must be initialized with ipc == 1.
// Note: send_handle must be a TCP socket or pipe, which is a server or a connection (listening or connected state). Bound sockets or pipes will be assumed to be servers.
func uv_write2(req *C.uv_write_t, stream *C.uv_stream_t, buf *C.uv_buf_t, bufcnt int, send_handle *C.uv_stream_t) C.int {
	return C._uv_write2(req, stream, buf, C.int(bufcnt), send_handle)
}

// uv_try_write (uv_try_write) write data to stream. Buffers are written in order.
func uv_try_write(stream *C.uv_stream_t, buf *C.uv_buf_t, bufcnt int) C.int {
	return C._uv_try_write(stream, buf, C.int(bufcnt))
}

// uv_is_readable (uv_is_readable) returns if the stream is readable.
func uv_is_readable(stream *C.uv_stream_t) bool {
	return C.uv_is_readable(stream) == 1
}

// uv_is_writable (uv_is_writable) returns if the stream is writable.
func uv_is_writable(stream *C.uv_stream_t) bool {
	return C.uv_is_writable(stream) == 1
}

// uv_stream_set_blocking (uv_stream_set_blocking) enable or disable blocking mode for a stream.
// When blocking mode is enabled all writes complete synchronously. The interface remains unchanged otherwise, e.g. completion or failure of the operation will still be reported through a callback which is made asynchronously.
func uv_stream_set_blocking(stream *C.uv_stream_t, blocking int) C.int {
	return C.uv_stream_set_blocking(stream, C.int(blocking))
}

func uv_udp_bind(udp *C.uv_udp_t, sa *C.struct_sockaddr, flags uint) C.int {
	return C.uv_udp_bind(udp, sa, C.uint(flags))
}

func uv_udp_send(req *C.uv_udp_send_t, udp *C.uv_udp_t, buf *C.uv_buf_t, bufcnt int, sa *C.struct_sockaddr) C.int {
	return C._uv_udp_send(req, udp, buf, C.uint(bufcnt), sa)
}

func uv_udp_recv_start(udp *C.uv_udp_t) C.int {
	return C._uv_udp_recv_start(udp)
}

func uv_udp_recv_stop(udp *C.uv_udp_t) C.int {
	return C.uv_udp_recv_stop(udp)
}

// BufInit (uv_buf_init) Constructor for uv_buf_t.
// Due to platform differences the user cannot rely on the ordering of the base and len members of the uv_buf_t struct. The user is responsible for freeing base after the uv_buf_t is done. Return struct passed by value.
func BufInit(b []byte) C.uv_buf_t {
	return C.uv_buf_init((*C.char)(unsafe.Pointer(&b[0])), C.uint(len(b)))
}

// BufInit2 (uv_buf_init) constructor for uv_buf_t from char*
func BufInit2(b *C.char, size C.uint) C.uv_buf_t {
	return C.uv_buf_init(b, C.uint(size))
}

// MallocUvBuf malloc in C memory of array of uv_buf_t
func MallocUvBuf(size C.uint) *C.uv_buf_t {
	return C.uv_buf_malloc(size)
}

// SetBuf set uv_buf_t inside bufs
func SetBuf(bufs *C.uv_buf_t, index uint, buf C.uv_buf_t) {
	C.uv_buf_set(bufs, C.uint(index), buf)
}

func uv_timer_start(timer *C.uv_timer_t, timeout uint64, repeat uint64) C.int {
	return C._uv_timer_start(timer, C.uint64_t(timeout), C.uint64_t(repeat))
}

func uv_poll_start(p *C.uv_poll_t, event int) C.int {
	return C._uv_poll_start(p, C.int(event))
}

func uv_signal_start(p *C.uv_signal_t, sigNum int) C.int {
	return C._uv_signal_start(p, C.int(sigNum))
}

func uv_signal_start_oneshot(p *C.uv_signal_t, sigNum int) C.int {
	return C._uv_signal_start_oneshot(p, C.int(sigNum))
}

func uv_prepare_start(p *C.uv_prepare_t) C.int {
	return C._uv_prepare_start(p)
}

func uv_check_start(p *C.uv_check_t) C.int {
	return C._uv_check_start(p)
}

func uv_idle_start(idle *C.uv_idle_t) C.int {
	return C._uv_idle_start(idle)
}

func uv_async_init(loop *C.uv_loop_t, async *C.uv_async_t) C.int {
	return C._uv_async_init(loop, async)
}

func uv_spawn(loop *C.uv_loop_t, process *C.uv_process_t, options *C.uv_process_options_t) C.int {
	return C._uv_spawn(loop, process, options)
}

// connect_socket initiate a connection on a socket.
// See also: https://linux.die.net/man/2/connect
func connect_socket(sock C.uv_os_sock_t, addr SockaddrIn) C.int {

	return C.connect_socket(sock, addr.GetSockAddr())
}

// close_socket return non zero if error occurs.
// Calling closesocket() on Window and close() on non-Window.
// See also: https://linux.die.net/man/2/close
func close_socket(sock C.uv_os_sock_t) C.int {
	return C.close_socket(sock)
}

// create_tcp_socket create socket with addr.
func create_tcp_socket(addr SockaddrIn, isBoundSocket int) C.uv_os_sock_t {
	return C.create_tcp_socket(addr.GetSockAddrIn(), C.int(isBoundSocket))
}

// test_sendAndRecv do test with send and recv over sock
func test_sendAndRecv(sock C.uv_os_sock_t) {
	defer func() {
		if e := recover(); e != nil {
		}
	}()
	C.test_sendAndRecv(sock)
}

func test_OpenFile(path string) C.int {
	p := C.CString(path)
	defer C.free(unsafe.Pointer(p))

	return C.test_Open(p)
}

/*
*			---------------------- EXPORT FUCTIONS ----------------------
 */

//export __uv_connect_cb
func __uv_connect_cb(c *C.uv_connect_t, status C.int) {
	if cbi := (*callbackInfo)(c.data); cbi.connect_cb != nil {
		cbi.connect_cb(&Request{
			(*C.uv_req_t)(unsafe.Pointer(c)),
			&Handle{(*C.uv_handle_t)(unsafe.Pointer(c)), cbi.data, cbi.ptr}}, int(status))
	}
}

//export __uv_connection_cb
func __uv_connection_cb(s *C.uv_stream_t, status C.int) {
	if cbi := (*callbackInfo)(s.data); cbi.connection_cb != nil {
		cbi.connection_cb(&Handle{(*C.uv_handle_t)(unsafe.Pointer(s)), cbi.data, cbi.ptr}, int(status))
	}
}

//export __uv_read_cb
func __uv_read_cb(s *C.uv_stream_t, nread C.ssize_t, buf *C.uv_buf_t) {
	if cbi := (*callbackInfo)(s.data); cbi.read_cb != nil {
		if nread == -1 || nread == C.UV_EOF {
			C.free(unsafe.Pointer(buf.base))
			return
		}

		cbi.read_cb(&Handle{(*C.uv_handle_t)(unsafe.Pointer(s)), cbi.data, cbi.ptr}, buf, nread)
	}
}

//export __uv_write_cb
func __uv_write_cb(w *C.uv_write_t, status C.int) {
	if cbi := (*callbackInfo)(w.data); cbi.write_cb != nil {
		cbi.write_cb(&Request{
			(*C.uv_req_t)(unsafe.Pointer(w)),
			&Handle{(*C.uv_handle_t)(unsafe.Pointer(w)), cbi.data, cbi.ptr}}, int(status))
	}
}

//export __uv_close_cb
func __uv_close_cb(h *C.uv_handle_t) {
	if cbi := (*callbackInfo)(h.data); cbi.close_cb != nil {
		cbi.close_cb(&Handle{h, cbi.data, cbi.ptr})
	}
}

//export __uv_prepare_cb
func __uv_prepare_cb(h *C.uv_prepare_t) {
	if cbi := (*callbackInfo)(h.data); cbi.prepare_cb != nil {
		cbi.prepare_cb(&Handle{(*C.uv_handle_t)(unsafe.Pointer(h)), cbi.data, cbi.ptr})
	}
}

//export __uv_check_cb
func __uv_check_cb(h *C.uv_check_t) {
	if cbi := (*callbackInfo)(h.data); cbi.check_cb != nil {
		cbi.check_cb(&Handle{(*C.uv_handle_t)(unsafe.Pointer(h)), cbi.data, cbi.ptr})
	}
}

//export __uv_async_cb
func __uv_async_cb(h *C.uv_prepare_t) {
	if cbi := (*callbackInfo)(h.data); cbi.async_cb != nil {
		cbi.async_cb(&Handle{(*C.uv_handle_t)(unsafe.Pointer(h)), cbi.data, cbi.ptr})
	}
}

//export __uv_shutdown_cb
func __uv_shutdown_cb(s *C.uv_shutdown_t, status C.int) {
	if cbi := (*callbackInfo)(s.data); cbi.shutdown_cb != nil {
		cbi.shutdown_cb(&Request{
			(*C.uv_req_t)(unsafe.Pointer(s)),
			&Handle{(*C.uv_handle_t)(unsafe.Pointer(s.handle)), cbi.data, cbi.ptr}}, int(status))
	}
}

//export __uv_udp_recv_cb
func __uv_udp_recv_cb(u *C.uv_udp_t, nread C.ssize_t, buf *C.uv_buf_t, sa *C.struct_sockaddr, flags C.uint) {
	if cbi := (*callbackInfo)(u.data); cbi.udp_recv_cb != nil {
		nRead := int(nread)
		if nRead < 0 {
			cbi.udp_recv_cb(&Handle{
				(*C.uv_handle_t)(unsafe.Pointer(u)), cbi.data, cbi.ptr}, nil, sa, uint(flags))
		} else {
			cbi.udp_recv_cb(&Handle{
				(*C.uv_handle_t)(unsafe.Pointer(u)), cbi.data, cbi.ptr}, (*[1 << 30]byte)(unsafe.Pointer(buf.base))[0:nRead], sa, uint(flags))
		}
	}
}

//export __uv_udp_send_cb
func __uv_udp_send_cb(us *C.uv_udp_send_t, status C.int) {
	if cbi := (*callbackInfo)(us.handle.data); cbi.udp_send_cb != nil {
		cbi.udp_send_cb(&Request{
			(*C.uv_req_t)(unsafe.Pointer(us)),
			&Handle{(*C.uv_handle_t)(unsafe.Pointer(us.handle)), cbi.data, cbi.ptr}}, int(status))
	}
}

//export __uv_timer_cb
func __uv_timer_cb(t *C.uv_timer_t, status C.int) {
	if cbi := (*callbackInfo)(t.data); cbi.timer_cb != nil {
		cbi.timer_cb(&Handle{(*C.uv_handle_t)(unsafe.Pointer(t)), cbi.data, cbi.ptr}, int(status))
	}
}

//export __uv_poll_cb
func __uv_poll_cb(p *C.uv_poll_t, status, events C.int) {
	if cbi := (*callbackInfo)(p.data); cbi.timer_cb != nil {
		cbi.poll_cb(&Handle{(*C.uv_handle_t)(unsafe.Pointer(p)), cbi.data, cbi.ptr}, int(status), int(events))
	}
}

//export __uv_signal_cb
func __uv_signal_cb(s *C.uv_signal_t, sigNum C.int) {
	if cbi := (*callbackInfo)(s.data); cbi.signal_cb != nil {
		cbi.signal_cb(&Handle{(*C.uv_handle_t)(unsafe.Pointer(s)), cbi.data, cbi.ptr}, int(sigNum))
	}
}

//export __uv_idle_cb
func __uv_idle_cb(i *C.uv_idle_t, status C.int) {
	if cbi := (*callbackInfo)(i.data); cbi.idle_cb != nil {
		cbi.idle_cb(&Handle{(*C.uv_handle_t)(unsafe.Pointer(i)), cbi.data, cbi.ptr}, int(status))
	}
}

//export __uv_exit_cb
func __uv_exit_cb(pc *C.uv_process_t, exit_status C.int, term_signal C.int) {
	if cbi := (*callbackInfo)(pc.data); cbi.exit_cb != nil {
		cbi.exit_cb(&Handle{(*C.uv_handle_t)(unsafe.Pointer(pc)), cbi.data, cbi.ptr}, int(exit_status), int(term_signal))
	}
}
