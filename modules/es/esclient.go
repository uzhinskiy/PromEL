// Copyright © 2020 Uzhinskiy Boris
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package es

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"math"
	"strings"
	"time"

	"github.com/olivere/elastic/v7"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/common/model"
	"github.com/prometheus/prometheus/prompb"
	cnf "github.com/uzhinskiy/PromEL/modules/config"
)

type ESClient struct {
	ec       *elastic.Client
	bps      *elastic.BulkProcessor
	ctx      context.Context
	config   esconf
	idxready bool
}

type promSample struct {
	Labels    model.Metric `json:"label"`
	Value     float64      `json:"value"`
	TimeStamp int64        `json:"timestamp"`
	Datetime  string       `json:"datetime"`
}

var (
	promel_docs_indexed_speed = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "promel_docs_indexed_speed",
		Help: "Speed of indexing documents",
	})
	promel_docs_indexed_failed_total = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "promel_docs_indexed_failed_total",
		Help: "Number of requests reported as failed",
	})
)

func NewESClient(in_cnf cnf.Config) (*ESClient, error) {
	esc := &ESClient{}
	esc.config = riseconfig(in_cnf)
	//Подключение к Эластик
	ehosts := esc.config.Hosts
	tmp, err := elastic.NewClient(
		elastic.SetURL(ehosts...),
		elastic.SetSniff(true),
		elastic.SetHealthcheckInterval(10*time.Second),
		elastic.SetRetrier(initretrier()),
		elastic.SetGzip(true),
	)
	if err != nil {
		return nil, err
	}

	if tmp.IsRunning() == true {
		esc.ec = tmp
		err = esc.SetupIndex(esc.config)
		if err != nil {
			esc.idxready = false
			return nil, errors.New("Setting up Index failed with: " + err.Error())
		}
		return esc, nil
	} else {
		return nil, err
	}
}

func (esc *ESClient) NewBulkService() error {
	p, err := esc.ec.BulkProcessor().
		//Name(esc.config.Bulk.Name).
		Workers(esc.config.Bulk.Workers).
		BulkActions(esc.config.Bulk.Size).
		FlushInterval(time.Duration(esc.config.Bulk.Flush) * time.Second). // commit every esc.config.Bulk.Flush seconds
		Stats(true).                                                       // enable collecting stats
		Do(context.Background())
	if err != nil {
		return errors.New("Setting up BulkProcessor failed with: " + err.Error())
	}
	esc.bps = p
	return nil
}

func (esc *ESClient) Stop() {
	esc.ec.Stop()
}

func (esc *ESClient) Close() error {
	return esc.bps.Close()
}

func (esc *ESClient) Statistics() {
	var t float64
	var c float64
	for {
		stats := esc.bps.Stats()

		c = float64(stats.Indexed/5) - t
		promel_docs_indexed_speed.Set(c)
		t = float64(stats.Indexed / 5)
		promel_docs_indexed_failed_total.Set(float64(stats.Failed))

		/*for i, w := range stats.Workers {
			log.Printf("Worker %d: Number of requests queued: %d\n", i, w.Queued)
			log.Printf("           Last response time       : %v\n", w.LastDuration)
		}*/
		time.Sleep(5 * time.Second)

	}
}

func (esc *ESClient) Insert(req prompb.WriteRequest) error {
	for _, ts := range req.Timeseries {
		metric := make(model.Metric, len(ts.Labels))
		for _, l := range ts.Labels {
			metric[model.LabelName(l.Name)] = model.LabelValue(l.Value)
		}
		for _, s := range ts.Samples {
			v := float64(s.Value)
			if math.IsNaN(v) || math.IsInf(v, 0) {
				continue
			}
			doc := promSample{
				metric,
				v,
				s.Timestamp,
				time.Now().Format(time.RFC3339Nano),
			}
			r := elastic.
				NewBulkIndexRequest().
				Index(esc.config.Index).
				Type("_doc").
				Doc(doc)
			esc.bps.Add(r)
		}

	}
	return nil
}

func (esc *ESClient) Select(q *prompb.Query) (*elastic.SearchHits, error) {

	query, err := buildquery(q)
	if err != nil {
		return nil, err
	}

	searchResult, err := esc.ec.Search().
		Index(esc.config.Index). // search in index "users"
		Query(query).            // specify the query
		Size(1000).              // TODO: сделать scroll
		Do(context.Background()) // execute
	if err != nil {
		return nil, err
	}

	if searchResult.TotalHits() > 0 {
		return searchResult.Hits, nil
	} else {
		return nil, errors.New("Not Found")
	}

}

func buildquery(q *prompb.Query) (*elastic.BoolQuery, error) {

	query := elastic.NewBoolQuery()

	for _, m := range q.Matchers {
		if m.Name == "__name__" {
			switch m.Type {
			case prompb.LabelMatcher_EQ:
				query = query.Must(elastic.NewTermQuery("label.__name__", escapeSlashes(m.Value)))
			case prompb.LabelMatcher_RE:
				query = query.MustNot(elastic.NewTermQuery("label.__name__", escapeSlashes(m.Value)))
			default:
				// TODO: Figure out how to support these efficiently.
				return query, errors.New("non-equal or regex-non-equal matchers are not supported on the metric name yet")
			}
			continue
		}

		switch m.Type {
		case prompb.LabelMatcher_EQ:
			query = query.Must(elastic.NewTermQuery("label."+m.Name, escapeSlashes(m.Value)))
		case prompb.LabelMatcher_NEQ:
			query = query.Must(elastic.NewTermQuery("label."+m.Name, escapeSlashes(m.Value)))
		case prompb.LabelMatcher_RE:
			query = query.Must(elastic.NewRegexpQuery("label."+m.Name, ".*["+escapeSlashes(m.Value)+"].*"))
		case prompb.LabelMatcher_NRE:
			query = query.MustNot(elastic.NewRegexpQuery("label."+m.Name, ".*["+escapeSlashes(m.Value)+"].*"))
		default:
			return query, errors.New(fmt.Sprintf("unknown match type %v", m.Type))
		}
	}
	query = query.Must(elastic.NewRangeQuery("timestamp").Gte(q.StartTimestampMs).Lte(q.EndTimestampMs))

	return query, nil
}

func (esc *ESClient) CreateTimeseries(results *elastic.SearchHits) ([]*prompb.TimeSeries, error) {
	tsMap := make(map[string]*prompb.TimeSeries)
	for _, r := range results.Hits {
		var s promSample
		if err := json.Unmarshal(r.Source, &s); err != nil {
			return nil, err
		}
		fingerprint := s.Labels.Fingerprint().String()

		ts, ok := tsMap[fingerprint]
		if !ok {
			labels := make([]*prompb.Label, 0, len(s.Labels))
			for k, v := range s.Labels {
				labels = append(labels, &prompb.Label{
					Name:  string(k),
					Value: string(v),
				})
			}
			ts = &prompb.TimeSeries{
				Labels: labels,
			}
			tsMap[fingerprint] = ts
		}
		ts.Samples = append(ts.Samples, prompb.Sample{
			Value:     s.Value,
			Timestamp: s.TimeStamp,
		})
	}
	ret := make([]*prompb.TimeSeries, 0, len(tsMap))

	for _, s := range tsMap {
		ret = append(ret, s)
	}
	return ret, nil
}

func escapeSingleQuotes(str string) string {
	return strings.Replace(str, `'`, `\'`, -1)
}

func escapeSlashes(str string) string {
	return strings.Replace(str, `/`, `\/`, -1)
}
