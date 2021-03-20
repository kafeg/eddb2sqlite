package main

import (
	"bufio"
	"compress/gzip"
	"database/sql"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
)

func downloadFile(filepath string, url string) (err error) {

	fmt.Printf("Download file '%s' from URL %s\n", filepath, url)

	// Create the file
	out, err := os.Create(filepath)
	if err != nil  {
		return err
	}
	defer out.Close()

	// Get the data
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return err
	}
	req.Header.Set("Accept-Encoding", "gzip, deflate, sdch")
	client := &http.Client{}
	resp, err := client.Do(req)
	defer resp.Body.Close()

	// Check server response
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("bad status: %s", resp.Status)
	}

	// Writer the body to file
	_, err = io.Copy(out, resp.Body)
	if err != nil  {
		return err
	}

	return nil
}

func gunzipFile(gzFilePath, dstFilePath string) (int64, error) {
	fmt.Printf("Gunzip file '%s' to '%s'\n", gzFilePath, dstFilePath)

	gzFile, err := os.Open(gzFilePath)
	if err != nil {
		return 0, fmt.Errorf("Failed to open file %s for unpack: %s", gzFilePath, err)
	}
	dstFile, err := os.OpenFile(dstFilePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0660)
	if err != nil {
		return 0, fmt.Errorf("Failed to create destination file %s for unpack: %s", dstFilePath, err)
	}

	ioReader, ioWriter := io.Pipe()

	go func() { // goroutine leak is possible here
		gzReader, _ := gzip.NewReader(gzFile)
		// it is important to close the writer or reading from the other end of the
		// pipe or io.copy() will never finish
		defer func(){
			gzFile.Close()
			gzReader.Close()
			ioWriter.Close()
		}()

		io.Copy(ioWriter, gzReader)
	}()

	written, err := io.Copy(dstFile, ioReader)
	if err != nil {
		return 0, err // goroutine leak is possible here
	}
	ioReader.Close()
	dstFile.Close()

	fi1, _ := os.Stat(gzFilePath)
	fi2, _ := os.Stat(dstFilePath)
	fmt.Printf("  %s (%.2fMB) -> %s (%.2fMB)\n", gzFilePath, float64(fi1.Size())/1024/1024, dstFilePath, float64(fi2.Size())/1024/1024)

	return written, nil
}

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

	referenceTables := []string{"allegiance", "government", "primary_economy", "security", "power_state", "controlling_minor_faction", "reserve_type"}

    tgtTblName := strings.ReplaceAll(csvFName, "data/", "")
	tgtTblName = strings.ReplaceAll(tgtTblName, ".csv", "")
	tgtTblName = strings.ReplaceAll(tgtTblName, "_recently", "")

    if !isTableExists(tgtTblName) {
    	return errors.New("Table " + tgtTblName + " does not exists")
	}

	file, err := os.Open(csvFName)
	if err != nil {
		return errors.New(err.Error())
	}
	defer file.Close()

	fmt.Printf("Processing file '%s'\n", csvFName)

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