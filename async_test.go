package gouv

import (
	"fmt"
	"testing"
)

func TestAsync(t *testing.T) {
	doTest(t, test_async, 2)
}

func test_async(t *testing.T, dfLoop *UvLoop) {
	async, err := UvAsyncInit(dfLoop, map[int]string{1: "a"}, func(h *Handle) {
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

	if r := async.Send(); r != 0 {
		t.Fatal(ParseUvErr(r))
	}
}
