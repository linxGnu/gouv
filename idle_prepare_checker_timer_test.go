package gouv

import (
	"fmt"
	"log"
	"testing"
	"time"
)

func TestIdlePrepareCheckerTimer(t *testing.T) {
	dfLoop := UvLoopDefault()

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
		time.Sleep(6 * time.Second)

		fmt.Println(timer2.GetRepeat())

		if err := timer2.Stop(); err != nil {
			panic(err)
		}

		timer2.Freemem()
	}()

	go func() {
		time.Sleep(10 * time.Second)

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

	go dfLoop.Run(UV_RUN_DEFAULT)
	time.Sleep(10 * time.Second)
	go dfLoop.Close()
}
