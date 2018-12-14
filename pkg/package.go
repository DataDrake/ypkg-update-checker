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
	"gopkg.in/yaml.v2"
	"os"
)

// PackageYML is a Go representation of a package.yml file
type PackageYML struct {
	Name    string              `yaml:"name"`
	Version string              `yaml:"version"`
	Sources []map[string]string `yaml:"source"`
}

// Open parses a package.yml into a struct and returns it
func Open(path string) (yml *PackageYML, err error) {
	ymlFile, err := os.Open(path)
	if err != nil {
		return
	}
	defer ymlFile.Close()
	dec := yaml.NewDecoder(ymlFile)
	yml = &PackageYML{}
	err = dec.Decode(yml)
	return
}
