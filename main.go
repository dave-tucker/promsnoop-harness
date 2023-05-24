package main

//	#include <stdint.h>
//	#include <sys/sdt.h>
//
//  #define SET_COUNTER(name) \
//	static void counter_##name(uint64_t arg) { \
//		DTRACE_PROBE1(harness, counter_##name, arg); \
//	}
//  SET_COUNTER(0)
//  SET_COUNTER(1)
//  SET_COUNTER(2)
//  SET_COUNTER(3)
//  SET_COUNTER(4)
//  SET_COUNTER(5)
//  SET_COUNTER(6)
//  SET_COUNTER(7)
//  SET_COUNTER(8)
//  SET_COUNTER(9)
//  SET_COUNTER(10)
//  SET_COUNTER(11)
//  SET_COUNTER(12)
//  SET_COUNTER(13)
//  SET_COUNTER(14)
//  SET_COUNTER(15)
//  SET_COUNTER(16)
//  SET_COUNTER(17)
//  SET_COUNTER(18)
//  SET_COUNTER(19)
//  SET_COUNTER(20)
//  SET_COUNTER(21)
//  SET_COUNTER(22)
//  SET_COUNTER(23)
//  SET_COUNTER(24)
//  SET_COUNTER(25)
//  SET_COUNTER(26)
//  SET_COUNTER(27)
//  SET_COUNTER(28)
//  SET_COUNTER(29)
//  SET_COUNTER(30)
//  SET_COUNTER(31)
//  SET_COUNTER(32)
//  SET_COUNTER(33)
//  SET_COUNTER(34)
//  SET_COUNTER(35)
//  SET_COUNTER(36)
//  SET_COUNTER(37)
//  SET_COUNTER(38)
//  SET_COUNTER(39)
//  SET_COUNTER(40)
//  SET_COUNTER(41)
//  SET_COUNTER(42)
//  SET_COUNTER(43)
//  SET_COUNTER(44)
//  SET_COUNTER(45)
//  SET_COUNTER(46)
//  SET_COUNTER(47)
//  SET_COUNTER(48)
//  SET_COUNTER(49)
//  SET_COUNTER(50)
//  SET_COUNTER(51)
//  SET_COUNTER(52)
//  SET_COUNTER(53)
//  SET_COUNTER(54)
//  SET_COUNTER(55)
//  SET_COUNTER(56)
//  SET_COUNTER(57)
//  SET_COUNTER(58)
//  SET_COUNTER(59)
//  SET_COUNTER(60)
//  SET_COUNTER(61)
//  SET_COUNTER(62)
//  SET_COUNTER(63)
//  SET_COUNTER(64)
//  SET_COUNTER(65)
//  SET_COUNTER(66)
//  SET_COUNTER(67)
//  SET_COUNTER(68)
//  SET_COUNTER(69)
//  SET_COUNTER(70)
//  SET_COUNTER(71)
//  SET_COUNTER(72)
//  SET_COUNTER(73)
//  SET_COUNTER(74)
//  SET_COUNTER(75)
//  SET_COUNTER(76)
//  SET_COUNTER(77)
//  SET_COUNTER(78)
//  SET_COUNTER(79)
//  SET_COUNTER(80)
//  SET_COUNTER(81)
//  SET_COUNTER(82)
//  SET_COUNTER(83)
//  SET_COUNTER(84)
//  SET_COUNTER(85)
//  SET_COUNTER(86)
//  SET_COUNTER(87)
//  SET_COUNTER(88)
//  SET_COUNTER(89)
//  SET_COUNTER(90)
//  SET_COUNTER(91)
//  SET_COUNTER(92)
//  SET_COUNTER(93)
//  SET_COUNTER(94)
//  SET_COUNTER(95)
//  SET_COUNTER(96)
//  SET_COUNTER(97)
//  SET_COUNTER(98)
//  SET_COUNTER(99)
import "C"

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"runtime/pprof"
	"sync/atomic"
	"syscall"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var cpuprofile = flag.String("cpuprofile", "", "write cpu profile to file")

type metrics struct {
	counters []prometheus.Counter
}

func NewMetrics(reg prometheus.Registerer) *metrics {
	var counters []prometheus.Counter
	for i := 0; i < 100; i++ {
		counter := prometheus.NewCounter(prometheus.CounterOpts{
			Name: fmt.Sprintf("counter_%d", i),
			Help: fmt.Sprintf("counter %d", i),
		})
		reg.MustRegister(counter)
		counters = append(counters, counter)
	}
	return &metrics{counters: counters}
}

