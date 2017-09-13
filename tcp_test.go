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

func testTCP(t *testing.T, dfLoop *UvLoop) {
	defer func() {
		if e := recover(); e != nil {
			fmt.Println(e)
		}
	}()

	defer os.Remove("/tmp/stderr.txt")
	defer os.Remove("/tmp/stdout.txt")

	server := initServer(t, dfLoop, nil, testServerPort)

	go runPythonClient(t, dfLoop)

	go runSockClient(t, dfLoop)

	go runUvTcpClient(t, dfLoop)

	go func() {
		time.Sleep(20 * time.Second)

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

func runPythonClient(t *testing.T, dfLoop *UvLoop) {
	defer func() {
		if e := recover(); e != nil {
			fmt.Println(e)
		}
	}()

	pythonClient, err := UvSpawnProcess(dfLoop, &UvProcessOptions{
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

		// Try to kill this proces
		pythonClient.Kill(9)
	}()
}

func runUvTcpClient(t *testing.T, dfLoop *UvLoop) {
	defer func() {
		if e := recover(); e != nil {
			fmt.Println(e)
		}
	}()

	//
	time.Sleep(1 * time.Second)

	serverAddr, err := IPv4Addr("127.0.0.1", testServerPort)
	if err != nil {
		t.Fatal(err)
		return
	}

	//
	tcp, _ := TCPInit(dfLoop, nil, nil)

	if r := tcp.Connect(NewUvConnect(nil), serverAddr, func(h *Request, status int) {
		conn := h.Handle.Ptr.(*UvTCP)
		fmt.Println("Connected to server", conn)
	}); r != 0 {
		t.Fatal(ParseUvErr(r))
	}

	sampleTCPReadOfClient(tcp)

	time.Sleep(5 * time.Second)

	shutDown := NewUvShutdown(nil)
	if r := tcp.Shutdown(shutDown.s, func(h *Request, status int) {
		fmt.Println("Shutdown uv_tcp_t client!", h, status)
	}); r != 0 {
		t.Fatal(ParseUvErr(r))
	} else {
		fmt.Println("Shutting down uv_tcp_t client!")
	}

	tcp.ReadStop()
}

func runSockClient(t *testing.T, dfLoop *UvLoop) {
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
	sock := create_tcp_socket(serverAddr, 0)

	// connect socket first
	if r := connect_socket(sock, serverAddr); r != 0 {
		t.Fatal(ParseUvErr(r))
	}

	//
	poller, err := UvPollInitSocket(dfLoop, sock, nil)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println("Poller sock:", poller.GetPollHandle())

	if r := poller.Start(int(UV_READABLE|UV_WRITABLE), func(h *Handle, status int, events int) {
		fmt.Println("Poll start callbacked!!!!!", status, events)
	}); r != 0 {
		t.Fatal(ParseUvErr(r))
	}

	// now try to send and recv
	test_sendAndRecv(sock)

	// Close socket
	if r := close_socket(sock); r != 0 {
		t.Fatal(ParseUvErr(r))
	}

	// Stop poller
	// if r := poller.Stop(); r != 0 {
	// 	t.Fatal(ParseUvErr(r))
	// }
}
