package gouv

import (
	"fmt"
	"testing"
	"time"
)

func TestTTY(t *testing.T) {
	doTest(t, testTTY, 3)
}

func testTTY(t *testing.T, loop *UvLoop) {
	go func() {
		tty, err := UvTTYInit(loop, 1, 1, nil) // stdout with readable
		if err != nil {
			t.Fatal(err)
		}

		if r := tty.SetMode(UV_TTY_MODE_NORMAL); r != 0 {
			t.Fatal(ParseUvErr(r))
		}

		sampleTTY(tty)

		fmt.Println(tty.GetWinsize())

		time.Sleep(2 * time.Second)

		tty.ResetMode()
	}()
}
