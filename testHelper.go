package gouv

// #cgo pkg-config: libuv
/*
#include <uv.h>
#include <stdlib.h>
char* testRead(uv_stream_t *client, ssize_t nread, uv_buf_t* buf) {
	char* tmp;
	tmp = malloc(nread + 1);
	memcpy(tmp, buf->base, nread);
	tmp[nread] = '\0';

	return tmp;
}
*/
import "C"

import (
	"fmt"
	"unsafe"
)

func sampleTCPReadHandling(h *Handle, buf *C.uv_buf_t, nRead C.ssize_t) {
	conn := h.Ptr.(*UvTCP)

	st := C.testRead(conn.s, nRead, buf)
	fmt.Println("Read from client: ", C.GoString(st))

	bufs := MallocUvBuf(1)
	SetBuf(bufs, 0, BufInit2(st, C.uint(C.strlen(st)+1)))

	conn.Write(NewUvWrite(nil).w, bufs, 1, func(h *Request, status int) {
		fmt.Println("Write done: ", h.Handle.Ptr.(*UvTCP), status)
	})
}

func sampleTCPReadOfClient(conn *UvTCP) {
	st := C.CString("Hello server")
	defer C.free(unsafe.Pointer(st))

	bufs := MallocUvBuf(1)
	SetBuf(bufs, 0, BufInit2(st, C.uint(C.strlen(st)+1)))

	conn.Write(NewUvWrite(nil).w, bufs, 1, func(h *Request, status int) {
		fmt.Println("Write done: ", h.Handle.Ptr.(*UvTCP), status)
	})

	conn.ReadStart(func(h *Handle, buf *C.uv_buf_t, nRead C.ssize_t) {
		st := C.testRead(conn.s, nRead, buf)
		fmt.Println("Read from server ______ :", C.GoString(st))
	})
}

func samplePipeReadHandling(h *Handle, buf *C.uv_buf_t, nRead C.ssize_t) {
	conn := h.Ptr.(*UvPipe)

	st := C.testRead(conn.s, nRead, buf)
	fmt.Println("Read from client: ", C.GoString(st))

	bufs := MallocUvBuf(1)
	SetBuf(bufs, 0, BufInit2(st, C.uint(C.strlen(st)+1)))

	conn.Write(NewUvWrite(nil).w, bufs, 1, func(h *Request, status int) {
		fmt.Println("Write done: ", h.Handle.Ptr.(*UvPipe), status)
	})
}

func samplePipeReadOfClient(conn *UvPipe) {
	st := C.CString("Hello pipe server")
	defer C.free(unsafe.Pointer(st))

	bufs := MallocUvBuf(1)
	SetBuf(bufs, 0, BufInit2(st, C.uint(C.strlen(st)+1)))

	conn.Write(NewUvWrite(nil).w, bufs, 1, func(h *Request, status int) {
		fmt.Println("Write pipe done: ", h.Handle.Ptr.(*UvPipe), status)
	})

	conn.ReadStart(func(h *Handle, buf *C.uv_buf_t, nRead C.ssize_t) {
		st := C.testRead(conn.s, nRead, buf)
		fmt.Println("Read from pipe server ______ :", C.GoString(st))
	})
}
