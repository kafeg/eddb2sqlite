package main

import (
	"db"
	"fmt"
)

func main() {
    fmt.Println(db.IsMigrated())

    if !db.IsMigrated() {
    	db.MigrateDatabase()
	}
}
