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

type userRetention struct {
	streaksByDay streaksByDay
}

func New() *userRetention {
	return &userRetention{*newStreaksByDay(defaultDateRange)}
}

// streaksByDay is a structure to store the different streaks grouped by date start.
// Key: dayNumber (firstDay = 1, secondDay = 2, etc...) which represents the day that the streak stared.
// Value: slice which each position represents a streak that took i+1 days (with i the slice index).
type streaksByDay map[int][]uint

// newStreaksByDay returns a streaksByDay with all the positions initialized and all the streaks with 0.
// It will have dateRange positions and each streak will be of dateRange length as well.
func newStreaksByDay(dateRange uint) *streaksByDay {
	streaksByDay := make(streaksByDay)

	for i := 1; i <= int(dateRange); i++ {
		streaksByDay[i] = make([]uint, dateRange)
	}

	return &streaksByDay
}

// ToString converts a map of streaks by day to a printable string separated by commas.
// Each line will be something like: key, value[0], value[1], value[2], ..., etc.
// It will be sorted by key ascendant
func (streaksByDay streaksByDay) ToString() string {
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

// ProcessRecords calculates the user retention of many users grouped by day.
// It receives the records through a channel as input. Each record should have the timestamp on the first position and the userID on the second position.
// userRetention is an immutable data structure, so this function returns a new userRetention.
func (userRetention userRetention) ProcessRecords(recordsChan <-chan []string) userRetention {
	startDate := time.Unix(startDateTimestamp, 0).UTC()
	streakByUserID := make(map[models.UserID]models.UserStreak)

	for rawRecord := range recordsChan {
		row := models.RowFromRecord(rawRecord)

		// This is actually the number of date in the range. So the first day will be 1, the second day will be 2, etc
		dayToInt := int((row.Date.Sub(startDate).Hours() / 24) + 1)

		userStreak, ok := streakByUserID[row.UserID]
		if !ok {
			// First streak for user
			userRetention.streaksByDay[dayToInt][0]++
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

				// Decrease current day streak, and add to the next day current streak
				userRetention.streaksByDay[userStreak.FirstDay][differenceSinceStart-1]--
				userRetention.streaksByDay[userStreak.FirstDay][differenceSinceStart]++

				streakByUserID[row.UserID] = models.UserStreak{FirstDay: userStreak.FirstDay, LastDay: dayToInt}
			} else {
				// newStreaksByDay streak
				userRetention.streaksByDay[dayToInt][0]++
				streakByUserID[row.UserID] = models.UserStreak{FirstDay: dayToInt, LastDay: dayToInt}
			}
		}
	}

	return userRetention
}

// Get returns the user retention processed with ProcessRecords in a printable way separated by comma.
func (userRetention userRetention) Get() string {
	return userRetention.streaksByDay.ToString()
}
