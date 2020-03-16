package work

import (
	"context"
	"fmt"
	"io"
	"log"
	"time"

	"es-index-monitor/alarm"

	"es-index-monitor/conf"

	"gopkg.in/olivere/elastic.v5"
)

type ES struct {
	conn      *elastic.Client
	indices   *conf.Indices
	interval  int
	threshold int
	alert     *alarm.Alert
}

func NewESConn(es *conf.Elastic, alert *conf.Alert, index *conf.Indices) *ES {
	conn, err := elastic.NewClient(elastic.SetSniff(false),
		elastic.SetBasicAuth(es.User, es.Password),
		elastic.SetURL(es.URL...),
		elastic.SetMaxRetries(3))
	if err != nil {
		log.Println("ERROR ", err)
	}

	return &ES{
		conn:      conn,
		indices:   index,
		interval:  es.Interval,
		threshold: es.Threshold,
		alert:     alarm.New(alert.URL, fmt.Sprintf("%v%v", index.App, alert.From), alert.Level),
	}
}

func (e *ES) Task() {
	start, end := lastInterval(e.interval)
	index := fmt.Sprintf("%v-%v-%v", e.indices.Name, e.indices.App, end.Format("2006.01.02"))
	result := e.search(start, end, index)
	if result > int64(e.threshold) {
		txt := fmt.Sprintf("%v分钟内%v级别的数量为%v,超过阈值%v", e.interval, e.indices.Lv, result, e.threshold)
		e.alert.Send(index, txt)
	}
}

func lastInterval(i int) (time.Time, time.Time) {
	m, _ := time.ParseDuration(fmt.Sprintf("-%vm", i))
	endT := time.Now()
	startT := endT.Add(m)
	return startT, endT
}

func (e *ES) search(start time.Time, end time.Time, index string) int64 {
	rangeQ := elastic.NewRangeQuery("@timestamp").Gte(start).Lte(end)
	termQ1 := elastic.NewTermQuery("lv", e.indices.Lv)
	termQ2 := elastic.NewTermQuery("fields.app", e.indices.App)
	query := elastic.NewBoolQuery().Must(rangeQ, termQ1, termQ2)
	return e.execute(index, query)
}

func (e *ES) execute(index string, query *elastic.BoolQuery) int64 {
	scroll := elastic.NewScrollService(e.conn)
	result, err := scroll.Index(index).Query(query).KeepAlive("1m").Size(4000).Do(context.Background())
	if err != nil && err != io.EOF {
		log.Println("ERROR ", err, index)
		return 0
	}
	return result.Hits.TotalHits
}
