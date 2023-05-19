package main

//		#include <stdint.h>
//		#include <sys/sdt.h>
//
//      const char* provider = "harness";
//      const char* requests_total = "requests_total";
//
//		static void set_requests_total(uint64_t arg) {
//			DTRACE_PROBE1(provider, requests_total, arg);
//		}
import "C"

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"runtime/pprof"
	"syscall"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var cpuprofile = flag.String("cpuprofile", "", "write cpu profile to file")
var useProm = flag.Bool("prometheus", false, "use prometheus")
var useUsdt = flag.Bool("usdt", false, "use usdt")

type metrics struct {
	requestsTotal prometheus.Counter
}

func NewMetrics(reg prometheus.Registerer) *metrics {
	m := &metrics{
		requestsTotal: prometheus.NewCounter(prometheus.CounterOpts{
			Name: "requests_total",
			Help: "Number of requests processed by the server, regardless of success or failure",
		}),
	}
	reg.MustRegister(m.requestsTotal)
	return m
}

var (
	metricRequestsTotal uint64 = 0
)

func udstInterceptor(ctx *gin.Context) {
	//	startTime := time.Now()

	// execute normal process.
	ctx.Next()

	// after request
	metricRequestsTotal++
	C.set_requests_total(C.ulong(metricRequestsTotal))
}

func prometheusInterceptor(m *metrics) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		if ctx.Request.URL.Path == "/metrics" {
			ctx.Next()
			return
		}
		ctx.Next()
		m.requestsTotal.Inc()
	}
}

func main() {
	var f *os.File
	flag.Parse()
	if *cpuprofile != "" {
		var err error
		f, err = os.Create(*cpuprofile)
		if err != nil {
			log.Fatal(err)
		}
		pprof.StartCPUProfile(f)
	}

	gin.SetMode(gin.ReleaseMode)

	r := gin.Default()

	if *useProm {
		reg := prometheus.NewRegistry()
		m := NewMetrics(reg)
		r.Use(prometheusInterceptor(m))
		r.GET("/metrics", func(ctx *gin.Context) {
			promhttp.HandlerFor(reg, promhttp.HandlerOpts{Registry: reg}).ServeHTTP(ctx.Writer, ctx.Request)
		})
	}

	if *useUsdt {
		r.Use(udstInterceptor)
	}

	r.GET("/widgets/:id", func(ctx *gin.Context) {
		ctx.JSON(200, map[string]string{
			"widgetId": ctx.Param("id"),
		})
	})

	go func() {
		_ = r.Run(":8080")
	}()

	c := make(chan os.Signal, 2)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM) // subscribe to system signals
	<-c
	fmt.Println("exiting...")
	pprof.StopCPUProfile()
	f.Close()
	os.Exit(0)
}
