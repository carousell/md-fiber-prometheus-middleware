package fiberprom

import (
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/gofiber/adaptor/v2"
	"github.com/gofiber/fiber/v2"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var defaultMetricPath = "/metrics"

// Prometheus contains the metrics gathered by the instance and its path
type Prometheus struct {
	reqDur      *prometheus.HistogramVec
	router      fiber.Router
	MetricsPath string
	// urlMapper is a map of url to be mapped to a different url to avoid too many labels
	urlMapper map[string]string
}

// NewPrometheus generates a new set of metrics with a certain subsystem name
func NewPrometheus(subsystem string) *Prometheus {
	p := &Prometheus{
		MetricsPath: defaultMetricPath,
	}
	p.registerMetrics(subsystem)

	return p
}

func (p *Prometheus) SetURLMapper(mapper map[string]string) {
	p.urlMapper = mapper
}

func (p *Prometheus) registerMetrics(subsystem string) {
	p.reqDur = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Subsystem: subsystem,
			Name:      "request_duration_seconds",
			Help:      "request latencies",
			Buckets:   []float64{.005, .01, .02, 0.04, .06, 0.08, .1, 0.15, .25, 0.4, .6, .8, 1, 1.5, 2, 3, 5},
		},
		[]string{"code", "path"},
	)

	err := prometheus.Register(p.reqDur)
	if err != nil {
		log.Printf("failed to register metrics: %v", err)
	}
}

// Use adds the middleware to a fiber
func (p *Prometheus) Use(r fiber.Router) {
	r.Use(p.HandlerFunc())
	r.Get(p.MetricsPath, prometheusHandler())
}

// HandlerFunc is onion or wrapper to handler for fasthttp listenandserve
func (p *Prometheus) HandlerFunc() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		uri := string(ctx.Request().URI().Path())
		if uri == p.MetricsPath {
			// next
			return ctx.Next()
		}

		if len(p.urlMapper) > 0 {
			for prefix, val := range p.urlMapper {
				if strings.HasPrefix(uri, prefix) {
					uri = val
					break
				}
			}
		}

		start := time.Now()

		defer func() {
			status := strconv.Itoa(ctx.Response().StatusCode())
			elapsed := float64(time.Since(start)) / float64(time.Second)
			ep := ctx.Method() + "_" + uri
			p.reqDur.WithLabelValues(status, ep).Observe(elapsed)
		}()
		// next
		return ctx.Next()
	}
}

// since prometheus/client_golang use net/http we need this net/http adapter for fiber
func prometheusHandler() fiber.Handler {
	return adaptor.HTTPHandler(promhttp.Handler())
}
