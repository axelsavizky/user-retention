package userretention

import (
	"embrace/models"
	"fmt"
	"strings"
	"time"
)

// These could be parameters
const startDateTimestamp int64 = 1609459200
const defaultDateRange int = 14

func Calculate(records [][]string) map[int][]uint {
	result := buildEmptyResult()

	startDate := time.Unix(startDateTimestamp, 0).UTC()
	streakByUserID := make(map[models.UserID]models.UserStreak)

	for _, rawRecord := range records {
		row := models.RowFromRecord(rawRecord)

		// This is actually the number of date in the range. So the first day will be 1, the second day will be 2, etc
		dayToInt := int((row.Date.Sub(startDate).Hours() / 24) + 1)

		userStreak, ok := streakByUserID[row.UserID]
		if !ok {
			// First streak for user
			result[dayToInt][0]++
			streakByUserID[row.UserID] = models.UserStreak{FirstDay: dayToInt, LastDay: dayToInt}
		} else {
			// User already has a streak
			if userStreak.LastDay == dayToInt {
				// Day already counted for this user. It might happen since a user can have many events at the same day

				// No op
				continue
			}

			if userStreak.LastDay == dayToInt-1 {
				// Currently, it has a streak, so we have to sum on the start date
				differenceSinceStart := dayToInt - userStreak.FirstDay

				result[userStreak.FirstDay][differenceSinceStart-1]--
				result[userStreak.FirstDay][differenceSinceStart]++

				streakByUserID[row.UserID] = models.UserStreak{FirstDay: userStreak.FirstDay, LastDay: dayToInt}
			} else {
				// New streak
				result[dayToInt][0]++
				streakByUserID[row.UserID] = models.UserStreak{FirstDay: dayToInt, LastDay: dayToInt}
			}
		}
	}

	return result
}

func PrintOutput(output map[int][]uint) {
	for i := 1; i <= len(output); i++ {
		fmt.Printf("%d,%s\n", i, sliceToSingleString(output[i]))
	}
}

func sliceToSingleString(s []uint) string {
	// This function joins all the element from the slice with a comma
	// Ex: [1,2,3] => "1,2,3"
	return strings.Trim(strings.Join(strings.Fields(fmt.Sprint(s)), ","), "[]")
}

func buildEmptyResult() map[int][]uint {
	result := make(map[int][]uint)

	for i := 1; i <= defaultDateRange; i++ {
		result[i] = make([]uint, defaultDateRange)
	}

	return result
}
