package gouv

import (
	"fmt"
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

	if r := connection.Bind(addr, 0); r != 0 {
		if connection != nil {
			connection.Freemem()
		}

		t.Fatal(ParseUvErr(r))
		return
	}

	fmt.Println(connection.IsReadable(), connection.IsWritable())

	if r := connection.SimultaneousAccepts(1); r != 0 {
		if connection != nil {
			connection.Freemem()
		}

		t.Fatal(ParseUvErr(r))
		return
	}

	if r := connection.Listen(128, func(h *Handle, status int) {
		client, _ := TCPInit(loop, nil, nil)
		client.NoDelay(1)
		client.KeepAlive(1, 10)

		if r := connection.ServerAccept(client.s); r != 0 {
			t.Fatal(ParseUvErr(r))
		}

		fmt.Println("Got connection", client)

		client.ReadStart(sampleTCPReadHandling)
	}); r != 0 {
		if connection != nil {
			connection.Freemem()
		}

		t.Fatal(ParseUvErr(r))
		return
	}

	return
}

func testTCP(t *testing.T, dfLoop *UvLoop) {
	// defer os.Remove("/tmp/stderr.txt")
	// defer os.Remove("/tmp/stdout.txt")

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
		time.Sleep(5 * time.Second)

		// try to close connection first
		shutDown := NewUvShutdown(nil)
		if r := connection.Shutdown(shutDown.s, func(h *Request, status int) {
			fmt.Println("Shutting down tcp server", h, status)
		}); r != 0 {
			t.Fatal(ParseUvErr(r))
		} else {
			fmt.Println("Shutting down tcp server")
		}

		// stop read
		if r := connection.ReadStop(); r != 0 {
			t.Fatal(ParseUvErr(r))
		}

		// Unref this process
		clientProcess.Unref()

		// Try to kill this proces
		clientProcess.Kill(9)
	}()
}