package gouv

// #cgo pkg-config: libuv
/*
#include <uv.h>
#include "common.h"
#include <stdlib.h>
uv_fs_poll_t* mallocFSPollT() {
	return (uv_fs_poll_t*)malloc(sizeof(uv_fs_poll_t));
}

char_result_t* _uv_fs_poll_getpath(uv_fs_poll_t* handle) {
	char_result_t* c;
	int r;

	c = (char_result_t*)malloc(sizeof(char_result_t));
	c->c = (char*)malloc(257 * sizeof(char));

	r = uv_fs_poll_getpath(handle, c->c, &(c->size));
	if (r == UV_ENOBUFS) {
		free(c->c); // free mem first

		// make new one
		c->c = (char*)malloc((c->size + 1) * sizeof(char));
		c->c[c->size] = '\0';

		// get again
		r = uv_fs_poll_getpath(handle, c->c, &c->size);
	}
	c->err = r;

	return c;
}
*/
import "C"
import "C"
import "unsafe"

// UvFSPoll handles allow the user to monitor a given path for changes. Unlike uv_fs_event_t, fs poll handles use stat to detect when a file has changed so they can work on file systems where fs event handles can’t.
type UvFSPoll struct {
	f *C.uv_fs_poll_t
	l *C.uv_loop_t
	Handle
}

// UvFSPollInit (uv_fs_poll_init) initialize the handle.
func UvFSPollInit(loop *UvLoop, data interface{}) (*UvFSPoll, error) {
	t := C.mallocFSPollT()

	if loop == nil {
		loop = UvLoopDefault()
	}

	res := &UvFSPoll{}
	t.data = unsafe.Pointer(&callbackInfo{data: data, ptr: res})
	res.f, res.l, res.Handle = t, loop.GetNativeLoop(), Handle{(*C.uv_handle_t)(unsafe.Pointer(t)), t.data, res}
	if r := C.uv_fs_poll_init(loop.GetNativeLoop(), t); r != 0 {
		C.free(unsafe.Pointer(t))
		return nil, ParseUvErr(r)
	}

	return res, nil
}

// Start (uv_fs_event_start) check the file at path for changes every interval milliseconds.
func (f *UvFSPoll) Start(cb func(h *Handle, status C.int, prev *C.uv_stat_t, current *C.uv_stat_t), path string, interval uint) C.int {
	cbi := (*callbackInfo)(f.f.data)
	cbi.fs_poll_cb = cb

	_path := C.CString(path)
	defer C.free(unsafe.Pointer(_path))

	return uv_fs_poll_start(f.f, _path, interval)
}

// Stop (uv_fs_event_stop) stop the handle, the callback will no longer be called.
func (f *UvFSPoll) Stop() C.int {
	return C.uv_fs_poll_stop(f.f)
}

// GetPath (uv_fs_poll_getpath) get the path being monitored by the handle.
// Returns 0 on success or an error code < 0 in case of failure.
// On success, buffer will contain the path and size its length.
// If the buffer is not big enough UV_ENOBUFS will be returned and size will be set to the required size, including the null terminator.
func (f *UvFSPoll) GetPath() (C.int, *C.char, C.size_t) {
	path := C._uv_fs_poll_getpath(f.f)
	defer C.free(unsafe.Pointer(path))

	return path.err, path.c, path.size
}

// Freemem freemem handle
func (f *UvFSPoll) Freemem() {
	C.free(unsafe.Pointer(f.f))
}
