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

package main

import (
	"runtime"
	"syscall"
	"time"
)

func maxOpenFiles() (rl int, err error) {
	if runtime.GOOS != "linux" {
		return
	} else {
		var rLimit syscall.Rlimit

		err = syscall.Getrlimit(syscall.RLIMIT_NOFILE, &rLimit)
		if err != nil {
			return
		}

		if rLimit.Cur < rLimit.Max {
			rLimit.Cur = rLimit.Max
			err = syscall.Setrlimit(syscall.RLIMIT_NOFILE, &rLimit)
			if err != nil {
				return
			}
			rl = int(rLimit.Cur)
		}
		return
	}
}

func nowFormatTime() string {
	t := time.Now().Local()
	return t.Format("2006-01-02")
}