var usdtCounters [100]uint64

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

	reg := prometheus.NewRegistry()
	m := NewMetrics(reg)

	r.GET("/baseline", func(c *gin.Context) {
		c.JSON(200, map[string]string{
			"message": "baseline",
		})
	})
	r.GET("/metrics", func(ctx *gin.Context) {
		promhttp.HandlerFor(reg, promhttp.HandlerOpts{Registry: reg}).ServeHTTP(ctx.Writer, ctx.Request)
	})

	prom := r.Group("/prom")
	prom.GET("/1", func(c *gin.Context) {
		m.counters[0].Inc()
		c.JSON(200, map[string]string{
			"message": "1 counter incremented",
		})
	})
	prom.GET("/10", func(c *gin.Context) {
		for i := 0; i < 10; i++ {
			m.counters[i].Inc()
		}
		c.JSON(200, map[string]string{
			"message": "10 counters incremented",
		})
	})
	prom.GET("/100", func(c *gin.Context) {
		for i := 0; i < 100; i++ {
			m.counters[i].Inc()
		}
		c.JSON(200, map[string]string{
			"message": "100 counters incremented",
		})
	})

	usdt := r.Group("/usdt")
	usdt.GET("/1", func(c *gin.Context) {
		atomic.AddUint64(&usdtCounters[0], 1)
		C.counter_0(C.ulong(usdtCounters[0]))
		c.JSON(200, map[string]string{
			"message": "1 counter incremented",
		})
	})
	usdt.GET("/10", func(c *gin.Context) {
		for i := 0; i < 10; i++ {
			atomic.AddUint64(&usdtCounters[i], 1)
		}
		C.counter_0(C.ulong(usdtCounters[0]))
		C.counter_1(C.ulong(usdtCounters[1]))
		C.counter_2(C.ulong(usdtCounters[2]))
		C.counter_3(C.ulong(usdtCounters[3]))
		C.counter_4(C.ulong(usdtCounters[4]))
		C.counter_5(C.ulong(usdtCounters[5]))
		C.counter_6(C.ulong(usdtCounters[6]))
		C.counter_7(C.ulong(usdtCounters[7]))
		C.counter_8(C.ulong(usdtCounters[8]))
		C.counter_9(C.ulong(usdtCounters[9]))
		c.JSON(200, map[string]string{
			"message": "10 counters incremented",
		})
	})
	usdt.GET("/100", func(c *gin.Context) {
		for i := 0; i < 100; i++ {
			atomic.AddUint64(&usdtCounters[i], 1)
		}
		C.counter_0(C.ulong(usdtCounters[0]))
		C.counter_1(C.ulong(usdtCounters[1]))
		C.counter_2(C.ulong(usdtCounters[2]))
		C.counter_3(C.ulong(usdtCounters[3]))
		C.counter_4(C.ulong(usdtCounters[4]))
		C.counter_5(C.ulong(usdtCounters[5]))
		C.counter_6(C.ulong(usdtCounters[6]))
		C.counter_7(C.ulong(usdtCounters[7]))
		C.counter_8(C.ulong(usdtCounters[8]))
		C.counter_9(C.ulong(usdtCounters[9]))
		C.counter_10(C.ulong(usdtCounters[10]))
		C.counter_11(C.ulong(usdtCounters[11]))
		C.counter_12(C.ulong(usdtCounters[12]))
		C.counter_13(C.ulong(usdtCounters[13]))
		C.counter_14(C.ulong(usdtCounters[14]))
		C.counter_15(C.ulong(usdtCounters[15]))
		C.counter_16(C.ulong(usdtCounters[16]))
		C.counter_17(C.ulong(usdtCounters[17]))
		C.counter_18(C.ulong(usdtCounters[18]))
		C.counter_19(C.ulong(usdtCounters[19]))
		C.counter_20(C.ulong(usdtCounters[20]))
		C.counter_21(C.ulong(usdtCounters[21]))
		C.counter_22(C.ulong(usdtCounters[22]))
		C.counter_23(C.ulong(usdtCounters[23]))
		C.counter_24(C.ulong(usdtCounters[24]))
		C.counter_25(C.ulong(usdtCounters[25]))
		C.counter_26(C.ulong(usdtCounters[26]))
		C.counter_27(C.ulong(usdtCounters[27]))
		C.counter_28(C.ulong(usdtCounters[28]))
		C.counter_29(C.ulong(usdtCounters[29]))
		C.counter_30(C.ulong(usdtCounters[30]))
		C.counter_31(C.ulong(usdtCounters[31]))
		C.counter_32(C.ulong(usdtCounters[32]))
		C.counter_33(C.ulong(usdtCounters[33]))
		C.counter_34(C.ulong(usdtCounters[34]))
		C.counter_35(C.ulong(usdtCounters[35]))
		C.counter_36(C.ulong(usdtCounters[36]))
		C.counter_37(C.ulong(usdtCounters[37]))
		C.counter_38(C.ulong(usdtCounters[38]))
		C.counter_39(C.ulong(usdtCounters[39]))
		C.counter_40(C.ulong(usdtCounters[40]))
		C.counter_41(C.ulong(usdtCounters[41]))
		C.counter_42(C.ulong(usdtCounters[42]))
		C.counter_43(C.ulong(usdtCounters[43]))
		C.counter_44(C.ulong(usdtCounters[44]))
		C.counter_45(C.ulong(usdtCounters[45]))
		C.counter_46(C.ulong(usdtCounters[46]))
		C.counter_47(C.ulong(usdtCounters[47]))
		C.counter_48(C.ulong(usdtCounters[48]))
		C.counter_49(C.ulong(usdtCounters[49]))
		C.counter_50(C.ulong(usdtCounters[50]))
		C.counter_51(C.ulong(usdtCounters[51]))
		C.counter_52(C.ulong(usdtCounters[52]))
		C.counter_53(C.ulong(usdtCounters[53]))
		C.counter_54(C.ulong(usdtCounters[54]))
		C.counter_55(C.ulong(usdtCounters[55]))
		C.counter_56(C.ulong(usdtCounters[56]))
		C.counter_57(C.ulong(usdtCounters[57]))
		C.counter_58(C.ulong(usdtCounters[58]))
		C.counter_59(C.ulong(usdtCounters[59]))
		C.counter_60(C.ulong(usdtCounters[60]))
		C.counter_61(C.ulong(usdtCounters[61]))
		C.counter_62(C.ulong(usdtCounters[62]))
		C.counter_63(C.ulong(usdtCounters[63]))
		C.counter_64(C.ulong(usdtCounters[64]))
		C.counter_65(C.ulong(usdtCounters[65]))
		C.counter_66(C.ulong(usdtCounters[66]))
		C.counter_67(C.ulong(usdtCounters[67]))
		C.counter_68(C.ulong(usdtCounters[68]))
		C.counter_69(C.ulong(usdtCounters[69]))
		C.counter_70(C.ulong(usdtCounters[70]))
		C.counter_71(C.ulong(usdtCounters[71]))
		C.counter_72(C.ulong(usdtCounters[72]))
		C.counter_73(C.ulong(usdtCounters[73]))
		C.counter_74(C.ulong(usdtCounters[74]))
		C.counter_75(C.ulong(usdtCounters[75]))
		C.counter_76(C.ulong(usdtCounters[76]))
		C.counter_77(C.ulong(usdtCounters[77]))
		C.counter_78(C.ulong(usdtCounters[78]))
		C.counter_79(C.ulong(usdtCounters[79]))
		C.counter_80(C.ulong(usdtCounters[80]))
		C.counter_81(C.ulong(usdtCounters[81]))
		C.counter_82(C.ulong(usdtCounters[82]))
		C.counter_83(C.ulong(usdtCounters[83]))
		C.counter_84(C.ulong(usdtCounters[84]))
		C.counter_85(C.ulong(usdtCounters[85]))
		C.counter_86(C.ulong(usdtCounters[86]))
		C.counter_87(C.ulong(usdtCounters[87]))
		C.counter_88(C.ulong(usdtCounters[88]))
		C.counter_89(C.ulong(usdtCounters[89]))
		C.counter_90(C.ulong(usdtCounters[90]))
		C.counter_91(C.ulong(usdtCounters[91]))
		C.counter_92(C.ulong(usdtCounters[92]))
		C.counter_93(C.ulong(usdtCounters[93]))
		C.counter_94(C.ulong(usdtCounters[94]))
		C.counter_95(C.ulong(usdtCounters[95]))
		C.counter_96(C.ulong(usdtCounters[96]))
		C.counter_97(C.ulong(usdtCounters[97]))
		C.counter_98(C.ulong(usdtCounters[98]))
		C.counter_99(C.ulong(usdtCounters[99]))
		c.JSON(200, map[string]string{
			"message": "100 counters incremented",
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
