package work

import (
	"context"
	"fmt"
	"io"
	"log"
	"strconv"
	"strings"
	"time"

	"es-index-monitor/alarm"

	"es-index-monitor/conf"

	"gopkg.in/olivere/elastic.v5"
)

type ES struct {
	conn      *elastic.Client
	index     *conf.Indices
	interval  int
	overcount int
	alert     *alarm.Alert
}

func NewConn(es *conf.Elastic, alert *conf.Alert, index *conf.Indices) *ES {
	conn, err := elastic.NewClient(elastic.SetSniff(false),
		elastic.SetBasicAuth(es.User, es.Password),
		elastic.SetURL(es.URL...),
		elastic.SetMaxRetries(3))
	if err != nil {
		log.Println(err)
	}

	return &ES{
		conn:  conn,
		index: index,
		alert: alarm.New(alert.URL, alert.From, alert.Level),
	}
}

func (e *ES) Task() {
	start, end := lastInterval(e.index.Interval)
	index := fmt.Sprintf("%v-%v", e.index.Name, end.Format("2006.01.02"))
	e.search(start, end, index)
}

func lastInterval(i int) (time.Time, time.Time) {
	m, _ := time.ParseDuration(fmt.Sprintf("-%vm", i))
	endT := time.Now()
	startT := endT.Add(m)
	return startT, endT
}

func (e *ES) search(start time.Time, end time.Time, index string) {

	// set the query time range according to the field
	rangeQ := elastic.NewRangeQuery("@timestamp").Gte(start).Lte(end)

	switch e.index.Type {
	case "eq":
		//  one or more key-value queries
		query := elastic.NewBoolQuery().Must(rangeQ)
		pairs := strings.Split(e.index.Include, ";")
		for _, v := range pairs {
			tmp := strings.Split(v, "=")
			query = query.Must(elastic.NewTermQuery(tmp[0], tmp[1]))
		}
		e.execute(index, query)

	case "gt":
		//  a key is greater than a value while excluing some key-value
		query := elastic.NewBoolQuery().Must(rangeQ)
		tmp := strings.Split(e.index.Include, ">")
		query = query.Must(elastic.NewRangeQuery(tmp[0]).Gte(tmp[1]))

		pairs := strings.Split(e.index.Exclude, ";")
		for _, v := range pairs {
			tmp := strings.Split(v, "=")
			query = query.MustNot(elastic.NewTermQuery(tmp[0], tmp[1]))
		}
		e.execute(index, query)

	case "agg":
		// aggregate query based on field, get the top size results
		tmp := strings.Split(e.index.Include, "=")
		size, err := strconv.ParseInt(tmp[1], 10, 0)
		if err != nil {
			log.Println(err)
			return
		}
		agg := elastic.NewTermsAggregation().Field(tmp[0]).Size(int(size)).OrderByCountDesc()
		results, err := e.conn.Search(index).Size(0).Query(rangeQ).Aggregation("tmps", agg).Do(context.Background())
		if err != nil {
			log.Println(err)
			return
		}
		e.handleAggResults(index, results)
	}
}

func (e *ES) execute(index string, query *elastic.BoolQuery) {
	scroll := elastic.NewScrollService(e.conn)
	result, err := scroll.Index(index).Query(query).KeepAlive("1m").Size(4000).Do(context.Background())
	if err != nil && err != io.EOF {
		log.Println(err)
		return
	}

	count := result.Hits.TotalHits
	if count > int64(e.index.Overcount) {
		txt := fmt.Sprintf("%v分钟内%v条件的数量为%v,超过阈值%v", e.index.Interval, e.index.Include, count, e.index.Overcount)
		e.alert.Send(fmt.Sprintf("%v (%v)", index, e.index.Include), txt)
	}
}

func (e *ES) handleAggResults(index string, aggResult *elastic.SearchResult) {
	items, ok := aggResult.Aggregations.Terms("tmps")
	if !ok {
		return
	}

	exclude := make(map[string]bool)
	pairs := strings.Split(e.index.Exclude, ";")
	for _, v := range pairs {
		tmp := strings.Split(v, "=")
		exclude[tmp[1]] = true
	}

	for _, item := range items.Buckets {
		field := fmt.Sprintf("%v", item.Key)

		if exclude[field] {
			continue
		}

		if item.DocCount > int64(e.index.Overcount) {
			txt := fmt.Sprintf("%v分钟内%v的数量为%v,超过阈值%v", e.index.Interval, field, item.DocCount, e.index.Overcount)
			e.alert.Send(fmt.Sprintf("%v (%v)", index, field), txt)
		}
	}
}
