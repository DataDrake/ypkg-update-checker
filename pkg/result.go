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
	"errors"
	"github.com/DataDrake/cuppa/providers"
	"github.com/DataDrake/cuppa/results"
)

// Result is a single result for use in a report
type Result struct {
	YML         *PackageYML
	NewVersions map[string]Version
}

// NewResult attempts to look up new sources for a package
func NewResult(path string) (r *Result, err error) {
	yml, err := Open(path)
	if err != nil {
		return
	}
	r = &Result{
		YML:         yml,
		NewVersions: make(map[string]Version),
	}
	return
}

// NotFound signifies that a matching upstream could not be identified
var NotFound error

func init() {
	NotFound = errors.New("Could not find a matching provider")
}

// Check attempts to find new sources for every source in the contained PackageYML
func (r *Result) Check() {
	for _, srcs := range r.YML.Sources {
		for src := range srcs {
			var found bool
			for _, p := range providers.All() {
				name := p.Match(src)
				if name == "" {
					continue
				}
				result, s := p.Latest(name)
				if s != results.OK || result == nil {
					continue
				}
				found = true
				v := Version{
					Number:   result.Version,
					Location: result.Location,
				}
				r.NewVersions[src] = v
				break
			}
			if !found {
				v := Version{
					Error: NotFound,
				}
				r.NewVersions[src] = v
			}
		}
	}
}
