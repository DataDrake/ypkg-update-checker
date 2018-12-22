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
	"os/user"
)

const getTablesQuery = "SELECT name FROM sqlite_master WHERE type='table'"

const releaseSchema = `
CREATE TABLE releases (
    package TEXT,
    source TEXT,
    current TEXT,
    latest TEXT,
    updated DATETIME,
    status INTEGER,
    idx  INTEGER
);
`

func CreateTables(db *sqlx.DB) error {
	found, err := db.Queryx(getTablesQuery)
	if err != nil {
		return err
	}
	missing := true
	for found.Next() {
		var table string
		err = found.Scan(&table)
		if err != nil {
			return err
		}
		if table == "releases" {
			missing = false
		}
	}
	if missing {
		_, err := db.Exec(releaseSchema)
		if err != nil {
			return err
		}
	}
	return nil
}

func Open() (db *sqlx.DB, err error) {
	u, err := user.Current()
	if err != nil {
		return
	}
	db, err = sqlx.Connect("sqlite3", u.HomeDir+"/.cache/ypkg-update.db")
	if err != nil {
		return
	}
	err = CreateTables(db)
	return
}
