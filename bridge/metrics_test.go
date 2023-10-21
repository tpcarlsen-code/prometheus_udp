package bridge_test

import (
	"fmt"
	"testing"

	"github.com/tpcarlsen-code/prometheus_udp/bridge"
)

func TestSetGet(t *testing.T) {
	metrics := bridge.NewMetrics()
	testSize := 50000
	for i := 0; i < testSize; i++ {
		m := fmt.Sprintf("metric_%d %d.%d\n", i, i, i)
		metrics.Set([]byte(m))
	}
	received := metrics.Get()
	if len(received) != testSize {
		t.Errorf("wrong length, got %d should be %d", len(received), testSize)
	}
	for i := 0; i < testSize; i++ {
		if received[i].Key != fmt.Sprintf("metric_%d", i) {
			t.Errorf("should be %s was %s", fmt.Sprintf("metric_%d", i), received[i].Key)
		}
		if received[i].Val != fmt.Sprintf("%d.%d", i, i) {
			t.Errorf("should be %s was %s", fmt.Sprintf("%d.%d", i, i), received[i].Val)
		}
	}
}

func BenchmarkSet(b *testing.B) {
	d := []byte("metric_1 89\n")
	metrics := bridge.NewMetrics()
	for i := 0; i < b.N; i++ {
		metrics.Set(d)
	}
}
