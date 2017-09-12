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

uv_stdio_container_t* mallocStdioContainerArrT(int size) {
	return malloc(size * sizeof(uv_stdio_container_t));
}

char** mallocProcessOptsArgs(int size) {
	char** result = malloc((size + 1) * sizeof(char*));
	result[size] = NULL;

	return result;
}

void set_process_args(char** args, int i, char* st) {
	args[i] = st;
}

void set_data_in_StdioContainer(uv_stdio_container_t* container, int i, uv_stdio_flags flags, uv_stream_t* stream, int fd) {
	container[i].flags = flags;
	container[i].data.stream = stream;
	container[i].data.fd = fd;
}
*/
import "C"
import (
	"unsafe"
)

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
	// ExitCb callback called after the process exits.
	ExitCb func(*Handle, int, int)

	// File path pointing to the program to be executed.
	File string

	// Args command line arguments. args[0] should be the path to the program. On Windows this uses CreateProcess which concatenates the arguments into a string this can cause some strange errors. See the UV_PROCESS_WINDOWS_VERBATIM_ARGUMENTS flag on uv_process_flags.
	Args []string

	// Env environment for the new process. If NULL the parents environment is used.
	Env []string

	// Cwd current working directory for the subprocess.
	Cwd string

	// Flags various flags that control how uv_spawn() behaves. See uv_process_flags.
	Flags UV_PROCESS_FLAGS

	// Stdio the stdio field points to an array of uv_stdio_container_t structs that describe the file descriptors that will be made available to the child process. The convention is that stdio[0] points to stdin, fd 1 is used for stdout, and fd 2 is stderr.
	Stdio []*UvStdioContainer

	// Uid libuv can change the child process’ user id. This happens only when the appropriate bits are set in the flags fields.
	UID uint8

	// GID libuv can change the child process’ group id. This happens only when the appropriate bits are set in the flags fields.
	GID uint8
}

// Freemem freemem of process handle
func (p *UvProcessOptions) Freemem() {
	if p.Stdio != nil {
		for _, v := range p.Stdio {
			v.Freemem()
		}
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
		opt.args = C.mallocProcessOptsArgs(C.int(len(options.Args)))
		for i := range options.Args {
			tmp := C.CString(options.Args[i])
			C.set_process_args(opt.args, C.int(i), tmp)
			defer C.free(unsafe.Pointer(tmp))
		}
		defer C.free(unsafe.Pointer(opt.args))
	}

	if len(options.Env) > 0 {
		opt.env = C.mallocProcessOptsArgs(C.int(len(options.Env)))
		for i := range options.Env {
			tmp := C.CString(options.Env[i])
			C.set_process_args(opt.env, C.int(i), tmp)
			defer C.free(unsafe.Pointer(tmp))
		}
		defer C.free(unsafe.Pointer(opt.env))
	}

	if len(options.Cwd) > 0 {
		opt.cwd = C.CString(options.Cwd)
		defer C.free(unsafe.Pointer(opt.cwd))
	}

	opt.flags = C.uint(options.Flags)

	if options.Stdio == nil {
		opt.stdio_count = 0
	} else {
		opt.stdio_count = C.int(len(options.Stdio))
		opt.stdio = C.mallocStdioContainerArrT(opt.stdio_count)
		for i, v := range options.Stdio {
			C.set_data_in_StdioContainer(opt.stdio, C.int(i), v.Flags, v.Data.Stream, C.int(v.Data.Fd))
		}
		defer C.free(unsafe.Pointer(opt.stdio))
	}

	opt.uid = C.uv_uid_t(options.UID)
	opt.gid = C.uv_gid_t(options.GID)

	p := C.mallocProcessT()
	p.data = unsafe.Pointer(&callbackInfo{exit_cb: options.ExitCb, data: data})

	if r := uv_spawn(loop.GetNativeLoop(), p, opt); r != 0 {
		return nil, ParseUvErr(r)
	}

	return &UvProcess{p, loop.GetNativeLoop(), Handle{(*C.uv_handle_t)(unsafe.Pointer(p)), p.data}}, nil
}

// Kill (uv_process_kill) sends the specified signal to the given process handle. Check the documentation on uv_signal_t — Signal handle for signal support, specially on Windows.
func (p *UvProcess) Kill(sigNum C.int) (err error) {
	if r := C.uv_process_kill(p.p, sigNum); r != 0 {
		err = ParseUvErr(r)
		return
	}

	return
}

// Freemem freemem of process handle
func (p *UvProcess) Freemem() {
	C.free(unsafe.Pointer(p.p))
}

// Unref unrefernce this process
func (p *UvProcess) Unref() {
	uv_unref((*C.uv_handle_t)(unsafe.Pointer(p.p)))
}

// UvKill (uv_kill) sends the specified signal to the given PID. Check the documentation on uv_signal_t — Signal handle for signal support, specially on Windows.
func UvKill(pid int, sigNum C.int) error {
	if r := C.uv_kill(C.int(pid), sigNum); r != 0 {
		return ParseUvErr(r)
	}

	return nil
}
