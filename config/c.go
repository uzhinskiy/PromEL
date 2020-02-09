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

package config

import (
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

type Config struct {
	App struct {
		Debug bool `yaml:"debug,omitempty"`
	} `yaml:"app"`
	Metric struct {
		Bind string `yaml:"bind"`
		Port string `yaml:"port"`
	} `yaml:"metric,omitempty"`
	Input struct {
		Bind string `yaml:"bind"`
		Port string `yaml:"port"`
	} `yaml:"input"`
	Output struct {
		Hosts    []string `yaml:"hosts"`
		Index    string   `yaml:"index"`
		Replicas int      `yaml:"replicas"`
		Shards   int      `yaml:"shards"`
		Bulk     struct {
			Size    int    `yaml:"size"`
			Name    string `yaml:"name"`
			Workers int    `yaml:"workers"`
		} `yaml:"bulk"`
		Ilm struct {
			Hot  string `yaml:"hot"`
			Warm string `yaml:"warm"`
			Cold string `yaml:"cold"`
		} `yaml:"ilm"`
	} `yaml:"elastic"`
	Logging struct {
		Enable bool   `yaml:"enable,omitempty"`
		Path   string `yaml:"path"`
		Size   int    `yaml:"size"`
	} `yaml:"logging"`
}

func Parse(f string) Config {
	var c Config
	yamlFile, err := ioutil.ReadFile(f)
	if err != nil {
		panic(err)
	}

	err = yaml.Unmarshal(yamlFile, &c)
	if err != nil {
		panic(err)
	}

	if c.Input.Port == "" {
		c.Input.Port = "9999"
	}

	if c.Input.Bind == "" {
		c.Input.Bind = "0.0.0.0"
	}

	if c.Output.Index == "" {
		c.Output.Index = "promel"
	}

	if len(c.Output.Hosts) == 0 {
		c.Output.Hosts[0] = "http://127.0.0.1:9200/"
	}

	if c.Output.Replicas == 0 {
		c.Output.Replicas = len(c.Output.Hosts) - 1
	}

	if c.Output.Shards == 0 {
		c.Output.Shards = 4
	}

	if c.Output.Bulk.Size == 0 {
		c.Output.Bulk.Size = 1000
	}
	if c.Output.Bulk.Name == "" {
		c.Output.Bulk.Name = "promelworker"
	}
	if c.Output.Bulk.Workers == 0 {
		c.Output.Bulk.Workers = 1
	}

	if c.Output.Ilm.Hot == "" {
		c.Output.Ilm.Hot = "12h"
	}

	if c.Output.Ilm.Warm == "" {
		c.Output.Ilm.Warm = "3d"
	}

	if c.Output.Ilm.Cold == "" {
		c.Output.Ilm.Cold = "30d"
	}

	return c
}
