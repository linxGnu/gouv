package gouv

import (
	"fmt"
	"testing"
)

func TestSignal(t *testing.T) {
	doTest(t, testSignal, 2)
}

func testSignal(t *testing.T, loop *UvLoop) {
	go func() {
		signal, err := UvSignalInit(loop, nil)
		if err != nil {
			t.Fatal(err)
		}

		fmt.Println("Signal handle:", signal.GetSignalHandle())
		signal.Start(func(h *Handle, sigNum int) {
			fmt.Println("Receive signal:", sigNum)
		}, 1)

		signal1, err := UvSignalInit(nil, nil)
		if err != nil {
			t.Fatal(err)
		}

		fmt.Println("Signal handle:", signal.GetSignalHandle())
		signal1.Start(func(h *Handle, sigNum int) {
			fmt.Println("Receive signal:", sigNum)
		}, 2)
	}()
}
