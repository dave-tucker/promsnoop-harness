package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"runtime/pprof"
	"syscall"

	"github.com/gin-gonic/gin"
	"github.com/penglongli/gin-metrics/ginmetrics"
)

var cpuprofile = flag.String("cpuprofile", "", "write cpu profile to file")

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

	appRouter := gin.Default()
	metricRouter := gin.Default()

	m := ginmetrics.GetMonitor()
	// use metric middleware without expose metric path
	m.UseWithoutExposingEndpoint(appRouter)
	// set metric path expose to metric router
	m.Expose(metricRouter)

	appRouter.GET("/widgets/:id", func(ctx *gin.Context) {
		ctx.JSON(200, map[string]string{
			"widgetId": ctx.Param("id"),
		})
	})
	go func() {
		_ = metricRouter.Run(":9090")
	}()
	go func() {
		_ = appRouter.Run(":8080")
	}()

	c := make(chan os.Signal, 2)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM) // subscribe to system signals
	<-c
	fmt.Println("exiting...")
	pprof.StopCPUProfile()
	f.Close()
	os.Exit(0)
}
