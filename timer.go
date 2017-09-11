package gouv

// #cgo pkg-config: libuv
/*
#include <uv.h>
#include <stdlib.h>
uv_timer_t* mallocTimeT() {
	return (uv_timer_t*)malloc(sizeof(uv_timer_t));
}
*/
import "C"
import "unsafe"

// UvTimer timer handles are used to schedule callbacks to be called in the future.
type UvTimer struct {
	t *C.uv_timer_t
	l *C.uv_loop_t
	Handle
}

// TimerInit initialize the timer handle
func TimerInit(loop *UvLoop, data interface{}) (timer *UvTimer, err error) {
	t := C.mallocTimeT()

	if loop == nil {
		loop = UvLoopDefault()
	}

	if r := C.uv_timer_init(loop.GetNativeLoop(), t); r != 0 {
		return nil, ParseUvErr(r)
	}

	t.data = unsafe.Pointer(&callback_info{data: data})
	return &UvTimer{t, loop.GetNativeLoop(), Handle{(*C.uv_handle_t)(unsafe.Pointer(t)), t.data}}, nil
}

// GetLoop get loop if this handle
func (t *UvTimer) GetLoop() *UvLoop {
	return &UvLoop{t.l}
}

// Start (uv_timer_start) start the timer. timeout and repeat are in milliseconds.
func (t *UvTimer) Start(cb func(*Handle, int), timeout uint64, repeat uint64) (err error) {
	cbi := (*callback_info)(t.t.data)
	cbi.timer_cb = cb

	if r := uv_timer_start(t.t, timeout, repeat); r != 0 {
		err = ParseUvErr(r)
		return
	}

	return nil
}

// Stop (uv_timer_stop) the timer, the callback will not be called anymore.
func (t *UvTimer) Stop() (err error) {
	if r := C.uv_timer_stop(t.t); r != 0 {
		err = ParseUvErr(r)
		return
	}

	return nil
}

// Again (uv_timer_again) stop the timer, and if it is repeating restart it using the repeat value as the timeout. If the timer has never been started before it returns UV_EINVAL
func (t *UvTimer) Again() (err error) {
	if r := C.uv_timer_again(t.t); r != 0 {
		err = ParseUvErr(r)
		return
	}

	return nil
}

// SetRepeat (uv_timer_set_repeat) set the repeat interval value in milliseconds. The timer will be scheduled to run on the given interval, regardless of the callback execution duration, and will follow normal timer semantics in the case of a time-slice overrun.
func (t *UvTimer) SetRepeat(repeat uint64) {
	C.uv_timer_set_repeat(t.t, C.uint64_t(repeat))
}

// GetRepeat (uv_timer_get_repeat) get the timer repeat value.
func (t *UvTimer) GetRepeat() uint64 {
	return uint64(C.uv_timer_get_repeat(t.t))
}

// Freemem freemem of timer
func (t *UvTimer) Freemem() {
	C.free(unsafe.Pointer(t.t))
}
