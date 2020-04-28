package main

import (
	"flag"
	"log"
	"sync"
	"time"

	"es-index-monitor/work"

	"es-index-monitor/conf"
)

var (
	c  string
	wg sync.WaitGroup
)

var cfg *conf.Config

func init() {
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
	flag.StringVar(&c, "c", "cfg.yaml", "configuration file, default cfg.yaml")
}

func main() {
	flag.Parse()
	cfg = conf.LoadFile(c)
	tickerT := time.NewTicker(time.Duration(cfg.Ticker) * time.Minute)
	defer tickerT.Stop()

	for {
		log.Println("check elastic")
		run()
		<-tickerT.C
	}
}

func run() {

	// consumer
	w := work.New(cfg.Goroutine)

	// concurrent tasks
	tasks := cfg.Indices
	// wait goroutine task to finish
	wg.Add(len(tasks))
	for i := 0; i < len(tasks); i++ {
		go func(index *conf.Indices) {
			// productions tasks
			task := work.NewConn(cfg.Elastic, cfg.Alert, index)
			// when the consumer is not enough, the task will block waiting.
			w.Run(task)
			wg.Done()
		}(tasks[i])
	}
	wg.Wait()

	// destory consumers at the end of all tasks.
	w.Shutdown()
}
