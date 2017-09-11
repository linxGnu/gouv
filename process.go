package gouv

// #cgo pkg-config: libuv
/*
#include <uv.h>
#include <stdlib.h>
uv_process_options_t* mallocProcessOptsT() {
	return (uv_process_options_t*)malloc(sizeof(uv_process_options_t));
}
uv_process_t* mallocProcessT() {
	return (uv_process_t*)malloc(sizeof(uv_process_t));
}
uv_stdio_container_t* mallocStdioContainerT() {
	return (uv_stdio_container_t*)malloc(sizeof(uv_stdio_container_t));
}
void set_data_in_StdioContainer(uv_stdio_container_t* container, uv_stream_t* stream, int fd) {
	container->data.stream = stream;
	container->data.fd = fd;
}
*/
import "C"
import "unsafe"

// UvStdioContainerData union of stream and fd
type UvStdioContainerData struct {
	Stream *C.uv_stream_t
	Fd     int
}

// Freemem of UvStdioContainerData
func (c *UvStdioContainerData) Freemem() {
	C.free(unsafe.Pointer(c.Stream))
}

// UvStdioContainer container for each stdio handle or fd passed to a child process.
type UvStdioContainer struct {
	Flags C.uv_stdio_flags
	Data  *UvStdioContainerData
}

// Freemem of UvStdioContainer
func (c *UvStdioContainer) Freemem() {
	if c.Data != nil {
		c.Data.Freemem()
	}
}

// UvProcessOptions options for spawning the process, passed to uv_spawn()
type UvProcessOptions struct {
	ExitCb     func(*Handle, int, int)
	File       string
	Args       []string
	Env        []string
	Cwd        string
	Flags      uint
	StdioCount int
	Stdio      *UvStdioContainer
	UID        uint8
	GID        uint8
}

// Freemem freemem of process handle
func (p *UvProcessOptions) Freemem() {
	if p.Stdio != nil {
		p.Stdio.Freemem()
	}
}

// UvProcess process handles will spawn a new process and allow the user to control it and establish communication channels with it using streams.
type UvProcess struct {
	p *C.uv_process_t
	l *C.uv_loop_t
	Handle
}

// UvDisableStdioInheritance (uv_disable_stdio_inheritance) disables inheritance for file descriptors / handles that this process inherited from its parent. The effect is that child processes spawned by this process don’t accidentally inherit these handles.
// It is recommended to call this function as early in your program as possible, before the inherited file descriptors can be closed or duplicated.
func UvDisableStdioInheritance() {
	C.uv_disable_stdio_inheritance()
}

// UvSpawnProcess initializes the process handle and starts the process. If the process is successfully spawned, this function will return 0. Otherwise, the negative error code corresponding to the reason it couldn’t spawn is returned.
// Possible reasons for failing to spawn would include (but not be limited to) the file to execute not existing, not having permissions to use the setuid or setgid specified, or not having enough memory to allocate for the new process.
func UvSpawnProcess(loop *UvLoop, options *UvProcessOptions, data interface{}) (*UvProcess, error) {
	if loop == nil {
		loop = UvLoopDefault()
	}

	// malloc in c memory space
	opt := C.mallocProcessOptsT()
	defer C.free(unsafe.Pointer(opt))

	if len(options.File) > 0 {
		opt.file = C.CString(options.File)
		defer C.free(unsafe.Pointer(opt.file))
	}

	if len(options.Args) > 0 {
		opt.args = (**C.char)(C.malloc(C.size_t(4 * (len(options.Args) + 1))))
		defer C.free(unsafe.Pointer(opt.args))

		for n := 0; n < len(options.Args); n++ {
			((*[1 << 24]*C.char)(unsafe.Pointer(&opt.args)))[n] = C.CString(options.Args[n])
		}
		((*[1 << 24]*C.char)(unsafe.Pointer(&opt.args)))[len(options.Args)] = nil
	}

	if len(options.Env) > 0 {
		opt.env = (**C.char)(C.malloc(C.size_t(4 * (len(options.Env) + 1))))
		defer C.free(unsafe.Pointer(opt.env))

		for n := 0; n < len(options.Args); n++ {
			((*[1 << 24]*C.char)(unsafe.Pointer(&opt.env)))[n] = C.CString(options.Args[n])
		}
		((*[1 << 24]*C.char)(unsafe.Pointer(&opt.env)))[len(options.Args)] = nil
	}

	if len(options.Cwd) > 0 {
		opt.cwd = C.CString(options.Cwd)
		defer C.free(unsafe.Pointer(opt.cwd))
	}

	opt.flags = C.uint(options.Flags)

	opt.stdio_count = C.int(options.StdioCount)
	if options.Stdio != nil {
		opt.stdio = C.mallocStdioContainerT()
		defer C.free(unsafe.Pointer(opt.stdio))

		opt.stdio.flags = C.uv_stdio_flags(options.Stdio.Flags)
		if options.Stdio.Data != nil {
			C.set_data_in_StdioContainer(opt.stdio, options.Stdio.Data.Stream, C.int(options.Stdio.Data.Fd))
		}
	}

	opt.uid = C.uv_uid_t(options.UID)
	opt.gid = C.uv_gid_t(options.GID)

	p := C.mallocProcessT()
	p.data = unsafe.Pointer(&callback_info{exit_cb: options.ExitCb, data: data})

	if r := uv_spawn(loop.GetNativeLoop(), p, opt); r != 0 {
		return nil, ParseUvErr(r)
	}

	return &UvProcess{p, loop.GetNativeLoop(), Handle{(*C.uv_handle_t)(unsafe.Pointer(p)), p.data}}, nil
}

// Kill (uv_process_kill) sends the specified signal to the given process handle. Check the documentation on uv_signal_t — Signal handle for signal support, specially on Windows.
func (p *UvProcess) Kill(sigNum C.int) error {
	if r := C.uv_process_kill(p.p, sigNum); r != 0 {
		return ParseUvErr(r)
	}

	return nil
}

// Freemem freemem of process handle
func (p *UvProcess) Freemem() {
	C.free(unsafe.Pointer(p.p))
}

// UvKill (uv_kill) sends the specified signal to the given PID. Check the documentation on uv_signal_t — Signal handle for signal support, specially on Windows.
func UvKill(pid int, sigNum C.int) error {
	if r := C.uv_kill(C.int(pid), sigNum); r != 0 {
		return ParseUvErr(r)
	}

	return nil
}
