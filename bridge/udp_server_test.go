package bridge

import (
	_ "embed"
	"strings"
	"testing"
	"time"
)

//go:embed test.txt
var data []byte

type mockReader struct {
	next chan []byte
}

func (m mockReader) Read(b []byte) (int, error) {
	n := <-m.next
	copy(b, n)
	return len(n), nil
}

func TestIntake(t *testing.T) {
	ch := make(chan []byte)
	nr := mockReader{next: ch}
	metrics := NewMetrics()
	server := NewIntakeServer(metrics, nr)
	go server.Run()
	chunkSize := 256
	for i := 0; i < len(data); i += chunkSize {
		end := i + chunkSize
		if end > len(data) {
			end = len(data)
		}
		ch <- data[i:end]
	}
	time.Sleep(10 * time.Millisecond) // wait for processing of last package segment.

	parsedLines := strings.Split(string(metrics.Bytes()), "\n")
	inputLines := strings.Split(string(data), "\n")

	// check all our input lines are present in output.
	for _, il := range inputLines {
		found := false
		for _, pl := range parsedLines {
			if pl == il {
				found = true
				break
			}
		}
		if !found {
			t.Fatalf("expected to find %s", il)
		}
	}
}
