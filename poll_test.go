package gouv

import (
	"fmt"
	"testing"
	"time"
)

func TestPollerFile(t *testing.T) {
	doTest(t, testPollerFile, 3)
}

func testPollerFile(t *testing.T, loop *UvLoop) {
	// setup poller
	poller, err := UvPollInit(loop, 1, nil)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println("Poller file:", poller.GetPollHandle())

	if r := poller.Start(int(UV_READABLE|UV_WRITABLE), func(h *Handle, status int, events int) {
		fmt.Println("Poll callbacked!!!!!", status, events)
	}); r != 0 {
		t.Fatal(ParseUvErr(r))
	}

	go func() {
		time.Sleep(2 * time.Second)

		if r := poller.Stop(); r != 0 {
			t.Fatal(ParseUvErr(r))
		}

		poller.Freemem()
	}()
}
