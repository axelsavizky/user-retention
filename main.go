package main

import (
	"embrace/userretention"
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"os"
)

const defaultChannelSize = 100

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

	// Put an arbitrary size to the channel so we can write without waiting someone to read
	records := make(chan []string, defaultChannelSize)
	go func() {
		csvReader := csv.NewReader(inputFile)
		defer close(records)

		for {
			record, err := csvReader.Read()
			if err == io.EOF {
				break
			}

			// Not needed
			if err != nil {
				log.Fatal("Unable to parse file as CSV for "+csvPath, err)
			}

			records <- record
		}
	}()

	userRetention := *userretention.New()
	userRetention = userRetention.ProcessRecords(records)

	fmt.Println(userRetention.Get())
}
