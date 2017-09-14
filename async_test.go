package gouv

import (
	"fmt"
	"testing"
	"time"
)

func TestAsync(t *testing.T) {
	doTest(t, testAsync, 2)

	doTestWithLoop(t, testAsync, nil, 2)
}

func testAsync(t *testing.T, loop *UvLoop) {
	async, err := UvAsyncInit(loop, map[int]string{1: "a"}, func(h *Handle) {
		if h.Data != nil {
			x := h.Data.(map[int]string)
			if st, ok := x[1]; ok && st == "a" {
				fmt.Println("Got async handle callback:", h.Ptr.(*UvAsync))
				return
			}
		}

		t.Fatalf("Fail on async")
	})
	if err != nil {
		t.Fatal(err)
	}

	go func() {
		if r := async.Send(); r != 0 {
			t.Fatal(ParseUvErr(r))
		}

		time.Sleep(1 * time.Second)

		async.Freemem()
	}()
}
