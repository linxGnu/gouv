package gouv

import (
	"fmt"
	"testing"
	"time"
)

func testAsync(t *testing.T, dfLoop *UvLoop) {
	async, err := UvAsyncInit(dfLoop, map[int]string{1: "a"}, func(h *Handle) {
		if h.Data != nil {
			x := h.Data.(map[int]string)
			if st, ok := x[1]; ok && st == "a" {
				fmt.Println(h.ptr.(*UvAsync))
				return
			}
		}

		t.Fatalf("Fail on async")
	})
	if err != nil {
		t.Fatal(err)
	}

	go func() {
		time.Sleep(100 * time.Millisecond)
		if r := async.Send(); r != 0 {
			t.Fatal(ParseUvErr(r))
		}
	}()
}
