package testutil

import (
	"database/sql"
	"log"

	"github.com/CzarSimon/httputil/dbutil"
)

// InMemoryDB creates an in memory sql database and applies migrations if desired.
func InMemoryDB(migrate bool, migrationsPath string) *sql.DB {
	cfg := dbutil.SqliteConfig{}
	db := dbutil.MustConnect(cfg)
	_, err := db.Exec("PRAGMA foreign_keys = ON")
	if err != nil {
		log.Panic("Failed to activate foregin keys ", err)
	}

	if !migrate {
		return db
	}

	err = dbutil.Upgrade(migrationsPath, cfg.Driver(), db)
	if err != nil {
		log.Panic("Failed to apply upgrade migratons", err)
	}

	return db
}
