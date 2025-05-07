package core

import (
	"database/sql"
)

func newDB(engine *Engine) (db *sql.DB, err error) {
	if ndb, nde := engine.NewDB(); nil == nde {
		db = ndb.DB
	}

	return
}
