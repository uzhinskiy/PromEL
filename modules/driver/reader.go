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
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/golang/protobuf/proto"
	"github.com/golang/snappy"
	"github.com/prometheus/prometheus/prompb"
)

func (drv *Driver) appRead(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	//w.Header().Add("X-Server", appConfig["version"]+" b."+vBuild)
	w.Header().Set("Content-Type", "application/x-protobuf")
	w.Header().Set("Content-Encoding", "snappy")

	if r.Method == "OPTIONS" {
		return
	}

	if r.Method != http.MethodPost {
		http.Error(w, "Service Unavailable", http.StatusServiceUnavailable)
		log.Println(r.RemoteAddr, "\t", r.Method, "\t", r.URL.Path, "\t", http.StatusServiceUnavailable, "\t", "Invalid request method ", "\t", r.UserAgent())
		return
	}

	compressed, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Println("Read error: ", err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	reqBuf, err := snappy.Decode(nil, compressed)
	if err != nil {
		log.Println("Decode error: ", err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var req prompb.ReadRequest
	if err := proto.Unmarshal(reqBuf, &req); err != nil {
		log.Println("Unmarshal error: ", err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	results := make([]*prompb.QueryResult, 0, len(req.Queries))
	for _, q := range req.Queries {
		sr, err := drv.esclient.Select(q)
		if err != nil {
			log.Println("Error while read from Elasticsearch: ", err.Error())
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}

		ts, err := drv.esclient.CreateTimeseries(sr)
		if err != nil {
			log.Println("Error while read from Elasticsearch: ", err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		results = append(results, &prompb.QueryResult{Timeseries: ts})
	}

	data, err := proto.Marshal(&prompb.ReadResponse{Results: results})
	if err != nil {
		log.Println(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, "%s", snappy.Encode(nil, data))

}
