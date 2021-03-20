package main

import (
	"bufio"
	"database/sql"
	"errors"
	"fmt"
	"os"
	"strings"
)

func sliceContains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

func splitAtCommas(s string) []string {
	res := []string{}
	var beg int
	var inString bool

	for i := 0; i < len(s); i++ {
		if s[i] == ',' && !inString {
			res = append(res, s[beg:i])
			beg = i+1
		} else if s[i] == '"' {
			if !inString {
				inString = true
			} else if i > 0 && s[i-1] != '\\' {
				inString = false
			}
		}
	}
	return append(res, s[beg:])
}

func csv2sqlite(csvFName string) error {

	referenceTables := []string{"allegiance", "government", "primary_economy", "security", "power_state", "controlling_minor_faction", "reserve_type", "category", "type"}

    tgtTblName := strings.ReplaceAll(csvFName, "data/", "")
	tgtTblName = strings.ReplaceAll(tgtTblName, ".csv", "")
	tgtTblName = strings.ReplaceAll(tgtTblName, "_recently", "")
	tgtTblName = strings.ReplaceAll(tgtTblName, "_populated", "")

    if !isTableExists(tgtTblName) {
    	return errors.New("Table " + tgtTblName + " does not exists")
	}

	file, err := os.Open(csvFName)
	if err != nil {
		return errors.New(err.Error())
	}
	defer file.Close()

	fmt.Printf("SQL Processing file '%s'\n", csvFName)

	scanner := bufio.NewScanner(file)

	//csv header for columns
	scanner.Scan()
	cols := scanner.Text()
	colsArr := splitAtCommas(cols)

	db, _ := openDb(dbFileName)
	defer db.Close()

	var batchValues []string
	batchCount := 100000

	db.Exec("BEGIN TRANSACTION")

	//data + autofill reference tables
	for scanner.Scan() {
		row := scanner.Text()
		rowArr := splitAtCommas(row)

		for i, col := range colsArr {

			//check for refernce tbl
            if strings.Contains(col, "_id") {
            	colRef := strings.ReplaceAll(col, "_id", "")

            	//fmt.Println(len(rowArr), len(colsArr), col, row)
            	if rowArr[i] == "" || rowArr[i+1] == "" {
            		continue //skip empty vals
				}

            	if sliceContains(referenceTables, colRef) {
            		id := 0
            		stmt := fmt.Sprintf("SELECT id FROM %s WHERE %s = %s;", colRef, col, rowArr[i])
					//fmt.Println(stmt)
					err := db.QueryRow(stmt).Scan(&id)

            		if err != nil {

						preparedVal := strings.ReplaceAll(rowArr[i+1], "'", "`")
						preparedVal = strings.ReplaceAll(preparedVal, "\"", "")
						stmt := fmt.Sprintf("REPLACE INTO %s (%s, %s) VALUES('%s', '%s');", colRef, col, colRef, rowArr[i], preparedVal)
						//fmt.Println(stmt)
						//fmt.Println(cols)
						//fmt.Println(row)
						_, err := db.Exec(stmt)

						if err != nil {
							return errors.New(err.Error())
						}
					}
				}
			}
		}

		var valuesToInsert []string
		for i, colName := range colsArr {
			if !sliceContains(referenceTables, colName) {
				preparedVal := strings.ReplaceAll(rowArr[i], "'", "`")
				preparedVal = strings.ReplaceAll(preparedVal, "\"", "")
				valuesToInsert = append(valuesToInsert, "'"+preparedVal+"'")
			}
		}

		values := fmt.Sprintf("(%s)", strings.Join(valuesToInsert[:], ","))
		batchValues = append(batchValues, values)

		if len(batchValues) >= batchCount || len(row) == 0 {
			err := insertBatchRows(referenceTables, colsArr, tgtTblName, db, batchValues)
			batchValues = []string{}
			if err != nil {
				return errors.New(err.Error())
			}
		}
	}

	//drop remains data
	err = insertBatchRows(referenceTables, colsArr, tgtTblName, db, batchValues)
	batchValues = []string{}
	if err != nil {
		return errors.New(err.Error())
	}
	db.Exec("END TRANSACTION")

	fmt.Println("")

	return nil
}

func insertBatchRows(referenceTables, colsArr []string, tgtTblName string, db *sql.DB, batchValues []string) error {
	var colsToInsert []string
	for _, colName := range colsArr {
		if !sliceContains(referenceTables, colName) {
			if colName == "id" {
				colsToInsert = append(colsToInsert, "eddb_id")
			} else {
				colsToInsert = append(colsToInsert, colName)
			}
		}
	}
	fmt.Printf("Insert %d rows...\n", len(batchValues))
	stmt := fmt.Sprintf("INSERT INTO %s (%s) VALUES %s;", tgtTblName, strings.Join(colsToInsert[:], ","), strings.Join(batchValues[:], ","))
	_, err := db.Exec(stmt)

	if err != nil {
		return errors.New(stmt + " -> " + err.Error())
	}
	batchValues = []string{}

	return nil
}
