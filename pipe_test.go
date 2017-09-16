package gouv

import (
	"fmt"
	"os"
	"testing"
	"time"
)

func TestPipe(t *testing.T) {
	doTest(t, testPipe, 5)

	// doTestWithLoop(t, testPipe, nil, 5)
}

func testPipe(t *testing.T, loop *UvLoop) {
	pServer, err := UvPipeInit(loop, 0, nil)
	if err != nil {
		t.Fatal(err)
	}

	os.Remove("echo.sock")

	// bind
	if r := pServer.Bind("echo.sock"); r != 0 {
		t.Fatal(ParseUvErr(r))
	}
	if r := pServer.Listen(128, func(h *Handle, status int) {
		server := h.Ptr.(*UvPipe)
		fmt.Println("Pipe server got connection", server, status)

		client, _ := UvPipeInit(loop, 0, nil)
		if r := server.ServerAccept(client.Stream); r != 0 {
			t.Fatal(ParseUvErr(r))
		}

		client.ReadStart(samplePipeReadHandling)
	}); r != 0 {
		t.Fatal(ParseUvErr(r))
	} else {
		fmt.Println("Pipe server is running!")
	}

	//
	go func() {
		defer func() {
			if e := recover(); e != nil {
				fmt.Println(e)
			}
		}()

		time.Sleep(2 * time.Second)

		pClient, _ := UvPipeInit(loop, 0, nil)
		pClient.Connect(NewUvConnect(nil), "echo.sock", func(r *Request, status int) {
			fmt.Println("Connected to pipe server with status: ", status)
		})

		samplePipeReadOfClient(pClient)

		time.Sleep(2 * time.Second)

		shutDown := NewUvShutdown(nil)
		if r := pClient.Shutdown(shutDown.Shutdown, func(h *Request, status int) {
			fmt.Println("Shutdown pipe client!", h, status)
		}); r != 0 {
			t.Fatal(ParseUvErr(r))
		} else {
			fmt.Println("Shutting down pipe client!")
		}

		pClient.ReadStop()

		pClient.Freemem()
	}()
}
