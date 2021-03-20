package main

import (
	"database/sql"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"github.com/rubenv/sql-migrate"
)

var dbFileName = "db/eddb.sqlite"
var dbDialect = "sqlite3"
var migrationsDir = "db"

func openDb(fileName string) (*sql.DB, error)  {
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
	defer db.Close()

	n, err := migrate.Exec(db, dbDialect, migrations, migrate.Up)
	if err != nil {
		fmt.Printf("Error on exec migrations: %s\n", err)
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
	defer db.Close()

	records, err := migrate.GetMigrationRecords(db, dbDialect)
	if err != nil {
		return false
	}

	allFound := false
	for _, m := range migrations {
		found := false
		for _, r := range records {
			//fmt.Println(m, r)
			if r.Id == m.Id {
				found = true
				break
			}
		}
		allFound = found
		if !allFound {
			break
		}
	}

	return allFound
}

func isTableExists(tblName string) bool {
	db, _ := openDb(dbFileName)
	defer db.Close()

	_, err := db.Query("SELECT 1 FROM " + tblName + ";")

	return err == nil
}