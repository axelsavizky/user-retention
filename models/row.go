package models

import (
	"strconv"
	"time"
)

type Row struct {
	Date   time.Time
	UserID UserID
}

func RowFromRecord(record []string) Row {
	// Ignore errors since we don't want to handle malformed data
	timestamp, _ := strconv.ParseInt(record[0], 10, 32)
	userID, _ := strconv.ParseInt(record[1], 10, 64)

	dateWithTime := time.Unix(timestamp, 0).UTC()
	date := time.Date(dateWithTime.Year(), dateWithTime.Month(), dateWithTime.Day(), 0, 0, 0, 0, dateWithTime.Location())

	return Row{date, UserID(userID)}
}
