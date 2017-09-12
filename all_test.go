package gouv

import (
	"fmt"
	"log"
	"testing"
	"time"
)

func TestAll(t *testing.T) {
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

	go loop.Run(UV_RUN_DEFAULT)

	// try to stop
	time.Sleep(200 * time.Millisecond)
	loop.Stop()
	time.Sleep(200 * time.Millisecond)
	if r := loop.Alive(); r > 0 {
		t.Fatalf("Loop not stop well")
	}

	go loop.Run(UV_RUN_NOWAIT)
	time.Sleep(100 * time.Millisecond)
	loop.Stop()

	go loop.Run(UV_RUN_ONCE)
	time.Sleep(100 * time.Millisecond)
	loop.Stop()

	// try to print backend fd and timeout
	fmt.Println(loop.BackendFD(), loop.BackendTimeout())

	//
	loop = UvLoopDefault()
	if err := loop.Init(); err != nil {
		t.Fatal(err)
	}

	//
	testAsync(t, nil)

	//
	testIdlePrepareCheckerTimer(t, nil)

	// //
	testSpawnChildProcess(t, nil)

	// //
	testKillProcess(t, nil)

	//
	testTCP(t, nil)

	go loop.Run(UV_RUN_DEFAULT)

	// try to update time
	loop.UpdateTime()
	n1 := loop.Now()
	time.Sleep(10 * time.Millisecond)
	loop.UpdateTime()
	n2 := loop.Now()
	if n1 == n2 {
		t.Fatalf("Update time fail")
	}

	time.Sleep(20 * time.Second)
	go loop.Close()
}

func testIdlePrepareCheckerTimer(t *testing.T, dfLoop *UvLoop) {
	timer1, err := TimerInit(dfLoop, map[int]string{1: "t1"})
	if err != nil {
		t.Fatal(err)
	}

	timer2, err := TimerInit(dfLoop, map[int]string{2: "t2"})
	if err != nil {
		t.Fatal(err)
	}

	prepare, err := UvPrepareInit(dfLoop, map[int]string{3: "p"})
	if err != nil {
		t.Fatal(err)
	}

	checker, err := UvCheckInit(dfLoop, map[int]string{4: "c"})
	if err != nil {
		t.Fatal(err)
	}

	idle, err := UvIdleInit(dfLoop, map[int]string{5: "i"})
	if err != nil {
		t.Fatal(err)
	}

	prepare.Start(func(h *Handle) {
		if h.Data != nil {
			x := h.Data.(map[int]string)
			if st, ok := x[3]; ok && st == "p" {
				return
			}
		}

		t.Fatalf("Failed on prepare")
	})

	checker.Start(func(h *Handle) {
		if h.Data != nil {
			x := h.Data.(map[int]string)
			if st, ok := x[4]; ok && st == "c" {
				return
			}
		}

		t.Fatalf("Failed on checker")
	})

	idle.Start(func(h *Handle, status int) {
		if h.Data != nil {
			x := h.Data.(map[int]string)
			if st, ok := x[5]; ok && st == "i" {
				return
			}
		}

		t.Fatalf("Failed on idle")
	})

	timer1.Start(func(h *Handle, status int) {
		log.Println(h.Data)
	}, 1000, 1500)

	timer2.Start(func(h *Handle, status int) {
		log.Println(h.Data)
	}, 1000, 1500)

	go func() {
		time.Sleep(2 * time.Second)
		timer1.Again()
		timer1.SetRepeat(1600)

		timer2.Again()
		timer2.SetRepeat(2700)
	}()

	go func() {
		time.Sleep(3 * time.Second)

		fmt.Println(timer2.GetRepeat())

		if err := timer2.Stop(); err != nil {
			t.Fatal(err)
		}

		timer2.Close(func(h *Handle) {
			fmt.Println("Timer is closed")
		})

		timer2.Freemem()
	}()

	go func() {
		time.Sleep(8 * time.Second)

		fmt.Println(timer1.GetRepeat())

		if err := timer1.Stop(); err != nil {
			panic(err)
		}

		timer1.Freemem()

		// now stop idle
		if err := idle.Stop(); err != nil {
			panic(err)
		}

		idle.Freemem()

		// now stop prepare
		if err := prepare.Stop(); err != nil {
			panic(err)
		}

		prepare.Freemem()

		// now stop checker
		if err := checker.Stop(); err != nil {
			panic(err)
		}

		checker.Freemem()
	}()
}
