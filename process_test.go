package gouv

import (
	"bufio"
	"fmt"
	"os"
	"testing"
	"time"
)

func TestSpawnChildProcess(t *testing.T) {
	dfLoop := UvLoopDefault()

	// spawn new process
	process, err := UvSpawnProcess(dfLoop, &UvProcessOptions{
		Args:  []string{"ls", "-lah"},
		Cwd:   "/tmp",
		Flags: UV_PROCESS_DETACHED,
		File:  "ls",
		ExitCb: func(h *Handle, status, sigNum int) {
			fmt.Printf("Process exited with status %d and signal %d\n", status, sigNum)
		},
	}, nil)
	if err != nil {
		t.Fatal(err)
	}

	go dfLoop.Run(UV_RUN_DEFAULT)
	time.Sleep(1 * time.Second)

	// Unref this process
	process.Unref()

	go dfLoop.Close()
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

func TestKillProcess(t *testing.T) {
	file, err := os.OpenFile("/tmp/runforever.py", os.O_CREATE|os.O_WRONLY, 0777)
	if err != nil {
		t.Fatal(err)
	}
	fileWriter := bufio.NewWriter(file)
	_, err = fileWriter.WriteString(`
import sys
path = "/tmp"
		
errorLog = open(path + "/stderr.txt", "w", 1)
errorLog.write("---Starting Error Log---\n")
sys.stderr = errorLog
stdoutLog = open(path + "/stdout.txt", "w", 1)
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
	defer os.Remove("/tmp/stderr.txt")
	defer os.Remove("/tmp/stdout.txt")

	dfLoop := UvLoopDefault()

	// spawn new process
	process, err := UvSpawnProcess(dfLoop, &UvProcessOptions{
		Args:  []string{"python", "/tmp/runforever.py"},
		Cwd:   "/tmp",
		Flags: UV_PROCESS_DETACHED,
		File:  "python",
		Env:   []string{"PATH"},
		ExitCb: func(h *Handle, status, sigNum int) {
			if sigNum != 9 {
				t.Fatal("Kill failed")
			}

			if !fileExists("/tmp/stderr.txt") {
				t.Fatalf("Process stderr not valid")
			}

			if !fileExists("/tmp/stdout.txt") {
				t.Fatalf("Process stdout not valid")
			}

			fmt.Printf("Process exited with status %d and signal %d\n", status, sigNum)
		},
	}, nil)
	if err != nil {
		t.Fatal(err)
	}

	go dfLoop.Run(UV_RUN_DEFAULT)
	time.Sleep(3 * time.Second)

	// Unref this process
	process.Unref()

	// Try to kill this proces
	process.Kill(9)

	time.Sleep(2 * time.Second)

	go dfLoop.Close()
}
