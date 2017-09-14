package gouv

import (
	"fmt"
	"testing"
	"time"
)

func TestPollerFile(t *testing.T) {
	doTest(t, testPoller, 5)
}

func testPoller(t *testing.T, loop *UvLoop) {
	// setup poller
	poller, err := UvPollInit(nil, int(test_OpenFile("test_pkg/tcp_client_sock.c")), nil)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println("Poller file:", poller.GetPollHandle())

	if r := poller.Start(int(UV_READABLE), func(h *Handle, status int, events int) {
		fmt.Println("Poll start callbacked!!!!!", status, events)
	}); r != 0 {
		t.Fatal(ParseUvErr(r))
	}

	go func() {
		time.Sleep(3 * time.Second)

		if r := poller.Stop(); r != 0 {
			t.Fatal(ParseUvErr(r))
		}
	}()
}
