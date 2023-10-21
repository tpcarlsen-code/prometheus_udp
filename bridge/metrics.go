package bridge

import (
	"fmt"
	"sync"
)

type metric struct {
	Key string
	Val string
}

func (m metric) String() string {
	return fmt.Sprintf("%s %s\n", m.Key, m.Val)
}

type Metrics struct {
	metrics []metric
	index   map[string]int
	sync.RWMutex
}

func NewMetrics() *Metrics {
	return &Metrics{
		metrics: make([]metric, 0, 500),
		index:   make(map[string]int, 500),
	}
}

func (m *Metrics) Set(b []byte) {
	m.Lock()
	defer m.Unlock()
	key := make([]byte, 0, 100)
	val := make([]byte, 0, 20)
	var inVal bool
	for i := 0; i < len(b); i++ {
		if b[i] == 10 && len(key) > 0 && len(val) > 0 { // LF
			m.addMetric(metric{Key: string(key), Val: string(val)})
			key = make([]byte, 0, 100)
			val = make([]byte, 0, 20)
			inVal = false
			continue
		}
		if b[i] == 32 { // space
			inVal = true
			continue
		}
		if inVal {
			val = append(val, b[i])
			continue
		}
		key = append(key, b[i])
	}
	if len(key) > 0 && len(val) > 0 {
		m.addMetric(metric{Key: string(key), Val: string(val)})
	}
}

func (m *Metrics) addMetric(m2 metric) {
	if idx, ok := m.index[m2.Key]; ok {
		m.metrics[idx] = m2
		return
	}
	m.metrics = append(m.metrics, m2)
	m.index[m2.Key] = len(m.metrics) - 1
}

func (m *Metrics) Get() []metric {
	m.RLock()
	defer m.RUnlock()
	r := make([]metric, len(m.metrics))
	copy(r, m.metrics)
	return r
}

func (m *Metrics) Bytes() []byte {
	var out []byte
	for _, m := range m.Get() {
		out = append(out, []byte(m.String())...)
	}
	return out
}

func (m *Metrics) Clean() {
	m.Lock()
	defer m.Unlock()
	m.metrics = make([]metric, 0, 500)
	m.index = make(map[string]int, 500)
}
