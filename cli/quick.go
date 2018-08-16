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
	"os"
)

// Quick gets the most recent release for a given source, without pretty printing
var Quick = cmd.CMD{
	Name:  "quick",
	Alias: "q",
	Short: "Get the version and location of the most recent release",
	Args:  &QuickArgs{},
	Run:   QuickRun,
}

// QuickArgs contains the arguments for the "quick" subcommand
type QuickArgs struct {
	Path string `desc:"Location of package.yml"`
}

// QuickRun carries out finding the latest release
func QuickRun(r *cmd.RootCMD, c *cmd.CMD) {
	args := c.Args.(*QuickArgs)

	yml, err := pkg.Open(args.Path)
	if err != nil {
		fmt.Printf("Failed to open package.yml, reason: \"%s\"\n", err.Error())
		os.Exit(1)
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
				fmt.Printf("%s %s %s %s\n", yml.Name, yml.Version, r.Version, r.Location)
			}
		}
	}
	if !found {
		fmt.Printf("No release(s) found for '%s'.\n", yml.Name)
	}
	os.Exit(0)
}
