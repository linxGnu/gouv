package binduv

import (
	"fmt"
	"strings"
	"testing"
)

func testAddrV4(host string, expected string) error {
	// create new address
	addr, err := IPv4Addr(host, 3306)
	if err != nil {
		return err
	}

	// check name
	st, err := addr.Name()
	if err != nil {
		return err
	}

	expected = strings.ToLower(expected)

	h := []byte(expected)
	if len(st) < len(h) {
		return fmt.Errorf("Fail addr %x %d", st, len(st))
	}

	// compare address
	for i := range h {
		if i >= len(h) {
			if st[i] != 0 {
				return fmt.Errorf("Fail addr %x %d", st, len(st))
			}
		} else if st[i] != h[i] {
			return fmt.Errorf("Fail addr %x %d", st, len(st))
		}
	}

	addr.Freemem()
	return nil
}

func TestIp4Addr(t *testing.T) {
	host := "121.222.123.144"
	if err := testAddrV4(host, host); err != nil {
		t.Fatal(err)
	}

	host = "121.256.123.144"
	if err := testAddrV4(host, "0.0.0.0"); err == nil {
		t.Fatalf("Fail checking error of testAddrV4")
	} else {
		fmt.Println(err)
	}
}

func testAddrV6(host string, expected string) error {
	// create new address
	addr, err := IPv6Addr(host, 3306)
	if err != nil {
		return err
	}

	// check name
	st, err := addr.Name()
	if err != nil {
		return err
	}

	expected = strings.ToLower(expected)

	h := []byte(expected)
	if len(st) < len(h) {
		return fmt.Errorf("Fail addr %s %d", string(st), len(st))
	}

	// compare address
	for i := range h {
		if i >= len(h) {
			if st[i] != 0 {
				return fmt.Errorf("Fail addr %s %d", string(st), len(st))
			}
		} else if st[i] != h[i] {
			return fmt.Errorf("Fail addr %s %d", string(st), len(st))
		}
	}

	addr.Freemem()
	return nil
}

func TestIp6Addr(t *testing.T) {
	host := "2001:dc8::1005:2f43:bcd:ffff"
	if err := testAddrV6(host, host); err != nil {
		t.Fatal(err)
	}

	// overflow checking
	host = "2001:dc8::1005:2f43:bcd:ffff2001:dc8::1005:2f43:bcd:ffff2001:dc8::1005:2f43:bcd:ffff2001:dc8::1005:2f43:bcd:ffff2001:dc8::1005:2f43:bcd:ffff2001:dc8::1005:2f43:bcd:ffff2001:dc8::1005:2f43:bcd:ffff2001:dc8::1005:2f43:bcd:ffff2001:dc8::1005:2f43:bcd:ffff2001:dc8::1005:2f43:bcdasdfasdf"
	if err := testAddrV6(host, "::"); err == nil {
		t.Fatalf("Fail checking error of testAddrV6")
	} else {
		fmt.Println(err)
	}
}
