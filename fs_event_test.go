package gouv

import (
	"bufio"
	"os"
	"testing"
	"time"
)

func TestFSEvent(t *testing.T) {
	doTest(t, testFSEvent, 2)
}

func testFSEvent(t *testing.T, loop *UvLoop) {

	ev, err := UvFSEventInit(loop, nil)
	if err != nil {
		t.Fatal(err)
	}

	if r, _, _ := ev.GetPath(); r == 0 {
		t.Fatal("Logic of FS Event invalid")
	}

	os.Remove("tmp.tmp")
	file, _ := os.Create("tmp.tmp")
	defer file.Close()

	sampleFSEvent(ev)

	fileWriter := bufio.NewWriter(file)
	fileWriter.WriteString("hello world!")
	fileWriter.Flush()

	if r, _, _ := ev.GetPath(); r != 0 {
		t.Fatal(ParseUvErr(r))
	}

	go func() {
		time.Sleep(2 * time.Second)

		// Stop fs event
		ev.Stop()

		// free mem
		ev.Freemem()
	}()
}
