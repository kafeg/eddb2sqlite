module eddb2sqlite

go 1.14

require (
	db v1.0.0
	github.com/mattn/go-sqlite3 v1.14.6 // indirect
	github.com/rubenv/sql-migrate v0.0.0-20210215143335-f84234893558 // indirect
)

replace db => ./db
