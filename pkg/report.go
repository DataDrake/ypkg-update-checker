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

package pkg

import (
	"fmt"
	"math"
	"sort"
)

// ReportStart is the header at the beginning of a report
const ReportStart = `
<html>
<head>
<style>
html {
    background-color: #333;
    color: #EEE;
    font-family: Hack, monospace;}
body { overflow: none; }
table {margin: 2em;};
th {text-align: left;
    border-bottom: 0.125rem solid #EEE;}
td {padding: 0 0.7rem;}
a { color: #eee; text-decoration: none;1}
.red {background-color: #F00; color: black;}
.green {background-color: #0F0; color: black;}
.blue {background-color: #0EF; color: black;}
</style>
</head>
<body>
`

// ReportSummary is a format string for a summary of a report
const ReportSummary = `
<h1 id="summary">Summary</h1>
<div style="display: flex; height: 1rem; padding: 0.7rem;">
<div class="red" style="flex: %d;"></div>
<div class="green" style="flex: %d;"></div>
<div class="blue" style="flex: %d;"></div>
</div>
<table>
<tr><td>Matched: </td><td>                                 </td><td>  </td></tr>
<tr><td>         </td><td class="red">  Out of Date        </td><td>%d</td></tr>
<tr><td>         </td><td class="green">Up to Date         </td><td>%d</td></tr>
<tr><td>         </td><td class="blue"> Newer than Upstream</td><td>%d</td></tr>
<tr><td>Unmatched</td><td>                                 </td><td>%d</td></tr>
<tr><td>Failed   </td><td>                                 </td><td>%d</td></tr>
<tr><td>Total    </td><td>                                 </td><td>%d</td></tr>
</table>
<h3><a href="#unmatched">Go to Unmatched Packages</a></h3>
`

// ReportMatchHeader is the header for matched packages
const ReportMatchHeader = `
<h1 id="matched">Matched Packages</h1>
<table>
<thead>
<tr><th>Name</th><th>Old Version</th><th>New Version</th><th>Location</th></tr>
</thead>
<tbody>
`

// ReportMatchRow is the format string for a row of the matched packages
const ReportMatchRow = "<tr><td>%s</td><td>%s</td><td class=\"%s\">%s</td><td><a href=\"%s\">%s</a></td></tr>\n"

// ReportTableClose terminates a table in the report
const ReportTableClose = "</tbody></table>\n"

// ReportUnmatchedHeader is the header for unmatched packages
const ReportUnmatchedHeader = `
<h1 id="unmatched">Unmatched Packages</h1>
<h3><a href="#summary">Back to Top</a></h3>
`
const ReportUnmatchedSectionStart = `
<h3>%s</h3>
<table>
<thead>
<tr><th>Name</th><th>Old Version</th><th>Location</th></tr>
</thead>
<tbody>
`

const ReportUnmatchedSectionStop = "</tbody></table>"

// ReportUnmatchedRow is the format strign for a row of the unmatched packages
const ReportUnmatchedRow = "<tr><td>%s</td><td>%s</td><td><a href=\"%s\">%s</a></td></tr>\n"

// ReportUnmatchedClose terminates the report
const ReportUnmatchedClose = `
<h3><a href="#summary">Back to Top</a></h3>
</body>
`

// Report is a record of multiple package checks
type Report []*Result

// Len is used for sorting
func (r Report) Len() int {
	return len(r)
}

// Less is used for sorting
func (r Report) Less(i, j int) bool {
	if r[i].First < r[j].First {
        return true
    }
	if r[i].First > r[j].First {
        return false
    }
    return r[i].YML.Name < r[j].YML.Name
}

// Swap is used for sorting
func (r Report) Swap(i, j int) {
	r[i], r[j] = r[j], r[i]
}

// Print generates an HTML report
func (r Report) Print(failed int) {
	sort.Sort(r)
	exact := 0
	greater := 0
	less := 0
	unmatched := 0
	for _, result := range r {
		for _, version := range result.NewVersions {
			if version.Error == nil {
				cmp := version.Compare(result.YML.Version)
				if cmp == 0 {
					exact++
				} else if cmp > 0 {
					greater++
				} else {
					less++
				}
			} else if version.Error == NotFound {
				unmatched++
			}
		}
	}
	fmt.Println(ReportStart)
	total := less + exact + greater
	lessP := int(math.Floor(float64(less) / float64(total) * 100.0))
	exactP := int(math.Floor(float64(exact) / float64(total) * 100.0))
	greaterP := int(math.Floor(float64(greater) / float64(total) * 100.0))
	fmt.Printf(ReportSummary, lessP, exactP, greaterP,
		less, exact, greater,
		unmatched, failed, less+exact+greater+unmatched+failed)
	fmt.Println(ReportMatchHeader)
	for _, result := range r {
		for _, version := range result.NewVersions {
			if version.Error == nil {
				status := "red"
				cmp := version.Compare(result.YML.Version)
				if cmp == 0 {
					status = "green"
				} else if cmp > 0 {
					status = "blue"
				}
				fmt.Printf(ReportMatchRow, result.YML.Name, result.YML.Version, status, version.Number, version.Location, version.Location)
			}
		}
	}
	fmt.Println(ReportTableClose)
	fmt.Println(ReportUnmatchedHeader)
    host := r[0].First
    fmt.Printf(ReportUnmatchedSectionStart, host)
	for _, result := range r {
        if result.First != host {
            fmt.Println(ReportUnmatchedSectionStop)
            host = result.First
            fmt.Printf(ReportUnmatchedSectionStart, host)
        }
		for src, version := range result.NewVersions {
			if version.Error == NotFound {
				fmt.Printf(ReportUnmatchedRow, result.YML.Name, result.YML.Version, src, src)
			}
		}
	}
    fmt.Println(ReportUnmatchedSectionStop)
	fmt.Println(ReportUnmatchedClose)
}
