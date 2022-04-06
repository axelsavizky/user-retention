package userretention

import (
	"embrace/models"
	"fmt"
	"strings"
	"time"
)

// These could be parameters
const startDateTimestamp int64 = 1609459200
const defaultDateRange uint = 14

// StreaksByDay is a structure to store the different streaks grouped by date start.
// Key: dayNumber (firstDay = 1, secondDay = 2, etc...) which represents the day that the streak stared.
// Value: slice which each position represents a streak that took i+1 days (with i the slice index).
type StreaksByDay map[int][]uint

// NewStreaksByDay returns a StreaksByDay with all the positions initialized and all the streaks with 0.
// It will have dateRange positions and each streak will be of dateRange length as well.
func NewStreaksByDay(dateRange uint) *StreaksByDay {
	streaksByDay := make(StreaksByDay)

	for i := 1; i <= int(dateRange); i++ {
		streaksByDay[i] = make([]uint, dateRange)
	}

	return &streaksByDay
}

// ToString converts a map of streaks by day to a printable string separated by commas.
// Each line will be something like: key, value[0], value[1], value[2], ..., etc.
// It will be sorted by key ascendant
func (streaksByDay StreaksByDay) ToString() string {
	stringBuilder := strings.Builder{}
	for i := 1; i <= len(streaksByDay); i++ {
		streaksByDayFormatted := fmt.Sprintf("%d,%s\n", i, sliceToSingleString(streaksByDay[i]))
		stringBuilder.WriteString(streaksByDayFormatted)
	}

	return stringBuilder.String()
}

func sliceToSingleString(s []uint) string {
	// This function joins all the element from the slice with a comma
	// Ex: [1,2,3] => "1,2,3"
	return strings.Trim(strings.Join(strings.Fields(fmt.Sprint(s)), ","), "[]")
}

// Calculate calculates the user retention of many users grouped by day.
// It takes a slice of records as input. Each record should have the timestamp on the first position and the userID on the second position.
// The return value is explained on the StreaksByDay documentation.
func Calculate(records [][]string) StreaksByDay {
	streaksByDay := *NewStreaksByDay(defaultDateRange)

	startDate := time.Unix(startDateTimestamp, 0).UTC()
	streakByUserID := make(map[models.UserID]models.UserStreak)

	for _, rawRecord := range records {
		row := models.RowFromRecord(rawRecord)

		// This is actually the number of date in the range. So the first day will be 1, the second day will be 2, etc
		dayToInt := int((row.Date.Sub(startDate).Hours() / 24) + 1)

		userStreak, ok := streakByUserID[row.UserID]
		if !ok {
			// First streak for user
			streaksByDay[dayToInt][0]++
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

				streaksByDay[userStreak.FirstDay][differenceSinceStart-1]--
				streaksByDay[userStreak.FirstDay][differenceSinceStart]++

				streakByUserID[row.UserID] = models.UserStreak{FirstDay: userStreak.FirstDay, LastDay: dayToInt}
			} else {
				// NewStreaksByDay streak
				streaksByDay[dayToInt][0]++
				streakByUserID[row.UserID] = models.UserStreak{FirstDay: dayToInt, LastDay: dayToInt}
			}
		}
	}

	return streaksByDay
}
