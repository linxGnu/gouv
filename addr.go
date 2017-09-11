package gouv

// #cgo pkg-config: libuv
/*
#include <uv.h>
#include <stdlib.h>
struct sockaddr_in* mallocSockaddr_in() {
	return (struct sockaddr_in*)malloc(sizeof(struct sockaddr_in));
}
struct sockaddr_in6* mallocSockaddr_in6() {
	return (struct sockaddr_in6*)malloc(sizeof(struct sockaddr_in6));
}
*/
import "C"

import (
	"unsafe"
)

// SockaddrIn sock addr interface
type SockaddrIn interface {
	Name() ([]byte, error)
	Freemem()
}

// Sockaddr general sock address
type Sockaddr struct {
	sa *C.struct_sockaddr
}

// SockaddrIn4 sock address ipv4
type SockaddrIn4 struct {
	sa *C.struct_sockaddr_in
}

// SockaddrIn6 sock address ipv6
type SockaddrIn6 struct {
	sa *C.struct_sockaddr_in6
}

// IPv4Addr make new SockAddr (ipv4) from host and port
func IPv4Addr(host string, port uint16) (SockaddrIn, error) {
	phost := C.CString(host)
	defer C.free(unsafe.Pointer(phost))

	addr := C.mallocSockaddr_in()
	if r := C.uv_ip4_addr(phost, C.int(port), addr); r != 0 {
		return nil, ParseUvErr(r)
	}

	return &SockaddrIn4{addr}, nil
}

// IPv6Addr make new SockAddr (ipv4) from host and port
func IPv6Addr(host string, port uint16) (SockaddrIn, error) {
	phost := C.CString(host)
	defer C.free(unsafe.Pointer(phost))

	addr := C.mallocSockaddr_in6()
	if r := C.uv_ip6_addr(phost, C.int(port), addr); r != 0 {
		return nil, ParseUvErr(r)
	}

	return &SockaddrIn6{addr}, nil
}

// Name name of sock addr
func (sa *SockaddrIn4) Name() (name []byte, err error) {
	b := make([]byte, 256)
	if r := C.uv_ip4_name(sa.sa, (*C.char)(unsafe.Pointer(&b[0])), C.size_t(len(b))); r != 0 {
		return nil, ParseUvErr(r)
	}

	return b, nil
}

// Freemem of addr
func (sa *SockaddrIn4) Freemem() {
	C.free(unsafe.Pointer(sa.sa))
}

// Name name of sock addr
func (sa *SockaddrIn6) Name() (name []byte, err error) {
	b := make([]byte, 256)
	if r := C.uv_ip6_name(sa.sa, (*C.char)(unsafe.Pointer(&b[0])), C.size_t(len(b))); r != 0 {
		return nil, ParseUvErr(r)
	}

	return b, nil
}

// Freemem of addr
func (sa *SockaddrIn6) Freemem() {
	C.free(unsafe.Pointer(sa.sa))
}
