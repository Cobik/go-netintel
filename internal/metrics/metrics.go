package metrics

import (
	"net/http"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var Requests = prometheus.NewCounterVec(
	prometheus.CounterOpts{Name: "requests_total", Help: "HTTP requests total"},
	[]string{"path", "method", "code"},
)

func init() { prometheus.MustRegister(Requests) }

func Handler() http.Handler { return promhttp.Handler() }
