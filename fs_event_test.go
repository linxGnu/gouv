package gouv

import (
	"bufio"
	"os"
	"testing"
	"time"
)

func TestFSEvent(t *testing.T) {
	doTest(t, testFSEvent, 10)
}

func testFSEvent(t *testing.T, loop *UvLoop) {

	ev, err := UvFSEventInit(loop, nil)
	if err != nil {
		t.Fatal(err)
	}

	if r, _, _ := ev.GetPath(); r == 0 {
		t.Fatal("Logic of FS Event invalid")
	}

	os.Remove("fs_event.tmp")
	os.Create("fs_event.tmp")

	sampleFSEvent(ev, "fs_event.tmp")

	if r, _, _ := ev.GetPath(); r != 0 {
		t.Fatal(ParseUvErr(r))
	}

	go func() {
		file, _ := os.OpenFile("fs_event.tmp", os.O_RDWR, 0777)
		fileWriter := bufio.NewWriter(file)
		fileWriter.WriteString("hello world!")
		fileWriter.Flush()
		file.Close()

		time.Sleep(2 * time.Second)

		file, _ = os.OpenFile("fs_event.tmp", os.O_RDWR, 0777)
		fileWriter = bufio.NewWriter(file)
		fileWriter.WriteString("hello world!")
		fileWriter.Flush()
		file.Close()

		time.Sleep(2 * time.Second)

		// Stop fs event
		ev.Stop()

		// free mem
		ev.Freemem()
	}()
}
