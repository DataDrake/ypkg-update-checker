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
    "sort"
)

const ReportMatchHeader =`
<html>
<head>
<style>
.red {background-color: #F00; color: black;}
.green {background-color: #0F0; color: black;}
.yellow {background-color: yellow; color: black;}
</style>
</head>
<body>
<h2>Matched Packages</h2>
<hr/>
<table>
<thead>
<tr><th>Name</th><th>Old Version</th><th>New Version</th><th>Location</th></tr>
</thead>
<tbody>
`
const ReportMatchRow = "<tr><td>%s</td><td>%s</td><td class=\"%s\">%s</td><td>%s</td></tr>\n"
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
const ReportSummary = `</table>
<h2>Summary</h2>
<table>
<tr><td>Matched: </td><td>                   </td><td>  </td></tr>
<tr><td>         </td><td>Out of Date        </td><td>%d</td></tr>
<tr><td>         </td><td>Up to Date         </td><td>%d</td></tr>
<tr><td>         </td><td>Newer than Upstream</td><td>%d</td></tr>
<tr><td>Unmatched</td><td>                   </td><td>%d</td></tr>
<tr><td>Failed   </td><td>                   </td><td>%d</td></tr>
<tr><td>Total    </td><td>                   </td><td>%d</td></tr>
</table>
</body></html>`


// Report is a record of multiple package checks
type Report []*Result

// Len is used for sorting
func (r Report) Len() int {
    return len(r);
}

// Less is used for sorting
func (r Report) Less(i, j int) bool {
    return r[i].YML.Name < r[j].YML.Name
}

// Swap is used for sorting
func (r Report) Swap(i, j int) {
    r[i], r[j] = r[j], r[i]
}

// Print generates an HTML report
func (r Report) Print(failed int) {
    sort.Sort(r)
    fmt.Println(ReportMatchHeader);
    exact := 0
    greater := 0
    less :=0
    unmatched := 0
    for _, result := range r {
        for _, version := range result.NewVersions {
            if version.Error == nil {
                status := "red"
                cmp := version.Compare(result.YML.Version)
                if cmp == 0 {
                    status = "green"
                    exact++
                } else if cmp > 0 {
                    status = "yellow"
                    greater++
                } else {
                    less++
                }
                fmt.Printf(ReportMatchRow, result.YML.Name, result.YML.Version, status, version.Number, version.Location );
            }
        }
    }
    fmt.Println(ReportTableClose);
    fmt.Println(ReportUnmatchedHeader);
    for _, result := range r {
        for src, version := range result.NewVersions {
            if version.Error == NotFound {
                fmt.Printf(ReportUnmatchedRow, result.YML.Name, result.YML.Version, src );
                unmatched++
            }
        }
    }
    fmt.Printf(ReportSummary, less, exact, greater, unmatched, failed, less+exact+greater+unmatched+failed);
}
