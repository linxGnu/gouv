package gouv

import (
	"fmt"
	"os"
	"testing"
	"time"
)

const (
	testServerPort = 9999
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

		if r := connection.ServerAccept(client.GetStreamHandle()); r != 0 {
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

func TestTCP(t *testing.T) {
	doTest(t, testTCP, 15)
}

func testTCP(t *testing.T, loop *UvLoop) {
	defer func() {
		if e := recover(); e != nil {
			fmt.Println(e)
		}
	}()

	defer os.Remove("/tmp/stderr.txt")
	defer os.Remove("/tmp/stdout.txt")

	server := initServer(t, loop, nil, testServerPort)

	go runPythonClient(t, loop)

	go runUvTcpClient(t, loop, testServerPort)

	go func() {
		time.Sleep(6 * time.Second)

		// try to close connection first
		shutDown := NewUvShutdown(nil)
		if r := server.Shutdown(shutDown.s, func(h *Request, status int) {
			fmt.Println("Shutting down tcp server", h, status)
		}); r != 0 {
			t.Fatal(ParseUvErr(r))
		} else {
			fmt.Println("Shutting down tcp server")
		}

		// stop read
		if r := server.ReadStop(); r != 0 {
			t.Fatal(ParseUvErr(r))
		}
	}()
}

func TestTCP2(t *testing.T) {
	doTestWithLoop(t, testTCP2, nil, 10)
}

func testTCP2(t *testing.T, loop *UvLoop) {
	defer func() {
		if e := recover(); e != nil {
			fmt.Println(e)
		}
	}()

	var flags uint = 0
	server := initServer(t, loop, &flags, 10000)

	go runUvTcpClient(t, loop, 10000)

	go func() {
		time.Sleep(6 * time.Second)

		// try to close connection first
		shutDown := NewUvShutdown(nil)
		if r := server.Shutdown(shutDown.s, func(h *Request, status int) {
			fmt.Println("Shutting down tcp server", h, status)
		}); r != 0 {
			t.Fatal(ParseUvErr(r))
		} else {
			fmt.Println("Shutting down tcp server")
		}

		// stop read
		if r := server.ReadStop(); r != 0 {
			t.Fatal(ParseUvErr(r))
		}
	}()
}

func runPythonClient(t *testing.T, loop *UvLoop) {
	defer func() {
		if e := recover(); e != nil {
			fmt.Println(e)
		}
	}()

	pythonClient, err := UvSpawnProcess(loop, &UvProcessOptions{
		Args:  []string{"python", "test_pkg/test_tcp_client.py"},
		Cwd:   "./",
		Flags: UV_PROCESS_DETACHED,
		File:  "python",
		Env:   []string{"PATH"},
		ExitCb: func(h *Handle, status, sigNum int) {
			fmt.Printf("Process client tcp server exited with status %d and signal %d\n", status, sigNum)
			fmt.Printf("%p\n", h.Ptr.(*UvProcess))
		},
	}, nil)
	if err != nil {
		t.Fatal(err)
	}

	go func() {
		time.Sleep(5 * time.Second)

		// Unref this process
		pythonClient.Unref()

		// Try to kill this process
		pythonClient.Kill(9)
	}()
}

func runUvTcpClient(t *testing.T, loop *UvLoop, testServerPort uint16) {
	defer func() {
		if e := recover(); e != nil {
			fmt.Println(e)
		}
	}()

	serverAddr, err := IPv4Addr("127.0.0.1", testServerPort)
	if err != nil {
		t.Fatal(err)
		return
	}

	//
	tcp, _ := TCPInit(loop, nil, nil)

	cnRe := NewUvConnect(nil)
	if r := tcp.Connect(cnRe, serverAddr, func(h *Request, status int) {
		conn := h.Handle.Ptr.(*UvTCP)

		fmt.Println("Connected to server", conn)

		cnRe.Freemem()
	}); r != 0 {
		t.Fatal(ParseUvErr(r))
	}

	sampleTCPReadOfClient(tcp)

	go func() {
		time.Sleep(2 * time.Second)

		shutDown := NewUvShutdown(nil)
		if r := tcp.Shutdown(shutDown.s, func(h *Request, status int) {
			fmt.Println("Shutdown uv_tcp_t client!", h, status)
		}); r != 0 {
			t.Fatal(ParseUvErr(r))
		} else {
			fmt.Println("Shutting down uv_tcp_t client!")
		}

		time.Sleep(1 * time.Second)

		shutDown.Freemem()

		tcp.ReadStop()
	}()
}
