package gouv

// #cgo pkg-config: libuv
/*
#include <uv.h>
#include "common.h"
#include <stdlib.h>
uv_tty_t* mallocTTY() {
	return (uv_tty_t*)malloc(sizeof(uv_tty_t));
}

window_size_t* get_tty_window_size(uv_tty_t* tty) {
	window_size_t* ws;
	ws = (window_size_t*)malloc(sizeof(window_size_t));
	ws->err = uv_tty_get_winsize(tty, &ws->width, &ws->height);

	return ws;
}
*/
import "C"
import "unsafe"

// UvTTY handles represent a stream for the console.
type UvTTY struct {
	t *C.uv_tty_t
	l *C.uv_loop_t
	UvStream
}

// UvTTYInit (uv_tty_init) initialize a new TTY stream with the given file descriptor. Usually the file descriptor will be:
// 0 = stdin
// 1 = stdout
// 2 = stderr
// readable, specifies if you plan on calling uv_read_start() with this stream. stdin is readable, stdout is not.
// On Unix this function will determine the path of the fd of the terminal using ttyname_r(3), open it, and use it if the passed file descriptor refers to a TTY. This lets libuv put the tty in non-blocking mode without affecting other processes that share the tty.
// This function is not thread safe on systems that don’t support ioctl TIOCGPTN or TIOCPTYGNAME, for instance OpenBSD and Solaris.
func UvTTYInit(loop *UvLoop, fd C.uv_file, readable int, data interface{}) (*UvTTY, error) {
	t := C.mallocTTY()

	if loop == nil {
		loop = UvLoopDefault()
	}

	if r := C.uv_tty_init(loop.GetNativeLoop(), t, fd, C.int(readable)); r != 0 {
		return nil, ParseUvErr(r)
	}

	res := &UvTTY{}
	t.data = unsafe.Pointer(&callbackInfo{ptr: res, data: data})
	res.s, res.l, res.t, res.Handle = (*C.uv_stream_t)(unsafe.Pointer(t)), loop.GetNativeLoop(), t, Handle{(*C.uv_handle_t)(unsafe.Pointer(t)), t.data, res}

	return res, nil
}

// SetMode (uv_tty_set_mode) set the TTY using the specified terminal mode.
func (t *UvTTY) SetMode(mode UV_TTY_MODE) C.int {
	return C.uv_tty_set_mode(t.t, C.uv_tty_mode_t(mode))
}

// ResetMode (uv_tty_reset_mode) to be called when the program exits. Resets TTY settings to default values for the next process to take over.
// This function is async signal-safe on Unix platforms but can fail with error code UV_EBUSY if you call it when execution is inside uv_tty_set_mode().
func (t *UvTTY) ResetMode() C.int {
	return C.uv_tty_reset_mode()
}

// GetWinsize (uv_tty_get_winsize) gets the current Window size. On success it returns 0.
func (t *UvTTY) GetWinsize() (err C.int, width C.int, height C.int) {
	ws := C.get_tty_window_size(t.t)
	defer C.free(unsafe.Pointer(ws))

	return ws.err, ws.width, ws.height
}
