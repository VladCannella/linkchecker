// Command linkcheck is a concurrent link checker.
//
// Flags (Этап 1-9):
//
//	-workers   int           размер пула воркеров (по умолчанию runtime.NumCPU())
//	-rps       int           запросов в секунду (rate limiter, Этап 3)
//	-timeout   duration      общий дедлайн на весь прогон (Этап 4)
//	-per-req   duration      таймаут одного запроса (Этап 4)
//	-out       string        файл отчёта в формате JSON
package main

import (
	"flag"
	"fmt"
	"os"
)

func main() {
	workers := flag.Int("workers", 0, "worker pool size (default: runtime.NumCPU())")
	rps := flag.Int("rps", 0, "requests per second (0 = unlimited)")
	timeout := flag.Duration("timeout", 0, "overall deadline for the whole run")
	perReq := flag.Duration("per-req", 0, "timeout for a single request")
	out := flag.String("out", "", "output JSON report file (empty = stdout)")
	flag.Parse()

	urls := flag.Args()
	if len(urls) == 0 {
		fmt.Fprintln(os.Stderr, "usage: linkcheck [flags] url [url...]")
		os.Exit(2)
	}

	_ = workers
	_ = rps
	_ = timeout
	_ = perReq
	_ = out

	// TODO (Этап 1): горутина на каждый URL, sync.WaitGroup, канал results.
	// TODO (Этап 2): заменить на worker pool (fan-out/fan-in) с флагом -workers.
	// TODO (Этап 3): семафор на хост + rate limiter на тикере (-rps).
	// TODO (Этап 4): context.WithTimeout на весь прогон и на каждый запрос.
	// TODO (Этап 9): signal.NotifyContext для graceful shutdown по Ctrl+C.
	// TODO: собрать []checker.Result и записать JSON-отчёт в -out (или stdout).
}
