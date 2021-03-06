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

package driver

import (
	"log"
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/uzhinskiy/PromEL/modules/config"
	"github.com/uzhinskiy/PromEL/modules/es"
)

type Driver struct {
	esclient *es.ESClient
	conf     config.Config
}

var (
	write_bad_request = promauto.NewCounter(prometheus.CounterOpts{
		Name: "promel_write_bad_request_total",
		Help: "The total number of bad requests for write",
	})
	write_snappy_corrupted = promauto.NewCounter(prometheus.CounterOpts{
		Name: "promel_write_snappy_corrupted_total",
		Help: "The total number of corrupted snappy packets",
	})
	write_protobuf_invalid = promauto.NewCounter(prometheus.CounterOpts{
		Name: "promel_write_protobuf_invalid_total",
		Help: "The total number of invalid protobuf packets",
	})
	write_elastic_failed = promauto.NewCounter(prometheus.CounterOpts{
		Name: "promel_write_elastic_failed_total",
		Help: "The total number of failed requests to Elastic",
	})
)

func Run(cnf config.Config) error {
	esc, err := es.NewESClient(cnf)
	if err != nil {
		log.Fatal("Bootstrap: ElasticSearch is not ready, cannot start: ", err)
		return err
	} else {
		log.Println("Bootstrap: ElasticSearch is ready")
	}
	defer esc.Stop()

	if err = esc.NewBulkService(); err != nil {
		log.Fatal("Bootstrap: ", err)
		return err
	} else {
		log.Println("Bootstrap: BulkProcessor is ready")

	}
	defer esc.Close()

	go esc.Statistics()

	drv := Driver{}
	drv.esclient = esc
	drv.conf = cnf

	r1 := http.NewServeMux()
	r1.HandleFunc("/write", drv.appWrite)
	r1.HandleFunc("/read", drv.appRead)

	r2 := http.NewServeMux()
	r2.Handle("/metrics", promhttp.Handler())
	go func() {
		log.Fatal("Bootstrap: external http listener failed with - ", http.ListenAndServe(cnf.Input.Bind+":"+cnf.Input.Port, r1))
	}()
	go func() {
		log.Fatal("Bootstrap: internal http listener failed with - ", http.ListenAndServe(cnf.Metric.Bind+":"+cnf.Metric.Port, r2))
	}()
	select {}

}

func limitNumClients(f http.HandlerFunc, maxClients int) http.HandlerFunc {
	sema := make(chan struct{}, maxClients)

	return func(w http.ResponseWriter, req *http.Request) {
		sema <- struct{}{}
		defer func() { <-sema }()
		f(w, req)
	}
}
