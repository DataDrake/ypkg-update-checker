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
    "github.com/jmoiron/sqlx"
    "time"
)

const (
    StatusMissingYML = -4
    StatusUnmatched  = -3
    StatusOutOfDate  = -2
    StatusHeldBack   = -1
    StatusUpToDate   = 0
    StatusAhead      = 1
)

const getReleasesQuery = "SELECT * FROM releases WHERE package=? ORDER BY index"
const getAllReleasesQuery = "SELECT * FROM releases WHERE package=? ORDER BY package, index"
const insertReleaseQuery = "INSERT INTO releases VALUES (:package, :source, :current, :latest, :updated, :status, :index)"
const updateReleaseQuery = `
UPDATE releases
SET
    source=:source,
    current=:current,
    latest=:latest,
    updated=:updated,
    status=:status,
WHERE package=:package AND index=:index`
const removeReleaseQuery = "DROP * FROM releases WHERE package=:package AND index=:index"

type Release struct {
    Package string
    Source string
    Current string
    Latest string
    Updated time.Time
    Status int
    Index  int
}

func GetReleases(db *sqlx.DB, name string) ([]Release, error) {
    releases := make([]Release,0)
    err := db.Select(&releases, getReleasesQuery, name)
    return releases, err
}

func GetAllReleases(db *sqlx.DB) ([]Release, error) {
    releases := make([]Release,0)
    err := db.Get(&releases, getAllReleasesQuery)
    return releases, err
}
