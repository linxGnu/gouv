package gouv

// #cgo pkg-config: libuv
/*
#include <uv.h>
#include <stdlib.h>
uv_pipe_t* mallocPipeT() {
	return (uv_pipe_t*)malloc(sizeof(uv_pipe_t));
}
*/
import "C"
import "unsafe"

// UvPipe handles provide an abstraction over local domain sockets on Unix and named pipes on Windows.
type UvPipe struct {
	p *C.uv_pipe_t
	l *C.uv_loop_t
	UvStream
}

// UvPipeInit (uv_pipe_init) initialize a pipe handle. The ipc argument is a boolean to indicate if this pipe will be used for handle passing between processes.
func UvPipeInit(loop *UvLoop, ipc int, data interface{}) (*UvPipe, error) {
	t := C.mallocPipeT()

	if loop == nil {
		loop = UvLoopDefault()
	}

	if r := C.uv_pipe_init(loop.GetNativeLoop(), t, C.int(ipc)); r != 0 {
		C.free(unsafe.Pointer(t))
		return nil, ParseUvErr(r)
	}

	res := &UvPipe{}
	t.data = unsafe.Pointer(&callbackInfo{ptr: res, data: data})
	res.s, res.l, res.p, res.Handle = (*C.uv_stream_t)(unsafe.Pointer(t)), loop.GetNativeLoop(), t, Handle{(*C.uv_handle_t)(unsafe.Pointer(t)), t.data, res}

	return res, nil
}

// Open (uv_pipe_open) open an existing file descriptor or HANDLE as a pipe.
// Note: the passed file descriptor or HANDLE is not checked for its type, but it’s required that it represents a valid pipe.
func (p *UvPipe) Open(file C.uv_file) C.int {
	return C.uv_pipe_open(p.p, file)
}

// Bind (uv_pipe_bind) bind the pipe to a file path (Unix) or a name (Windows).
// Note: paths on Unix get truncated to sizeof(sockaddr_un.sun_path) bytes, typically between 92 and 108 bytes.
func (p *UvPipe) Bind(name string) C.int {
	tmp := C.CString(name)
	defer C.free(unsafe.Pointer(tmp))

	return C.uv_pipe_bind(p.p, tmp)
}

// Connect (uv_pipe_connect) connect to the Unix domain socket or the named pipe.
// Note: paths on Unix get truncated to sizeof(sockaddr_un.sun_path) bytes, typically between 92 and 108 bytes.
func (p *UvPipe) Connect(req *UvConnect, name string, cb func(*Request, int)) {
	cbi := (*callbackInfo)(req.c.data)
	cbi.connect_cb = cb
	cbi.ptr = p

	uv_pipe_connect(req.c, p.p, name)
}

// PendingInstances (uv_pipe_pending_instances) set the number of pending pipe instance handles when the pipe server is waiting for connections.
// Note: this setting applies to Windows only.
func (p *UvPipe) PendingInstances(count int) {
	C.uv_pipe_pending_instances(p.p, C.int(count))
}

// PendingCount (uv_pipe_pending_count) return number of pending instances.
func (p *UvPipe) PendingCount() C.int {
	return C.uv_pipe_pending_count(p.p)
}

// PendingType (uv_pipe_pending_type) used to receive handles over IPC pipesu
// First - call uv_pipe_pending_count(), if it’s > 0 then initialize a handle of the given type, returned by uv_pipe_pending_type() and call uv_accept(pipe, handle).
func (p *UvPipe) PendingType() C.uv_handle_type {
	return C.uv_pipe_pending_type(p.p)
}

// GetPipeHandle get handle
func (p *UvPipe) GetPipeHandle() *C.uv_pipe_t {
	return p.p
}

// Freemem freemem pipe
func (p *UvPipe) Freemem() {
	C.free(unsafe.Pointer(p.p))
}

// TODO: uv_pipe_getsockname
// TODO: uv_pipe_getpeername
