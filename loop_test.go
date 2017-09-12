package gouv

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"testing"
	"time"
)

func TestAll(t *testing.T) {
	loop := UvLoopDefault()
	if loop == nil {
		t.Fatalf("NewUvLoopDefault failed")
	}

	if err := loop.Init(); err != nil {
		t.Fatal(err)
	}

	if r := loop.GetNativeLoop(); r == nil {
		t.Fatalf("NewUvLoopDefault failed")
	}

	// go loop.Run(UV_RUN_DEFAULT)

	// // try to stop
	// time.Sleep(100 * time.Millisecond)
	// loop.Stop()
	// time.Sleep(100 * time.Millisecond)
	// if r := loop.Alive(); r > 0 {
	// 	t.Fatalf("Loop not stop well")
	// }

	// go loop.Run(UV_RUN_NOWAIT)
	// time.Sleep(100 * time.Millisecond)
	// loop.Stop()

	// go loop.Run(UV_RUN_ONCE)
	// time.Sleep(100 * time.Millisecond)
	// loop.Stop()

	// // try to print backend fd and timeout
	// fmt.Println(loop.BackendFD(), loop.BackendTimeout())

	//
	testAsync(t, loop)

	//
	testIdlePrepareCheckerTimer(t, loop)

	//
	testSpawnChildProcess(t, loop)

	//
	testKillProcess(t, loop)

	//
	testTCP(t, loop)

	//
	testTCPEx(t, loop)

	loop.Run(UV_RUN_DEFAULT)

	// try to update time
	// loop.UpdateTime()
	// n1 := loop.Now()
	// time.Sleep(10 * time.Millisecond)
	// loop.UpdateTime()
	// n2 := loop.Now()
	// if n1 == n2 {
	// 	t.Fatalf("Update time fail")
	// }

	time.Sleep(20 * time.Second)
	loop.Close()
}

func testAsync(t *testing.T, dfLoop *UvLoop) {
	async, err := UvAsyncInit(dfLoop, map[int]string{1: "a"}, func(h *Handle) {
		if h.Data != nil {
			x := h.Data.(map[int]string)
			if st, ok := x[1]; ok && st == "a" {
				fmt.Println(h.ptr.(*UvAsync))
				return
			}
		}

		t.Fatalf("Fail on async")
	})
	if err != nil {
		t.Fatal(err)
	}

	go func() {
		time.Sleep(1 * time.Second)
		if err = async.Send(); err != nil {
			t.Fatal(err)
		}
	}()
}

func testIdlePrepareCheckerTimer(t *testing.T, dfLoop *UvLoop) {
	timer1, err := TimerInit(dfLoop, map[int]string{1: "t1"})
	if err != nil {
		t.Fatal(err)
	}

	timer2, err := TimerInit(dfLoop, map[int]string{2: "t2"})
	if err != nil {
		t.Fatal(err)
	}

	prepare, err := UvPrepareInit(dfLoop, map[int]string{3: "p"})
	if err != nil {
		t.Fatal(err)
	}

	checker, err := UvCheckInit(dfLoop, map[int]string{4: "c"})
	if err != nil {
		t.Fatal(err)
	}

	idle, err := UvIdleInit(dfLoop, map[int]string{5: "i"})
	if err != nil {
		t.Fatal(err)
	}

	prepare.Start(func(h *Handle) {
		if h.Data != nil {
			x := h.Data.(map[int]string)
			if st, ok := x[3]; ok && st == "p" {
				return
			}
		}

		t.Fatalf("Failed on prepare")
	})

	checker.Start(func(h *Handle) {
		if h.Data != nil {
			x := h.Data.(map[int]string)
			if st, ok := x[4]; ok && st == "c" {
				return
			}
		}

		t.Fatalf("Failed on checker")
	})

	idle.Start(func(h *Handle, status int) {
		if h.Data != nil {
			x := h.Data.(map[int]string)
			if st, ok := x[5]; ok && st == "i" {
				return
			}
		}

		t.Fatalf("Failed on idle")
	})

	timer1.Start(func(h *Handle, status int) {
		log.Println(h.Data)
	}, 1000, 1500)

	timer2.Start(func(h *Handle, status int) {
		log.Println(h.Data)
	}, 1000, 1500)

	go func() {
		time.Sleep(2 * time.Second)
		timer1.Again()
		timer1.SetRepeat(1600)

		timer2.Again()
		timer2.SetRepeat(2700)
	}()

	go func() {
		time.Sleep(3 * time.Second)

		fmt.Println(timer2.GetRepeat())

		if err := timer2.Stop(); err != nil {
			panic(err)
		}

		timer2.Freemem()
	}()

	go func() {
		time.Sleep(2 * time.Second)

		fmt.Println(timer1.GetRepeat())

		if err := timer1.Stop(); err != nil {
			panic(err)
		}

		timer1.Freemem()

		// now stop idle
		if err := idle.Stop(); err != nil {
			panic(err)
		}

		idle.Freemem()

		// now stop prepare
		if err := prepare.Stop(); err != nil {
			panic(err)
		}

		prepare.Freemem()

		// now stop checker
		if err := checker.Stop(); err != nil {
			panic(err)
		}

		checker.Freemem()
	}()
}

