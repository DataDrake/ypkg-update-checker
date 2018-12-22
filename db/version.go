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

package db

import (
	"strconv"
	"strings"
	"unicode"
)

// Version is a record of a new version for a single source
type Version []string

func NewVersion(raw string) Version {
	dots := strings.Split(raw, ".")
	dashes := make([]string, 0)
	for _, dot := range dots {
		dashes = append(dashes, strings.Split(dot, "-")...)
	}
	pieces := make([]string, 0)
	for _, dash := range dashes {
		pieces = append(pieces, strings.Split(dash, "_")...)
	}
	v := make(Version, 0)
	i := 0
	if pieces[i][0] == 'v' || pieces[i][0] == 'V' {
		v = append(v, strings.TrimLeft(pieces[i], "vV"))
		i++
	}
	for i < len(pieces) && unicode.IsDigit(rune(pieces[i][0])) {
		v = append(v, pieces[i])
		i++
	}
	return v
}

// Compare allows to version nubmers to be compared to see which is newer (higher)
func (v Version) Compare(old Version) int {
	result := 0
	var curr, prev int
	var err error
	for i, piece := range v {
		if len(old) == i {
			return result
		}
		if old[i] == piece {
			continue
		}
		curr, err = strconv.Atoi(piece)
		if err != nil {
			goto HARD
		}
		prev, err = strconv.Atoi(old[i])
		if err != nil {
			goto HARD
		}
		result = prev - curr
		goto CHECK
	HARD:
		result = strings.Compare(piece, old[i])
	CHECK:
		if result != 0 {
			return result
		}
	}
	return result
}
