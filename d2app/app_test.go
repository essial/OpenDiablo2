package d2app

import (
	"bytes"
	"io"
	"log"
	"os"
	"strings"
	"sync"
	"testing"
)

func captureOutput(f func()) string {
	reader, writer, err := os.Pipe()
	if err != nil {
		panic(err)
	}

	stdout := os.Stdout
	stderr := os.Stderr

	defer func() {
		os.Stdout = stdout
		os.Stderr = stderr
		log.SetOutput(os.Stderr)
	}()

	os.Stdout = writer
	os.Stderr = writer

	log.SetOutput(writer)

	out := make(chan string)
	wg := new(sync.WaitGroup)
	wg.Add(1)

	go func() {
		var buf bytes.Buffer

		wg.Done()
		if _, err := io.Copy(&buf, reader); err != nil {
			log.Fatalf("Could not copy output buffer!")
		}

		out <- buf.String()
	}()

	wg.Wait()
	f()

	if err := writer.Close(); err != nil {
		log.Fatalf("Could not close writer.")
	}

	return <-out
}

func TestCommandLineHelp(t *testing.T) {
	os.Args = []string{os.Args[0], "-h"}

	output := captureOutput(func() {
		_ = Create("branch", "commit")
	})

	lines := strings.Split(output, "\n")

	if len(lines) < 4 || lines[3] != "Flags:" {
		t.Errorf("Flags were not returned.")
		t.Fail()
	}
}

func TestCommandLineVersion(t *testing.T) {
	os.Args = []string{os.Args[0], "-v"}

	output := captureOutput(func() {
		_ = Create("branch", "commit")
	})

	lines := strings.Split(output, "\n")

	if len(lines) < 2 || !strings.Contains(lines[1], "OpenDiablo2 (branch commit)") {
		t.Errorf("Version was not properly formatted (or not returned at all).")
		t.Fail()
	}
}
