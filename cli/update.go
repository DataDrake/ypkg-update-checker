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

func updateCheck(rdb *sqlx.DB, in chan string, quit chan bool) {
	for {
		select {
		case p := <-in:
            prev := db.GetReleases(rdb, p)
            curr := make([]db.Release)
		    yml, err := pkg.Open(filepath.Join(".", p, "package.yml"))
		    if err != nil {
                if err == os.ErrNotExist {
                    curr = append(curr,
                        db.Release{
                            Package: p,
                            Updated: time.Now(),
                            Index: 0,
                            Status: db.StatusMissingYML
                        }
                    )
                } else {
			        fmt.Fprintf(os.Stderr, "%s failed, reason: %s\n", file.Name(), err.Error())
        	        continue
                }
            } else {
                for index, src := range r.YML.Sources {
                    for location := range src {
                        
                    }
                }
		    }
		r, err := pkg.NewResult(filepath.Join(".", p, "package.yml"))
		if err != nil {
			fmt.Fprintf(os.Stderr, "%s failed, reason: %s\n", file.Name(), err.Error())
			fail++
			continue
		}
			r.Check()
			out <- r
		case <-quit:
			return
		}
	}
}

const updateWorkers = 4

// UpdateRun carries out finding the latest releases
func UpdateRun(r *cmd.RootCMD, c *cmd.CMD) {
    rdb, err := db.Open()
    if err != nil {
		fmt.Printf("Failed to open database, reason: \"%s\"\n", err.Error())
		os.Exit(1)
    }
    defer rdb.Close()
	fail := 0
	files, err := ioutil.ReadDir(".")
	if err != nil {
		fmt.Printf("Failed to get packages, reason: \"%s\"\n", err.Error())
		os.Exit(1)
	}
    packages := make([]string,0)
	for _, file := range files {
		if !file.IsDir() {
			continue
		}
        if file.Name() == "common" {
            continue
        }
		packages = append(packages, file.Name())
	}
    err = db.CleanPackages(rdb, packages)
    if err != nil {
		fmt.Printf("Failed to clean up packages, reason: \"%s\"\n", err.Error())
		os.Exit(1)
    }
	in := make(chan string)
	quit := make(chan bool)
	for i := 0; i < updateWorkers; i++ {
		go updateCheck(rdb, in, quit)
	}
	for _, p := range packages {
		fmt.Fprintf(os.Stderr, "Starting %s\n", p)
		in <- p
	}
	for i := 0; i < updateWorkers; i++ {
		quit <- true
	}
	os.Exit(0)
}
