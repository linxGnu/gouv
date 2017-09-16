package gouv

// #cgo pkg-config: libuv
/*
#include <uv.h>
#include "common.h"
#include <stdlib.h>
uv_fs_event_t* mallocFSEventT() {
	return (uv_fs_event_t*)malloc(sizeof(uv_fs_event_t));
}

char_result_t* _uv_fs_event_getpath(uv_fs_event_t* handle) {
	char_result_t* c;
	int r;

	c = (char_result_t*)malloc(sizeof(char_result_t));
	c->c = (char*)malloc(257 * sizeof(char));

	r = uv_fs_event_getpath(handle, c->c, &(c->size));
	if (r == UV_ENOBUFS) {
		free(c->c); // free mem first

		// make new one
		c->c = (char*)malloc((c->size + 1) * sizeof(char));
		c->c[c->size] = '\0';

		// get again
		r = uv_fs_event_getpath(handle, c->c, &c->size);
	}
	c->err = r;

	return c;
}
*/
import "C"
import (
	"unsafe"
)

// UvFSEvent (uv_fs_event_t) FS Event handles allow the user to monitor a given path for changes, for example, if the file was renamed or there was a generic change in it.
// This handle uses the best backend for the job on each platform.
type UvFSEvent struct {
	FSEvent *C.uv_fs_event_t
	Loop    *C.uv_loop_t
	Handle
}

// UvFSEventInit (uv_fs_event_init) initialize the handle.
func UvFSEventInit(loop *UvLoop, data interface{}) (*UvFSEvent, error) {
	t := C.mallocFSEventT()

	if loop == nil {
		loop = UvLoopDefault()
	}

	res := &UvFSEvent{}
	t.data = unsafe.Pointer(&callbackInfo{data: data, ptr: res})
	res.FSEvent, res.Loop, res.Handle = t, loop.GetNativeLoop(), Handle{(*C.uv_handle_t)(unsafe.Pointer(t)), t.data, res}
	if r := C.uv_fs_event_init(loop.GetNativeLoop(), t); r != 0 {
		C.free(unsafe.Pointer(t))
		return nil, ParseUvErr(r)
	}

	return res, nil
}

// Start (uv_fs_event_start) start the handle with the given callback, which will watch the specified path for changes. flags can be an ORed mask of uv_fs_event_flags.
func (f *UvFSEvent) Start(cb func(*Handle, *C.char, int, int), path string, flags uint) C.int {
	cbi := (*callbackInfo)(f.FSEvent.data)
	cbi.fs_event_cb = cb

	_path := C.CString(path)
	defer C.free(unsafe.Pointer(_path))

	return uv_fs_event_start(f.FSEvent, _path, flags)
}

// Stop (uv_fs_event_stop) stop the handle, the callback will no longer be called.
func (f *UvFSEvent) Stop() C.int {
	return C.uv_fs_event_stop(f.FSEvent)
}

// GetPath (uv_fs_event_getpath) get the path being monitored by the handle.
// Returns 0 on success or an error code < 0 in case of failure.
// On success, buffer will contain the path and size its length.
// If the buffer is not big enough UV_ENOBUFS will be returned and size will be set to the required size, including the null terminator.
func (f *UvFSEvent) GetPath() (C.int, *C.char, C.size_t) {
	path := C._uv_fs_event_getpath(f.FSEvent)
	defer C.free(unsafe.Pointer(path))

	return path.err, path.c, path.size
}

// Freemem freemem handle
func (f *UvFSEvent) Freemem() {
	C.free(unsafe.Pointer(f.FSEvent))
}
