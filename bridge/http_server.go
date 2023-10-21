package bridge

import (
	"fmt"
	"net/http"
)

type MetricsServer struct {
	metrics *Metrics
}

func NewMetricsServer(metrics *Metrics) *MetricsServer {
	return &MetricsServer{metrics: metrics}
}

func (ms *MetricsServer) Run(port int) error {
	http.HandleFunc("/metrics", ms.serveMetrics)
	http.HandleFunc("/-/health", ms.healthCheck)
	http.HandleFunc("/-/purge", ms.purge)
	return http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
}

func (ms *MetricsServer) serveMetrics(rw http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodGet {
		ms.writeNotImplemented(rw)
		return
	}
	b := writeAll(ms.metrics.Get())
	rw.Header().Add("content-type", "text/plain")
	rw.WriteHeader(200)
	rw.Write(b)
}

func (ms *MetricsServer) healthCheck(rw http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodGet {
		ms.writeNotImplemented(rw)
		return
	}
	rw.Header().Add("content-type", "text/plain")
	rw.WriteHeader(200)
	rw.Write([]byte("OK"))
}

func (ms *MetricsServer) purge(rw http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPost {
		ms.writeNotImplemented(rw)
		return
	}
	ms.metrics.Clean()
	rw.Header().Add("content-type", "text/plain")
	rw.WriteHeader(200)
	rw.Write([]byte("OK"))
}

func (ms *MetricsServer) writeNotImplemented(rw http.ResponseWriter) {
	rw.WriteHeader(http.StatusNotImplemented)
	rw.Write([]byte("Not Implemented"))
}

func writeAll(metrics []metric) []byte {
	var s string
	for _, m := range metrics {
		s += m.String()
	}
	return []byte(s)
}
