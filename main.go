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
	tickerT := time.NewTicker(time.Duration(cfg.Elastic.Interval) * time.Minute)
	defer tickerT.Stop()

	for {
		log.Println("periodic check log ...")
		run()
		<-tickerT.C
	}
}

func run() {

	// 设定固定数量的goroutine处理任务
	w := work.New(cfg.Goroutine)
	tasks := cfg.Indices

	wg.Add(len(tasks))

	for i := 0; i < len(tasks); i++ {
		go func(index *conf.Indices) {
			esTask := work.NewESConn(cfg.Elastic, cfg.Alert, index)
			w.Run(esTask)
			wg.Done()
		}(tasks[i])
	}
	wg.Wait()
	w.Shutdown()
}
