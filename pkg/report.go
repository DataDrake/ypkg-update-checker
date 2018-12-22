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
	"github.com/DataDrake/ypkg-update-checker/db"
	"math"
	"net/url"
	"sort"
	"strings"
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
.behind {background-color: #F00; color: black;}
.held {background-color: #F93; color: black;}
.ok {background-color: #0F0; color: black;}
.ahead {background-color: #0EF; color: black;}
</style>
</head>
<body>
`

// ReportSummary is a format string for a summary of a report
const ReportSummary = `
<h1 id="summary">Summary</h1>
<div style="display: flex; height: 1rem; padding: 0.7rem;">
<div class="behind" style="flex: %d;"></div>
<div class="held" style="flex: %d;"></div>
<div class="ok" style="flex: %d;"></div>
<div class="ahead" style="flex: %d;"></div>
</div>
<table>
<tr><td>Matched: </td><td>                                    </td><td>  </td></tr>
<tr><td>         </td><td class="behind"> Out of Date         </td><td>%d</td></tr>
<tr><td>         </td><td class="held">   Held Behind         </td><td>%d</td></tr>
<tr><td>         </td><td class="ok">     Up to Date          </td><td>%d</td></tr>
<tr><td>         </td><td class="ahead">  Newer than Upstream </td><td>%d</td></tr>
<tr><td>Unmatched</td><td>                                    </td><td>%d</td></tr>
<tr><td>Failed   </td><td>                                    </td><td>%d</td></tr>
<tr><td>Total    </td><td>                                    </td><td>%d</td></tr>
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
type Report struct {
	matched        []db.Release
	unmatched      map[string][]db.Release
	failed         []db.Release
	unmatchedCount int
	outOfDateCount int
	heldBackCount  int
	upToDateCount  int
	aheadCount     int
}

func NewReport(releases []db.Release) *Report {
	r := &Report{
        unmatched: make(map[string][]db.Release),
    }
	for _, release := range releases {
		switch release.Status {
		case db.StatusUnmatched:
			r.unmatchedCount++
			hostname := "N/A"
			pieces := strings.Split(release.Source, "|")
			loc := pieces[len(pieces)-1]
			host, err := url.Parse(loc)
			if err == nil {
				pieces := strings.Split(host.Hostname(), ".")
				hostname = pieces[len(pieces)-2]
			}
			r.unmatched[hostname] = append(r.unmatched[hostname], release)
		case db.StatusOutOfDate:
			r.outOfDateCount++
			r.matched = append(r.matched, release)
		case db.StatusHeldBack:
			r.heldBackCount++
			r.matched = append(r.matched, release)
		case db.StatusUpToDate:
			r.upToDateCount++
			r.matched = append(r.matched, release)
		case db.StatusAhead:
			r.aheadCount++
			r.matched = append(r.matched, release)
		default:
			r.failed = append(r.failed, release)
		}
	}
	return r
}

// Print generates an HTML report
func (r Report) Print() {
	fmt.Println(ReportStart)
	matched := r.outOfDateCount + r.heldBackCount + r.upToDateCount + r.aheadCount
	behindP := int(math.Floor(float64(r.outOfDateCount) / float64(matched) * 100.0))
	heldP := int(math.Floor(float64(r.heldBackCount) / float64(matched) * 100.0))
	okP := int(math.Floor(float64(r.upToDateCount) / float64(matched) * 100.0))
	aheadP := int(math.Floor(float64(r.aheadCount) / float64(matched) * 100.0))
	fmt.Printf(ReportSummary, behindP, heldP, okP, aheadP,
		r.outOfDateCount, r.heldBackCount, r.upToDateCount, r.aheadCount,
		r.unmatchedCount, len(r.failed), matched+r.unmatchedCount+len(r.failed))
	fmt.Println(ReportMatchHeader)
	for _, release := range r.matched {
		var color string
		switch release.Status {
		case db.StatusOutOfDate:
			color = "behind"
		case db.StatusHeldBack:
			color = "held"
		case db.StatusAhead:
			color = "ahead"
		default:
			color = "ok"
		}
		fmt.Printf(ReportMatchRow, release.Package, release.Current, color, release.Latest, release.Source, release.Source)
	}
	fmt.Println(ReportTableClose)
	fmt.Println(ReportUnmatchedHeader)
	hosts := make([]string, 0)
	for host := range r.unmatched {
		hosts = append(hosts, host)
	}
	sort.Strings(hosts)
	for _, host := range hosts {
		fmt.Printf(ReportUnmatchedSectionStart, host)
		for _, release := range r.unmatched[host] {
			fmt.Printf(ReportUnmatchedRow, release.Package, release.Current, release.Source, release.Source)
		}
		fmt.Println(ReportUnmatchedSectionStop)
	}
	fmt.Println(ReportUnmatchedClose)
}
