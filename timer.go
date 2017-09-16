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
	Timer *C.uv_timer_t
	Loop  *C.uv_loop_t
	Handle
}

// TimerInit initialize the timer handle
func TimerInit(loop *UvLoop, data interface{}) (*UvTimer, error) {
	t := C.mallocTimeT()

	if loop == nil {
		loop = UvLoopDefault()
	}

	if r := C.uv_timer_init(loop.GetNativeLoop(), t); r != 0 {
		C.free(unsafe.Pointer(t))
		return nil, ParseUvErr(r)
	}

	res := &UvTimer{}
	t.data = unsafe.Pointer(&callbackInfo{data: data, ptr: res})
	res.Timer, res.Loop, res.Handle = t, loop.GetNativeLoop(), Handle{(*C.uv_handle_t)(unsafe.Pointer(t)), t.data, res}

	return res, nil
}

// Start (uv_timer_start) start the timer. timeout and repeat are in milliseconds.
func (t *UvTimer) Start(cb func(*Handle, int), timeout uint64, repeat uint64) C.int {
	cbi := (*callbackInfo)(t.Timer.data)
	cbi.timer_cb = cb

	return uv_timer_start(t.Timer, timeout, repeat)
}

// Stop (uv_timer_stop) the timer, the callback will not be called anymore.
func (t *UvTimer) Stop() C.int {
	return C.uv_timer_stop(t.Timer)
}

// Again (uv_timer_again) stop the timer, and if it is repeating restart it using the repeat value as the timeout. If the timer has never been started before it returns UV_EINVAL
func (t *UvTimer) Again() C.int {
	return C.uv_timer_again(t.Timer)
}

// SetRepeat (uv_timer_set_repeat) set the repeat interval value in milliseconds. The timer will be scheduled to run on the given interval, regardless of the callback execution duration, and will follow normal timer semantics in the case of a time-slice overrun.
func (t *UvTimer) SetRepeat(repeat uint64) {
	C.uv_timer_set_repeat(t.Timer, C.uint64_t(repeat))
}

// GetRepeat (uv_timer_get_repeat) get the timer repeat value.
func (t *UvTimer) GetRepeat() uint64 {
	return uint64(C.uv_timer_get_repeat(t.Timer))
}

// Freemem freemem of timer
func (t *UvTimer) Freemem() {
	C.free(unsafe.Pointer(t.Timer))
}
