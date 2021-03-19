package main

import (
	"fmt"
	"os"
	"strings"
)

func main() {

	if len(os.Args) < 2 {
		fmt.Println("Wrong args count. Available commands: migrate, status, getdata, generate")
		os.Exit(1)
	}

	switch os.Args[1] {
	case "status": {
		fmt.Printf("DB is migrated: %t, DB main data imported: %s, DB recent data updated: %s", IsMigrated(), "01.01.2020", "10.10.2021")
	}
	case "migrate": {
		if !IsMigrated() {
			fmt.Println("Migrating DB...")
			MigrateDatabase()
			fmt.Println("Done!")
		} else {
			fmt.Println("DB alrteady migrated")
		}
	}
	case "getdata": {
		force := len(os.Args) > 2 && os.Args[2] == "--force"
		//recentOnly := len(os.Args) > 2 && os.Args[2] == "--recent"
		os.MkdirAll("data", os.ModePerm)

		urls := make(map[string]string)
		urls["data/commodities.json.gz"] = "https://eddb.io/archive/v6/commodities.json"
		urls["data/systems.csv.gz"] = "https://eddb.io/archive/v6/systems.csv"
		urls["data/systems_populated.json.gz"] = "https://eddb.io/archive/v6/systems_populated.json"
		urls["data/stations.json.gz"] = "https://eddb.io/archive/v6/stations.json"
		urls["data/attractions.json.gz"] = "https://eddb.io/archive/v6/attractions.json"
		urls["data/factions.csv.gz"] = "https://eddb.io/archive/v6/factions.csv"
		urls["data/listings.csv.gz"] = "https://eddb.io/archive/v6/listings.csv"
		urls["data/systems_recently.csv.gz"] = "https://eddb.io/archive/v6/systems_recently.csv"
		urls["data/modules.json.gz"] = "https://eddb.io/archive/v6/modules.json"

		for k, v := range urls {
			if _, err := os.Stat(k); err != nil || force {
				if force {
					os.Remove(k)
				}
				err := downloadFile(k, v)
				if err != nil {
					fmt.Printf("Error: %s", err)
				}
				unpackFName := strings.ReplaceAll(k, ".gz", "")
				_, err = gunzipFile(k, unpackFName)
				if err != nil {
					fmt.Printf("Error: %s", err)
				}
			}
		}
	}

	}

    os.Exit(0)
}
