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
    "strings"
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
const ReportSummary = "</table><p>Failed: %d</p><p>Unmatched: %d</p><p>Total: %d</p></body></html>"


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
    matched := 0
    unmatched := 0
    for _, result := range r {
        for _, version := range result.NewVersions {
            if version.Error == nil {
                status := "red"
                if result.YML.Version == version.Number {
                    status = "green"
                } else if strings.Contains(version.Number, result.YML.Version) {
                    status = "green"
                }
                fmt.Printf(ReportMatchRow, result.YML.Name, result.YML.Version, status, version.Number, version.Location );
            }
            matched++
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
    fmt.Printf(ReportSummary, failed, unmatched, failed+matched+unmatched);
}
