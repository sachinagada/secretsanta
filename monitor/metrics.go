package monitor

import (
	"fmt"
	"log"
	"net"
	"net/http"

	// pprof handlers will register with DefaultServeMux at start up
	_ "net/http/pprof"

	"github.com/sachinagada/secretsanta/pick"

	prom "github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/collectors"

	"contrib.go.opencensus.io/exporter/prometheus"

	"go.opencensus.io/plugin/ochttp"
	"go.opencensus.io/stats/view"
	"go.opencensus.io/zpages"
)

// NewServer returns an http.Server used for monitoring the application. The
// tracez debug pages, pprof profiles, and metrics endpoint are all configured
// on this server.
func NewServer(port string) (http.Server, error) {
	regErr := register()
	if regErr != nil {
		log.Fatalf("error registering opencensus views: %s", regErr)
	}

	promReg := prom.NewRegistry()
	promReg.Register(collectors.NewProcessCollector(collectors.ProcessCollectorOpts{}))
	promReg.Register(collectors.NewGoCollector())
	pe, peErr := prometheus.NewExporter(
		prometheus.Options{
			Registry: promReg,
		},
	)
	if peErr != nil {
		return http.Server{}, fmt.Errorf("error initializing prometheus exporter: %w", peErr)
	}

	http.Handle("/metrics", pe)
	zpages.Handle(http.DefaultServeMux, "/debug")

	return http.Server{
		Addr: net.JoinHostPort("", port),
	}, nil
}

func register() error {
	regErr := view.Register(pick.ViewParticipants, pick.ViewCommunicationLatency)
	if regErr != nil {
		return fmt.Errorf("error registering custom server views: %w", regErr)
	}

	if ocregErr := view.Register(
		ochttp.ServerRequestBytesView,
		ochttp.ServerResponseBytesView,
		ochttp.ServerLatencyView,
		ochttp.ServerRequestCountByMethod,
		ochttp.ServerResponseCountByStatusCode,
	); ocregErr != nil {
		return fmt.Errorf("error registering ochttp server views")
	}
	return nil
}
