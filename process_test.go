package gouv

import (
	"fmt"
	"os"
	"testing"
	"time"
)

func testSpawnChildProcess(t *testing.T, dfLoop *UvLoop) {
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
	// spawn new process
	process, err := UvSpawnProcess(dfLoop, &UvProcessOptions{
		Args:  []string{"sleep", "100000"},
		Cwd:   "/tmp",
		Flags: UV_PROCESS_DETACHED,
		File:  "sleep",
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
