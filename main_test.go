package main

import (
	"bytes"
	"encoding/csv"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestReadCsv(t *testing.T) {
	var buffer bytes.Buffer

	buffer.WriteString("10,1\n20,2\n30,3")
	csvReader := csv.NewReader(&buffer)

	recordsChan := make(chan []string, 10)

	readCsv(csvReader, recordsChan, "aPath")

	expectedRecords := [][]string{
		{"10", "1"}, {"20", "2"}, {"30", "3"},
	}

	for _, expectedRecord := range expectedRecords {
		record := <-recordsChan

		assert.Equal(t, expectedRecord, record)
	}

}
