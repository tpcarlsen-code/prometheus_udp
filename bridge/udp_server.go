package bridge

import (
	"fmt"
	"net"
	"strconv"
	"time"
)

const (
	metricReceivedMetrics = "prometheus_udp_received_metrics"
	metricReceivedBytes   = "prometheus_udp_received_bytes"
)

type NetworkReader interface {
	Read(b []byte) (int, error)
}

var readBuffer []byte
var processBuffer []byte
var metricsReceived int
var bytesReceived int

type IntakeServer struct {
	metrics       *Metrics
	networkReader NetworkReader
}

func DefaultUDPNetworkReader(port int) (NetworkReader, error) {
	addr, err := net.ResolveUDPAddr("udp4", ":"+strconv.Itoa(port))
	if err != nil {
		return nil, err
	}
	return net.ListenUDP("udp4", addr)
}

func NewIntakeServer(m *Metrics, nr NetworkReader) *IntakeServer {
	return &IntakeServer{
		metrics:       m,
		networkReader: nr,
	}
}

func (is *IntakeServer) Run() error {
	go func() {
		is.setOwnMetrics()
		time.Sleep(10 * time.Second)
	}()
	for {
		readBuffer = make([]byte, 1<<16)
		read, err := is.networkReader.Read(readBuffer)
		if err != nil {
			return err
		}
		is.readLines(readBuffer, read)
		is.metrics.Set(is.makeMetricBytes(metricReceivedBytes, bytesReceived))
		is.metrics.Set(is.makeMetricBytes(metricReceivedMetrics, metricsReceived))
	}
}

func (is *IntakeServer) readLines(b []byte, n int) {
	bytesReceived += n
	temp := make([]byte, 0, 256)

	for i := 0; i < len(processBuffer); i++ {
		temp = append(temp, processBuffer[i])
	}

	for i := 0; i < n; i++ {
		temp = append(temp, b[i])
		if b[i] == 10 { // LF
			is.metrics.Set(temp)
			temp = make([]byte, 0, 256)
			metricsReceived++
		}
	}
	processBuffer = temp
}

func (is *IntakeServer) makeMetricBytes(name string, value int) []byte {
	return []byte(fmt.Sprintf("%s %d\n", name, value))
}

func (is *IntakeServer) setOwnMetrics() {
	is.metrics.Set(is.makeMetricBytes(metricReceivedBytes, bytesReceived))
	is.metrics.Set(is.makeMetricBytes(metricReceivedMetrics, metricsReceived))
}
