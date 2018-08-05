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
	"github.com/DataDrake/cuppa/providers"
	"github.com/DataDrake/cuppa/results"
	"github.com/DataDrake/ypkg-update-checker/pkg"
    "io/ioutil"
  	"os"
    "path/filepath"
)

// Report gets the most recent release for all package.yml files in a directory
var Report = cmd.CMD{
	Name:  "report",
	Alias: "r",
	Short: "Get the version and location for all identifiable sources",
	Args:  &ReportArgs{},
	Run:   ReportRun,
}

// ReportArgs contains the arguments for the "report" subcommand
type ReportArgs struct {}


const ReportMatchHeader =`
<html>
<body>
<h2>Matched Packages</h2>
<hr/>
<table>
<thead>
<tr><th>Name</th><th>Old Version</th><th>New Version</th><th>Location</th></tr>
</thead>
<tbody>
`
const ReportMatchRow = "<tr><td>%s</td><td>%s</td><td>%s</td><td>%s</td></tr>\n"
const ReportTableClose = "</tbody></table>\n"

const ReportUnmatchedHeader =`
<h2>Unmatched Packaged</h2>
<hr/>
<table>
<thead>
<tr><th>Name</th><th>Old Version</th><th>Location</th></tr>
</thead>
<tbody>
`

const ReportUnmatchedRow = "<tr><td>%s</td><td>%s</td><td>%s</td></tr>\n"
const ReportSummary = "</table><p>Failed: %d</p><p>Unmatched: %d</p><p>Total: %d</p></body></html>"

// ReportRun carries out finding the latest releases
func ReportRun(r *cmd.RootCMD, c *cmd.CMD) {

    fail  := 0
    total := 0

    unmatched := make([]*pkg.PackageYML,0)

    files, err := ioutil.ReadDir(".")
    if err != nil {
        fmt.Printf("Failed to get files in directory, reason: \"%s\"\n", err.Error())
        os.Exit(1)
    }
    fmt.Println(ReportMatchHeader)
    for _, file := range files {
        if !file.IsDir() {
            continue
        }
        yml, err := pkg.Open(filepath.Join(".", file.Name(), "package.yml"))
        if err != nil {
            //fmt.Printf("Failed to open package.yml, reason: \"%s\"\n", err.Error())
            fail++
            continue
        }
	    found := false
	    for _, p := range providers.All() {
            for _, srcs := range yml.Sources {
                for src, _ := range srcs {
		            name := p.Match(src)
		            if name == "" {
			            continue
		            }
		            r, s := p.Latest(name)
		            if s != results.OK || r == nil {
			            continue
		            }
                    found = true
                    fmt.Printf(ReportMatchRow, yml.Name, yml.Version, r.Version, r.Location)
                }
		    }
        }
        if found {
            total++
        } else {
            unmatched = append(unmatched, yml)
            fail++
        }
	}
    fmt.Printf(ReportTableClose)
    fmt.Printf(ReportUnmatchedHeader)
    for _, yml := range unmatched {
        for _, srcs := range yml.Sources {
            for src, _ := range srcs {
                fmt.Printf(ReportUnmatchedRow, yml.Name, yml.Version, src)
            }
        }
    }
    fmt.Printf(ReportTableClose)
    fmt.Printf(ReportSummary, fail, len(unmatched), total+fail+len(unmatched))
	os.Exit(0)
}
