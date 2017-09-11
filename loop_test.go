package gouv

import (
	"fmt"
	"testing"
	"time"
)

func TestNewUvLoop(t *testing.T) {
	loop := UvLoopDefault()
	if loop == nil {
		t.Fatalf("NewUvLoopDefault failed")
	}

	if err := loop.Init(); err != nil {
		t.Fatal(err)
	}

	if r := loop.GetNativeLoop(); r == nil {
		t.Fatalf("NewUvLoopDefault failed")
	}

	if err := loop.Run(UVRUNDEFAULT); err != nil {
		t.Fatal(err)
	}

	// try to stop
	loop.Stop()
	time.Sleep(2 * time.Second)
	if r := loop.Alive(); r > 0 {
		t.Fatalf("Loop not stop well")
	}

	// Try to rerun
	if err := loop.Run(UVRUNONCE); err != nil {
		t.Fatal(err)
	}

	if err := loop.Run(UVRUNNOWAIT); err != nil {
		t.Fatal(err)
	}

	// try to update time
	loop.UpdateTime()
	n1 := loop.Now()
	time.Sleep(1 * time.Second)
	loop.UpdateTime()
	n2 := loop.Now()
	if n1 == n2 {
		t.Fatalf("Update time fail")
	}

	// try to print backend fd and timeout
	fmt.Println(loop.BackendFD(), loop.BackendTimeout())

	if err := loop.Close(); err != nil {
		t.Fatal(err)
	}
}
