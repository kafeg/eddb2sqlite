package main

import (
	"flag"
	"fmt"
	"os"
	"strings"
)

func main() {

	// migrate database
	if !isMigrated() {
		fmt.Println("Migrating DB...")
		migrateDatabase()
		fmt.Println("Done!")
	}

	//download, unpack and convert files
	os.MkdirAll("data", os.ModePerm)

	urls := make(map[string]string)
	urls["data/commodities.json.gz"] = "https://eddb.io/archive/v6/commodities.json"
	//urls["data/systems.csv.gz"] = "https://eddb.io/archive/v6/systems.csv"
	urls["data/systems_populated.json.gz"] = "https://eddb.io/archive/v6/systems_populated.json"
	urls["data/stations.json.gz"] = "https://eddb.io/archive/v6/stations.json"
	//urls["data/attractions.json.gz"] = "https://eddb.io/archive/v6/attractions.json"
	urls["data/factions.csv.gz"] = "https://eddb.io/archive/v6/factions.csv"
	urls["data/listings.csv.gz"] = "https://eddb.io/archive/v6/listings.csv"
	//urls["data/systems_recently.csv.gz"] = "https://eddb.io/archive/v6/systems_recently.csv"
	//urls["data/modules.json.gz"] = "https://eddb.io/archive/v6/modules.json"

	for k, v := range urls {
		if _, err := os.Stat(k); err != nil {

			//download
			err := downloadFile(k, v)
			if err != nil {
				fmt.Printf("Error: %s", err)
			}

			//decompress
			unpackFName := strings.ReplaceAll(k, ".gz", "")
			_, err = gunzipFile(k, unpackFName)
			if err != nil {
				fmt.Printf("Error: %s", err)
			}

			//convert jsons to csvs
			if strings.Contains(unpackFName, ".json") {
				tgtFName := strings.ReplaceAll(unpackFName, ".json", ".csv")
				json2csv(unpackFName, tgtFName)
				flag.Parse()
			}
			fmt.Println("")
		}
	}

    //import csvs to sqlite
	csvs := []string{
		"data/factions.csv",
		"data/commodities.csv",
		"data/systems_populated.csv",
		"data/stations.csv",
		"data/listings.csv",
	}

	for _, s := range csvs {
		err := csv2sqlite(s)

		if err != nil {
			fmt.Printf("Error: %s\n", err)
			os.Exit(1)
		}
	}

    os.Exit(0)
}
