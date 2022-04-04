package models

import (
	"strconv"
	"time"
)

const keyFormat string = "2006-01-02"

type Row struct {
	Timestamp time.Time
	UserID    int64
}

func RowFromRecord(record []string) Row {
	// Ignore errors since we don't want to handle malformed data
	timestamp, _ := strconv.ParseInt(record[0], 10, 32)
	userID, _ := strconv.ParseInt(record[1], 10, 64)

	dateWithTime := time.Unix(timestamp, 0).UTC()
	date := time.Date(dateWithTime.Year(), dateWithTime.Month(), dateWithTime.Day(), 0, 0, 0, 0, dateWithTime.Location())

	return Row{date, userID}
}

func (row Row) ToKey() string {
	return row.Timestamp.Format(keyFormat)
}
