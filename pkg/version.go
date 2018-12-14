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
	"regexp"
	"strconv"
	"strings"
)

var versionRegex *regexp.Regexp

func init() {
	versionRegex = regexp.MustCompile("(\\d+(?:[._]\\d+)*[a-zA-z]?)")
}

// Version is a record of a new version for a single source
type Version struct {
	Number   string
	Location string
	Error    error
	pieces   []string
}

func versionToPieces(version string) []string {
	all := versionRegex.FindString(version)
	return strings.Split(all, ".")
}

// Compare allows to version nubmers to be compared to see which is newer (higher)
func (v Version) Compare(old string) int {
	if v.pieces == nil || len(v.pieces) == 0 {
		v.pieces = versionToPieces(v.Number)
	}
	piecesOld := versionToPieces(old)
	result := 0
	var curr, prev int
	var err error
	for i, piece := range v.pieces {
		if len(piecesOld) == i {
			return result
		}
		if piecesOld[i] == piece {
			continue
		}
		curr, err = strconv.Atoi(piece)
		if err != nil {
			goto HARD
		}
		prev, err = strconv.Atoi(piecesOld[i])
		if err != nil {
			goto HARD
		}
		result = prev - curr
		goto CHECK
	HARD:
		result = strings.Compare(piece, piecesOld[i])
	CHECK:
		if result != 0 {
			return result
		}
	}
	return result
}
