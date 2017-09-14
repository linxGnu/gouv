package gouv

import (
	"fmt"
	"testing"
	"time"
)

func doTest(t *testing.T, run func(t *testing.T, loop *UvLoop), timeout int) {
	loop := UvLoopCreate()

	fmt.Println("Loop at:", loop.GetNativeLoop())

	// do real test
	run(t, loop)

	// run but not blocking
	go loop.Run(UV_RUN_DEFAULT)
	// loop.Run(UV_RUN_NOWAIT)

	// wait to stop
	time.Sleep(time.Duration(timeout) * time.Second)
	loop.Stop()
	loop.Close()
}

func doTestWithLoop(t *testing.T, run func(t *testing.T, loop *UvLoop), loop *UvLoop, timeout int) {
	// do real test
	run(t, loop)

	loop = UvLoopDefault()

	// run but not blocking
	go loop.Run(UV_RUN_DEFAULT)
	// loop.Run(UV_RUN_NOWAIT)

	// wait to stop
	time.Sleep(time.Duration(timeout) * time.Second)
	loop.Stop()
	loop.Close()
}
