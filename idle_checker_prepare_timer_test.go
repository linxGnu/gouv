package gouv

import (
	"fmt"
	"log"
	"testing"
	"time"
)

func TestIdlePrepareCheckerTimer(t *testing.T) {
	doTest(t, testIdlePrepareCheckerTimer, 20)
}

func testIdlePrepareCheckerTimer(t *testing.T, dfLoop *UvLoop) {
	defer func() {
		if e := recover(); e != nil {
			fmt.Println(e)
		}
	}()

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
	}, 200, 300)

	timer2.Start(func(h *Handle, status int) {
		log.Println(h.Data)
	}, 200, 300)

	go func() {
		time.Sleep(2 * time.Second)
		timer1.Again()
		timer1.SetRepeat(500)

		timer2.Again()
		timer2.SetRepeat(300)

		// try to stop timer 2
		time.Sleep(2 * time.Second)

		fmt.Println(timer2.GetRepeat())

		if r := timer2.Stop(); r != 0 {
			t.Fatal(ParseUvErr(r))
		}

		timer2.Close(func(h *Handle) {
			fmt.Println("Timer is closed")
		})

		fmt.Println(timer2.IsActive(), timer2.IsClosing())

		timer2.Freemem()

		// try to stop timer 1 and others
		time.Sleep(2 * time.Second)

		fmt.Println(timer1.GetRepeat())

		if r := timer1.Stop(); r != 0 {
			t.Fatal(ParseUvErr(r))
		}

		timer1.Freemem()

		// now stop idle
		if r := idle.Stop(); r != 0 {
			t.Fatal(ParseUvErr(r))
		}

		idle.Freemem()

		// now stop prepare
		if r := prepare.Stop(); r != 0 {
			t.Fatal(ParseUvErr(r))
		}

		prepare.Freemem()

		// now stop checker
		if r := checker.Stop(); r != 0 {
			t.Fatal(ParseUvErr(r))
		}

		checker.Freemem()
	}()
}
