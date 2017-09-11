package gouv

import (
	"testing"
	"time"
)

func TestAsync(t *testing.T) {
	dfLoop := UvLoopDefault()

	async, err := UvAsyncInit(dfLoop, map[int]string{1: "a"}, func(h *Handle) {
		if h.Data != nil {
			x := h.Data.(map[int]string)
			if st, ok := x[1]; ok && st == "a" {
				return
			}
		}

		t.Fatalf("Fail on async")
	})
	if err != nil {
		t.Fatal(err)
	}

	go func() {
		time.Sleep(1 * time.Second)
		if err = async.Send(); err != nil {
			t.Fatal(err)
		}
	}()

	go dfLoop.Run(UVRUNDEFAULT)
	time.Sleep(2 * time.Second)
	go dfLoop.Close()
}
