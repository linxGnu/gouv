package gouv

import (
	"fmt"
	"testing"
	"time"
)

func initServer(t *testing.T, loop *UvLoop) (connection *UvTCP) {
	addr, err := IPv4Addr("0.0.0.0", 9999)
	if err != nil {
		t.Fatal(err)
		return
	}

	if connection, err = TCPInit(loop, nil); err != nil {
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

	if err = connection.Listen(128, func(h *Handle, status int) {
		client, _ := TCPInit(loop, nil)
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

func TestTCP(t *testing.T) {
	dfLoop := UvLoopDefault()

	connection := initServer(t, dfLoop)

	go dfLoop.Run(UV_RUN_DEFAULT)

	time.Sleep(5 * time.Second)

	// try to close connection first
	shutDown := NewUvShutdown(nil)
	if err := connection.Shutdown(shutDown.s, func(h *Request, status int) {
		fmt.Println(h, status)
	}); err != nil {
		t.Fatal(err)
	}

	time.Sleep(2 * time.Second)

	go dfLoop.Close()
}