func testSpawnChildProcess(t *testing.T, dfLoop *UvLoop) {
	UvDisableStdioInheritance()

	// spawn new process
	process, err := UvSpawnProcess(dfLoop, &UvProcessOptions{
		Args:  []string{"ls", "-lah"},
		Cwd:   "/tmp",
		Flags: UV_PROCESS_DETACHED,
		File:  "ls",
		ExitCb: func(h *Handle, status, sigNum int) {
			if status != 0 {
				t.Fatalf("Failed spawn child process")
			}

			fmt.Printf("Process exited with status %d and signal %d\n", status, sigNum)
		},
	}, nil)
	if err != nil {
		t.Fatal(err)
	}

	time.Sleep(1 * time.Second)

	// Unref this process
	process.Unref()

	// kill this process
	// process.Kill(9)
}

// fileExists reports whether the named file or directory exists.
func fileExists(name string) bool {
	if _, err := os.Stat(name); err != nil {
		if os.IsNotExist(err) {
			return false
		}
	}
	return true
}

func testKillProcess(t *testing.T, dfLoop *UvLoop) {
	file, err := os.OpenFile("/tmp/runforever.py", os.O_CREATE|os.O_WRONLY, 0777)
	if err != nil {
		t.Fatal(err)
	}
	fileWriter := bufio.NewWriter(file)
	_, err = fileWriter.WriteString(`
import sys
path = "/tmp"
			
errorLog = open(path + "/stderr3.txt", "w", 1)
errorLog.write("---Starting Error Log---\n")
sys.stderr = errorLog
stdoutLog = open(path + "/stdout3.txt", "w", 1)
stdoutLog.write("---Starting Standard Out Log---\n")
sys.stdout = stdoutLog
	
a = 1
while True:
	print a
	a = -a`)
	if err != nil {
		t.Fatal(err)
	}
	fileWriter.Flush()
	file.Close()
	defer os.Remove("/tmp/runforever.py")
	defer os.Remove("/tmp/stderr3.txt")
	defer os.Remove("/tmp/stdout3.txt")

	// spawn new process
	process, err := UvSpawnProcess(dfLoop, &UvProcessOptions{
		Args:  []string{"python", "/tmp/runforever.py"},
		Cwd:   "/tmp",
		Flags: UV_PROCESS_DETACHED,
		File:  "python",
		Env:   []string{"PATH"},
		ExitCb: func(h *Handle, status, sigNum int) {
			fmt.Printf("Process exited with status %d and signal %d\n", status, sigNum)
			fmt.Printf("%p\n", h.ptr.(*UvProcess))
		},
	}, nil)
	if err != nil {
		t.Fatal(err)
	}

	time.Sleep(1 * time.Second)

	// Unref this process
	process.Unref()

	// Try to kill this proces
	process.Kill(9)

	fmt.Printf("%p\n", process)
}
