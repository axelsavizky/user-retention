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
	inputFile, err := os.Open(csvPath)
	// Not needed
	if err != nil {
		log.Fatalf("Unable to read input file: %s. Error: %s\n ", csvPath, err)
	}
	defer inputFile.Close()

	csvReader := csv.NewReader(inputFile)

	// Put an arbitrary size to the channel, so we can write without waiting someone to read
	recordsChan := make(chan []string, defaultChannelSize)
	go readCsv(csvReader, recordsChan, csvPath)

	userRetention := *userretention.New()
	userRetention = userRetention.ProcessRecords(recordsChan)

	fmt.Println(userRetention.Get())
}

func readCsv(csvReader *csv.Reader, recordsChan chan<- []string, csvPath string) {

	defer close(recordsChan)

	for {
		record, err := csvReader.Read()
		if err == io.EOF {
			break
		}

		// Not needed
		if err != nil {
			log.Fatal("Unable to parse file as CSV for "+csvPath, err)
		}

		recordsChan <- record
	}
}
