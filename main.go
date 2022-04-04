package main

import (
	"embrace/models"
	"encoding/csv"
	"fmt"
	"log"
	"os"
)

func main() {
	if len(os.Args) != 2 {
		log.Fatalln("It should receive one argument which is the path to the input CSV")
	}

	csvPath := os.Args[1]
	// File will not fit in memory. TODO: Stream it!!! https://medium.com/swlh/processing-16gb-file-in-seconds-go-lang-3982c235dfa2
	inputFile, err := os.Open(csvPath)
	// Not needed
	if err != nil {
		log.Fatalf("Unable to read input file: %s. Error: %s\n ", csvPath, err)
	}
	defer inputFile.Close()

	csvReader := csv.NewReader(inputFile)
	records, err := csvReader.ReadAll()
	// Not needed
	if err != nil {
		log.Fatal("Unable to parse file as CSV for "+csvPath, err)
	}

	dayToUser := make(map[string]models.Set)

	for _, rawRecord := range records {
		row := models.RowFromRecord(rawRecord)

		if dayToUser[row.ToKey()] == nil {
			dayToUser[row.ToKey()] = make(models.Set)
		}
		set := dayToUser[row.ToKey()]
		set.Add(row.UserID)
	}

	fmt.Println(dayToUser)
}
