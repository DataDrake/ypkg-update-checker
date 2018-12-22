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
    "sort"
)

const removePackageQuery = "DROP * FROM releases WHERE package IN (?)"

func UpdatePackage(db *sqlx.DB, releases []Release) error {
    prev, err := GetReleases(db, releases[0].Package)
    if err != nil {
        return err
    }
    tx := db.MustBegin()
    for index, release := range releases {
        if index >= len(prev) {
            _, err = tx.NamedExec(insertReleaseQuery, release)
        } else if release.Updated.After(prev[index].Updated) {
            _, err = tx.NamedExec(updateReleaseQuery, release)
        } else {
            // do nothing
        }
        if err != nil {
            tx.Rollback()
            return err
        }
    }
    for i := len(releases); i < len(prev); i++ {
        _, err := tx.NamedExec(removeReleaseQuery, prev[i])
        if err != nil {
            tx.Rollback()
            return err
        }
    }
    return tx.Commit()
}

const getPackagesQuery = "SELECT package FROM releases GROUP BY package"

func CleanPackages(db *sqlx.DB, curr []string) error {
    sort.Strings(curr)
    prev := make([]string,0)
    err := db.Select(&prev, getPackagesQuery)
    if err != nil {
        return err
    }
    sort.Strings(prev)
    deletions := make([]string,0)
    for _, p := range prev {
        if sort.SearchStrings(curr, p) == len(curr) {
            deletions = append(deletions, p)
        }
    }
    if len(deletions) == 0 {
        return nil
    }
    query, args, err := sqlx.In(removePackageQuery, deletions)
    if err != nil {
        return err
    }
    query = db.Rebind(query)
    _, err = db.Exec(query, args)
    return err
}
