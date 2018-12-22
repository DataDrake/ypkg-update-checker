//
// Copyright 2018 Bryan T. Meyers <bmeyers@datadrake.com>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//

package cli

import (
	"fmt"
	"github.com/DataDrake/cli-ng/cmd"
	"github.com/DataDrake/ypkg-update-checker/db"
	"github.com/DataDrake/ypkg-update-checker/pkg"
	"io/ioutil"
	"os"
	"path/filepath"
)

// Update gets the most recent release for all package.yml files in a directory
var Update = cmd.CMD{
	Name:  "update",
	Alias: "u",
	Short: "Get the version and location for all identifiable sources",
	Args:  &UpdateArgs{},
	Run:   UpdateRun,
}

// UpdateArgs contains the arguments for the "update" subcommand
type UpdateArgs struct{}

func updateCheck(in, out chan pkg.Result, quit chan bool) {
	for {
		select {
		case r := <-in:
			r.Check()
			out <- r
		case <-quit:
			return
		}
	}
}

func updateGather(in chan pkg.Result, out chan pkg.Report, quit chan bool) {
	results := make(pkg.Report, 0)
	for {
		select {
		case r := <-in:
			results = append(results, &r)
		case <-quit:
			out <- results
			return
		}
	}
}

const updateWorkers = 4

// UpdateRun carries out finding the latest releases
func UpdateRun(r *cmd.RootCMD, c *cmd.CMD) {
    releases, err := db.Open()
    if err != nil {
		fmt.Printf("Failed to open database, reason: \"%s\"\n", err.Error())
		os.Exit(1)
    }
    defer releases.Close()
	fail := 0
	files, err := ioutil.ReadDir(".")
	if err != nil {
		fmt.Printf("Failed to get files in directory, reason: \"%s\"\n", err.Error())
		os.Exit(1)
	}
	in := make(chan pkg.Result)
	out := make(chan pkg.Result)
	final := make(chan pkg.Report)
	quit := make(chan bool)
	quit2 := make(chan bool)
	go updateGather(out, final, quit2)
	for i := 0; i < updateWorkers; i++ {
		go updateCheck(in, out, quit)
	}
	for _, file := range files {
		if !file.IsDir() {
			continue
		}
		fmt.Fprintf(os.Stderr, "Processing %s\n", file.Name())
		r, err := pkg.NewResult(filepath.Join(".", file.Name(), "package.yml"))
		if err != nil {
			fmt.Fprintf(os.Stderr, "%s failed, reason: %s\n", file.Name(), err.Error())
			fail++
			continue
		}
		in <- *r
	}
	for i := 0; i < updateWorkers; i++ {
		quit <- true
	}
	quit2 <- true
	results := <-final
	results.Print(fail)
	os.Exit(0)
}
