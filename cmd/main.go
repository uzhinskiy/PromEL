// Copyright © 2020 Uzhinskiy Boris <boris.ujinsky@gmail.com>
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

package main

import (
	"flag"
	"log"
	"os"

	"github.com/uzhinskiy/PromEL/modules/config"
	"github.com/uzhinskiy/PromEL/modules/driver"
)

var (
	configfile string
	vBuild     string
	cnf        config.Config
	hostname   string
)

func init() {
	flag.StringVar(&configfile, "config", "main.yml", "Read configuration from this file")
	flag.StringVar(&configfile, "f", "main.yml", "Read configuration from this file")
	version := flag.Bool("V", false, "Show version")
	flag.Parse()
	if *version {
		print("Build num: ", vBuild, "\n")
		os.Exit(0)
	}

	hostname, _ = os.Hostname()
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	log.SetPrefix(hostname + "\t")

	log.Println("Bootstrap: build num.", vBuild)

	cnf = config.Parse(configfile)
	log.Println("Bootstrap: successful parsing config file. Items: ", cnf)
}

func main() {
	if cnf.Logging.Enable == true {
		logTo, err := os.OpenFile(cnf.Logging.Path, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0640)
		if err != nil {
			log.Fatal("Bootstrap: cannot open logfile - ", err)
		} else {
			log.Println("Bootstrap: logs switching to file - ", cnf.Logging.Path)
			defer logTo.Close()
		}
		log.SetOutput(logTo)
	}

	driver.Run(cnf)

}
