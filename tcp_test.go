package gouv

import (
	"fmt"
	"os"
	"testing"
	"time"
)

func initServer(t *testing.T, loop *UvLoop, flag *uint, port uint16) (connection *UvTCP) {
	addr, err := IPv4Addr("0.0.0.0", port)
	if err != nil {
		t.Fatal(err)
		return
	}

	if connection, err = TCPInit(loop, flag, nil); err != nil {
		t.Fatal(err)
		return
	}

	if err = connection.Bind(addr, 0); err != nil {
		if connection != nil {
			connection.Freemem()
		}

		t.Fatal(err)
		return
	}

	if err = connection.SimultaneousAccepts(1); err != nil {
		if connection != nil {
			connection.Freemem()
		}

		t.Fatal(err)
		return
	}

	if err = connection.Listen(128, func(h *Handle, status int) {
		client, _ := TCPInit(loop, nil, nil)
		client.NoDelay(1)
		client.KeepAlive(1, 10)

		if e := connection.ServerAccept(client.s); e != nil {
			t.Fatal(e)
		}

		fmt.Println("Got connection", client)

		client.ReadStart(sampleTCPReadHandling)
	}); err != nil {
		if connection != nil {
			connection.Freemem()
		}

		t.Fatal(err)
		return
	}

	return
}

func testTCP(t *testing.T, dfLoop *UvLoop) {
	defer os.Remove("/tmp/stderr.txt")
	defer os.Remove("/tmp/stdout.txt")

	connection := initServer(t, dfLoop, nil, 9999)

	clientProcess, err := UvSpawnProcess(dfLoop, &UvProcessOptions{
		Args:  []string{"python", "test_pkg/test_tcp_client.py"},
		Cwd:   "./",
		Flags: UV_PROCESS_DETACHED,
		File:  "python",
		Env:   []string{"PATH"},
		ExitCb: func(h *Handle, status, sigNum int) {
			fmt.Printf("Process client tcp server exited with status %d and signal %d\n", status, sigNum)
			fmt.Printf("%p\n", h.ptr.(*UvProcess))
		},
	}, nil)
	if err != nil {
		t.Fatal(err)
	}

	go func() {
		time.Sleep(2 * time.Second)

		// try to close connection first
		shutDown := NewUvShutdown(nil)
		if err := connection.Shutdown(shutDown.s, func(h *Request, status int) {
			fmt.Println("Shutting down tcp server", h, status)
		}); err != nil {
			t.Fatal(err)
		} else {
			fmt.Println("Shutting down tcp server")
		}

		// Unref this process
		clientProcess.Unref()

		// Try to kill this proces
		clientProcess.Kill(9)
	}()
}
