package gouv

import (
	"bufio"
	"os"
	"testing"
	"time"
)

func TestFSPoll(t *testing.T) {
	doTest(t, testFSPoll, 10)
}

func testFSPoll(t *testing.T, loop *UvLoop) {
	ev, err := UvFSPollInit(loop, nil)
	if err != nil {
		t.Fatal(err)
	}

	if r, _, _ := ev.GetPath(); r == 0 {
		t.Fatal("Logic of FS Poll invalid")
	}

	os.Remove("fs_poll.tmp")
	os.Create("fs_poll.tmp")

	sampleFSPoll(ev, "fs_poll.tmp")

	if r, _, _ := ev.GetPath(); r != 0 {
		t.Fatal(ParseUvErr(r))
	}

	go func() {
		file, _ := os.OpenFile("fs_poll.tmp", os.O_RDWR, 0777)
		fileWriter := bufio.NewWriter(file)
		fileWriter.WriteString("hello world!")
		fileWriter.Flush()
		file.Close()

		time.Sleep(2 * time.Second)

		file, _ = os.OpenFile("fs_poll.tmp", os.O_RDWR, 0777)
		fileWriter = bufio.NewWriter(file)
		fileWriter.WriteString("hello world!")
		fileWriter.Flush()
		file.Close()

		time.Sleep(6 * time.Second)

		// Stop fs event
		ev.Stop()

		// free mem
		ev.Freemem()
	}()
}
