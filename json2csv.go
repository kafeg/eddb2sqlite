package main

import (
	"encoding/csv"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"reflect"
	"sort"
	"strconv"
)

var (
	inputFile   = flag.String("i", "", "/path/to/input.json (optional; default is stdin)")
	outputFile  = flag.String("o", "", "/path/to/output.json (optional; default is stdout)")
	outputDelim = flag.String("d", ",", "delimiter used for output values")
	verbose     = flag.Bool("v", false, "verbose output (to stderr)")
	showVersion = flag.Bool("version", false, "print version string")
	printHeader = flag.Bool("p", false, "prints header to output")
	//keys        = StringArray{}
)

func reverseStrings(input []string) []string {
	if len(input) == 0 {
		return input
	}
	return append(reverseStrings(input[1:]), input[0])
}

func json2csv(inFile, outFile string) {
	//for _, v := range fields {
	//	fmt.Println(v)
	//}
	fmt.Printf("Converting file '%s' to %s\n", inFile, outFile)

	fields := []string{}

	jsonObj := make([]map[string]interface{}, 0)

	data, err := ioutil.ReadFile(inFile)
	if err != nil {
		fmt.Println(err)
	}
	//fmt.Println(string(data))

	err = json.Unmarshal(data, &jsonObj)
	if err != nil {
		fmt.Println(err)
	}
	//fmt.Println(jsonObj)
	f, err := os.Create(outFile)
	if err != nil {
		fmt.Println(err)
	}
	defer f.Close()
	w := csv.NewWriter(f)

	// fill fields arr
	for _, v := range jsonObj {
		for k, _ := range v {
			//dirty
			if k == "states" || k == "minor_faction_presences" || k == "import_commodities" || k == "export_commodities" || k == "prohibited_commodities" || k == "economies" || k == "selling_ships" || k == "selling_modules" || k == "settlement_security" {
				continue
			}
			fields = append(fields, k)
		}
		break;
	}
	sort.Strings(fields)
	fields = reverseStrings(fields)

	//fmt.Println(fields)

	record := make([]string, 0)
	for _, v := range fields {
		record = append(record, v)
	}

	if err := w.Write(record); err != nil {
		log.Fatalln("error writing record to csv:", err)
	}

	for _, v := range jsonObj {
		record = make([]string, 0)
		for _, v2 := range fields {
			if value, ok := v[v2]; ok == true {
				//fmt.Println(" ")
				if value == nil {
					value = "0" // like '"max_buy_price": null'
				}
				//fmt.Println(value, v2, reflect.TypeOf(value).String())
				switch reflect.TypeOf(value).String() {
				case "string":
					record = append(record, value.(string))
				case "bool":
					if (value.(bool)) {
						record = append(record, "1")
					} else {
						record = append(record, "0")
					}
				case "float64":
					record = append(record, strconv.Itoa(int(value.(float64))))
				case "map[string]interface {}":
					fmt.Println(value, v2, reflect.TypeOf(value).String())
					record = append(record, value.(map[string]interface{})["name"].(string))
				default:
					log.Fatalln("Unhandled field type",  reflect.TypeOf(value).String())
				}
			}

		}
		if err := w.Write(record); err != nil {
			log.Fatalln("error writing record to csv:", err)
		}
	}
	w.Flush()
	if err := w.Error(); err != nil {
		log.Fatal(err)
	}
}
