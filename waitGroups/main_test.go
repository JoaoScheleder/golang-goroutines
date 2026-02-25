package main

import (
	"io"
	"os"
	"strings"
	"sync"
	"testing"
)

func TestPrintSomething(t *testing.T) {
	stdOut := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	var wg sync.WaitGroup
	wg.Add(1)
	go printSomething("Test", &wg)
	wg.Wait()

	w.Close()
	out, _ := io.ReadAll(r)
	output := string(out)

	os.Stdout = stdOut

	if !strings.Contains(output, "Test") {
		t.Errorf("Expected output to contain 'Test', got '%s'", output)
	}
}
