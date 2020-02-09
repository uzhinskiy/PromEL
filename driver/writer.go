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
	"io/ioutil"
	"log"
	"net/http"

	"github.com/golang/protobuf/proto"
	"github.com/golang/snappy"
	"github.com/prometheus/prometheus/prompb"
)

func (drv *Driver) appWrite(w http.ResponseWriter, r *http.Request) {
	var err error
	var req prompb.WriteRequest

	defer r.Body.Close()

	//w.Header().Add("Server", appConfig["version"]+" b."+vBuild)

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
		log.Println("Read error: ", err.Error(), "Content length is: ", r.Header["Content-Length"])
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	reqBuf, err := snappy.Decode(nil, compressed)
	if err != nil {
		log.Println("Decode error: ", err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := proto.Unmarshal(reqBuf, &req); err != nil {
		log.Println("Unmarshal error: ", err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Add("Access-Control-Allow-Origin", "*")
	w.Header().Add("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE") //!!!
	w.Header().Add("Access-Control-Allow-Headers", "Accept, Content-Type, Authorization")
	w.Header().Add("Connection", "Close")
	w.WriteHeader(http.StatusOK)

	err = drv.esclient.Insert(req)
	if err != nil {
		log.Printf("Indexing into Elastic is failed with: %s\n", err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

}
