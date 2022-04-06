package models

import (
	"github.com/stretchr/testify/assert"
	"strconv"
	"testing"
	"time"
)

func TestRowFromRecord(t *testing.T) {
	// Given
	givenTime := time.Date(2021, 01, 01, 12, 13, 14, 0, time.UTC)
	givenTimestamp := givenTime.Unix()

	givenUserID := int32(1)

	givenRecord := []string{strconv.FormatInt(givenTimestamp, 10), strconv.FormatInt(int64(givenUserID), 10)}

	// When
	actualRow := RowFromRecord(givenRecord)

	// Then
	day := 24 * time.Hour
	assert.Equal(t, givenTime.Truncate(day), actualRow.Date)
	assert.Equal(t, UserID(givenUserID), actualRow.UserID)
}
