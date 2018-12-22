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
)

// Report generates a report of the last update
var Report = cmd.CMD{
	Name:  "report",
	Alias: "r",
	Short: "Generates a report of the last update",
	Args:  &ReportArgs{},
	Run:   ReportRun,
}

// ReportArgs contains the arguments for the "report" subcommand
type ReportArgs struct{}

// ReportRun carries out finding the latest releases
func ReportRun(r *cmd.RootCMD, c *cmd.CMD) {
    rdb, err := db.Open()
	if err != nil {
		fmt.Printf("Failed to open database, reason: \"%s\"\n", err.Error())
		os.Exit(1)
	}
    releases, err := db.GetAllReleases(rdb)
	if err != nil {
		fmt.Printf("Failed to read database, reason: \"%s\"\n", err.Error())
		os.Exit(1)
	}
    report := pkg.NewReport(releases)
    report.Print()
}
