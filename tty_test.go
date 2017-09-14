package gouv

import (
	"fmt"
	"testing"
)

func TestTTY(t *testing.T) {
	doTest(t, testTTY, 3)
}

func testTTY(t *testing.T, loop *UvLoop) {
	go func() {
		defer func() {
			if e := recover(); e != nil {
				fmt.Println(e)
			}
		}()

		tty, err := UvTTYInit(loop, 1, 1, nil) // stdout with readable
		if err != nil {
			t.Fatal(err)
		}

		fmt.Println(tty.GetWinsize())

		// if r := tty.SetMode(UV_TTY_MODE_NORMAL); r != 0 {
		// 	t.Fatal(ParseUvErr(r))
		// }

		// sampleTTY(tty)

		// time.Sleep(2 * time.Second)

		// tty.ResetMode()

		// tty.Freemem()
	}()
}
