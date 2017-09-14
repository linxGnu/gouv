package gouv

import (
	"fmt"
	"testing"
	"time"
)

func TestLoopStartStop(t *testing.T) {
	loop := UvLoopDefault()
	if loop == nil {
		t.Fatalf("NewUvLoopDefault failed")
	}

	if r := loop.Init(); r != 0 {
		t.Fatal(ParseUvErr(r))
	}

	go loop.Run(UV_RUN_DEFAULT)

	// try to print backend fd and timeout
	fmt.Println(loop.BackendFD(), loop.BackendTimeout())

	// try to stop
	time.Sleep(200 * time.Millisecond)
	loop.Stop()
	time.Sleep(200 * time.Millisecond)
	if r := loop.Alive(); r > 0 {
		t.Fatalf("Loop not stop well")
	}

	loop.Run(UV_RUN_NOWAIT)
	time.Sleep(200 * time.Millisecond)
	loop.Stop()

	go loop.Run(UV_RUN_ONCE)
	time.Sleep(200 * time.Millisecond)
	loop.Stop()

	//
	loop = UvLoopDefault()
	if r := loop.Init(); r != 0 {
		t.Fatal(ParseUvErr(r))
	}

	// try to update time
	loop.UpdateTime()
	n1 := loop.Now()
	time.Sleep(10 * time.Millisecond)
	loop.UpdateTime()
	n2 := loop.Now()
	if n1 == n2 {
		t.Fatalf("Update time fail")
	}

	time.Sleep(10 * time.Millisecond)
	loop.Stop()
	loop.Close()
}
