package db

import (
	"database/sql"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"github.com/rubenv/sql-migrate"
)

var dbFileName = "db/eddb.sqlite"
var dbDialect = "sqlite3"
var migrationsDir = "db"

func openDb(fileName string) (*DB, error)  {
	db, err := sql.Open("sqlite3", fileName)
	if err != nil {
		fmt.Print("Error on open DB file!")
	}

	return db, err
}

func migrateDatabase() {

	migrations := &migrate.FileMigrationSource{
		Dir: migrationsDir,
	}

	db, err := openDb(dbFileName)

	n, err := migrate.Exec(db, dbDialect, migrations, migrate.Up)
	if err != nil {
		fmt.Print("Error on exec migrations!")
	}
	fmt.Printf("Applied %d migrations!\n", n)
}

func isMigrated() bool {
	src := &migrate.FileMigrationSource{
		Dir: migrationsDir,
	}

	migrations, err := src.FindMigrations()
	if err != nil {
		return false
	}

	db, err := openDb(dbFileName)

	records, err := migrate.GetMigrationRecords(db, dbDialect)
	if err != nil {
		return false
	}

	for _, m := range migrations {
		for _, r := range records {
			
		}
	}

	return true
}