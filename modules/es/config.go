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
	cnf "github.com/uzhinskiy/PromEL/modules/config"
)

type esconf struct {
	Hosts    []string
	SSL      bool
	Cert     string
	Index    string
	Replicas int
	Shards   int
	Bulk     struct {
		Name    string
		Size    int
		Workers int
		Flush   int
	}
}

func riseconfig(in cnf.Config) esconf {
	c := esconf{}

	if in.Output.Index == "" {
		c.Index = "promel"
	} else {
		c.Index = in.Output.Index
	}

	if !in.Output.SSL {
		c.SSL = false
	} else {
		c.SSL = true
		c.Cert = in.Output.Cert
	}

	if len(in.Output.Hosts) == 0 {
		c.Hosts[0] = "http://127.0.0.1:9200/"
	} else {
		c.Hosts = in.Output.Hosts
	}

	if in.Output.Replicas == 0 {
		c.Replicas = len(in.Output.Hosts) - 1
	} else {
		c.Replicas = in.Output.Replicas
	}

	if in.Output.Shards == 0 {
		c.Shards = 4
	} else {
		c.Shards = in.Output.Shards
	}

	if in.Output.Bulk.Size == 0 {
		c.Bulk.Size = 1000
	} else {
		c.Bulk.Size = in.Output.Bulk.Size
	}

	if in.Output.Bulk.Flush == 0 {
		c.Bulk.Flush = 5
	} else {
		c.Bulk.Flush = in.Output.Bulk.Flush
	}

	if in.Output.Bulk.Name == "" {
		c.Bulk.Name = "promelworker"
	} else {
		c.Bulk.Name = in.Output.Bulk.Name
	}

	if in.Output.Bulk.Workers == 0 {
		c.Bulk.Workers = 1
	} else {
		c.Bulk.Workers = in.Output.Bulk.Workers
	}

	return c
}
