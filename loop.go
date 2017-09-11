package binduv

// #cgo pkg-config: libuv
/*
#include <uv.h>
#include <stdlib.h>
*/
import "C"
import "fmt"

// UvRunMode wrapper of uv_run_mode
type UvRunMode int

const (
	// UVRUNDEFAULT runs the event loop until there are no more active and referenced handles or requests. Returns non-zero if uv_stop() was called and there are still active handles or requests. Returns zero in all other cases.
	UVRUNDEFAULT UvRunMode = 1

	// UVRUNONCE poll for i/o once. Note that this function blocks if there are no pending callbacks. Returns zero when done (no active handles or requests left), or non-zero if more callbacks are expected (meaning you should run the event loop again sometime in the future).
	UVRUNONCE UvRunMode = 2

	// UVRUNNOWAIT poll for i/o once but don’t block if there are no pending callbacks. Returns zero if done (no active handles or requests left), or non-zero if more callbacks are expected (meaning you should run the event loop again sometime in the future).
	UVRUNNOWAIT UvRunMode = 3
)

// UvLoop wrapper of C.uv_loop_t
type UvLoop struct {
	loop *C.uv_loop_t
}

// GetNativeLoop get native c pointer of loop
func (uvl *UvLoop) GetNativeLoop() *C.uv_loop_t {
	return uvl.loop
}

// UvLoopDefault (uv_default_loop) return uv_loop default
func UvLoopDefault() *UvLoop {
	return &UvLoop{loop: C.uv_default_loop()}
}

// Init (uv_loop_init) initialize uv_loop
func (uvl *UvLoop) Init() error {
	if r := C.uv_loop_init(uvl.loop); r != 0 {
		return ParseUvErr(r)
	}

	return nil
}

// Close (uv_loop_close) releases all internal loop resources.
// Call this function only when the loop has finished executing and all open handles and requests have been closed, or it will return UV_EBUSY.
// After this function returns, the user can free the memory allocated for the loop.
func (uvl *UvLoop) Close() error {
	if r := C.uv_loop_close(uvl.loop); r != 0 {
		return ParseUvErr(r)
	}

	return nil
}

// Alive (uv_loop_alive) returns non-zero if there are active handles or request in the loop.
func (uvl *UvLoop) Alive() int {
	return int(C.uv_loop_alive(uvl.loop))
}

// Fork (uv_loop_fork) reinitialize any kernel state necessary in the child process after a fork(2) system call.
func (uvl *UvLoop) Fork() error {
	if r := C.uv_loop_fork(uvl.loop); r != 0 {
		return ParseUvErr(r)
	}

	return nil
}

// Run (uv_run) This function runs the event loop. It will act differently depending on the specified mode.
func (uvl *UvLoop) Run(mode UvRunMode) error {
	switch mode {
	case UVRUNDEFAULT:
		if r := C.uv_run(uvl.loop, C.UV_RUN_DEFAULT); r != 0 {
			return ParseUvErr(r)
		}
		return nil
	case UVRUNONCE:
		if r := C.uv_run(uvl.loop, C.UV_RUN_ONCE); r != 0 {
			return ParseUvErr(r)
		}
		return nil
	case UVRUNNOWAIT:
		if r := C.uv_run(uvl.loop, C.UV_RUN_NOWAIT); r != 0 {
			return ParseUvErr(r)
		}
		return nil
	}

	return fmt.Errorf("Unknown run mode")
}

// Stop (uv_stop) Stop the event loop, causing uv_run() to end as soon as possible.
// This will happen not sooner than the next loop iteration.
// If this function was called before blocking for i/o, the loop won’t block for i/o on this iteration.
func (uvl *UvLoop) Stop() {
	C.uv_stop(uvl.loop)
}

// UpdateTime (uv_update_time) Update the event loop’s concept of “now”. Libuv caches the current time at the start of the event loop tick in order to reduce the number of time-related system calls.
// You won’t normally need to call this function unless you have callbacks that block the event loop for longer periods of time, where “longer” is somewhat subjective but probably on the order of a millisecond or more.
func (uvl *UvLoop) UpdateTime() {
	C.uv_update_time(uvl.loop)
}

// Now (uv_now) return the current timestamp in milliseconds. The timestamp is cached at the start of the event loop tick, see uv_update_time() for details and rationale.
func (uvl *UvLoop) Now() uint64 {
	return uint64(C.uv_now(uvl.loop))
}

// BackendFD (uv_backend_fd) get backend file descriptor. Only kqueue, epoll and event ports are supported.
func (uvl *UvLoop) BackendFD() int {
	return int(C.uv_backend_fd(uvl.loop))
}

// BackendTimeout (uv_backend_timeout) Get the poll timeout. The return value is in milliseconds, or -1 for no timeout.
func (uvl *UvLoop) BackendTimeout() uint64 {
	return uint64(C.uv_backend_timeout(uvl.loop))
}

// Configure (uv_loop_configure) TODO in the future. Set additional loop options
// func (uvl *UvLoop) Configure(option C.uv_loop_option, p ...interface{}) error {
// 	return nil
// }
