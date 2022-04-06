package userretention

import (
	"embrace/models"
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"math/rand"
	"strconv"
	"testing"
	"time"
)

func TestNewStreaksByDay(t *testing.T) {
	// Given
	testCases := []struct {
		name           string
		givenDateRange uint
	}{
		{"empty date range", 0},
		{"single date range", 1},
		{"default date range", defaultDateRange},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// When
			streaksByDay := *newStreaksByDay(tc.givenDateRange)

			// Then
			require.Equal(t, int(tc.givenDateRange), len(streaksByDay), "It should have same entries as days in range")

			emptySlice := make([]uint, tc.givenDateRange)
			for i := 1; i <= int(tc.givenDateRange); i++ {
				assert.Equal(t, emptySlice, streaksByDay[i], "Day %d should have an empty uint slice of length %d", i, tc.givenDateRange)
			}
		})
	}
}

func TestToString(t *testing.T) {
	// Given
	randomStreak := uint(rand.Int())

	singleDateRange := generateStreaksWithRandomValue(1, randomStreak)

	fiveDaysRange := generateStreaksWithRandomValue(5, randomStreak)
	fiveDaysExpectedString := fmt.Sprintf("1,%[1]d,%[1]d,%[1]d,%[1]d,%[1]d\n2,0,%[1]d,%[1]d,%[1]d,%[1]d\n3,0,0,%[1]d,%[1]d,%[1]d\n4,0,0,0,%[1]d,%[1]d\n5,0,0,0,0,%[1]d\n", randomStreak)

	testCases := []struct {
		name     string
		given    streaksByDay
		expected string
	}{
		{"empty date range", *newStreaksByDay(0), ""},
		{"single day range", singleDateRange, fmt.Sprintf("1,%d\n", randomStreak)},
		{"five days range", fiveDaysRange, fiveDaysExpectedString},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			assert.Equal(t, tc.expected, tc.given.ToString())
		})
	}
}

func TestProcessRecords(t *testing.T) {
	// Given
	singleUserActivity := [][]models.UserID{
		{1},
	}
	singleUserStreak := *newStreaksByDay(defaultDateRange)
	singleUserStreak[1][0] = 1

	twoDaysActivity := [][]models.UserID{
		{1}, {1},
	}
	twoDaysStreak := *newStreaksByDay(defaultDateRange)
	twoDaysStreak[1][1] = 1

	twoDaysDifferentUsersActivity := [][]models.UserID{
		{1, 2}, {1},
	}

	twoDaysDifferentUsersStreak := *newStreaksByDay(defaultDateRange)
	twoDaysDifferentUsersStreak[1][0] = 1
	twoDaysDifferentUsersStreak[1][1] = 1

	twoDaysUserAppearManyTimesActivity := [][]models.UserID{
		{1, 2, 1, 2, 1, 2, 1, 2}, {1, 1, 1, 1, 1, 1, 1},
	}

	manyStreaksForSameUserActivity := [][]models.UserID{
		{1, 2}, {1}, {2},
	}
	manyStreaksForSameUser := *newStreaksByDay(defaultDateRange)
	manyStreaksForSameUser[1][0] = 1
	manyStreaksForSameUser[1][1] = 1
	manyStreaksForSameUser[3][0] = 1

	moreDaysMoreUsersActivity := [][]models.UserID{
		{1, 2, 3, 4, 1}, {1, 3}, {1, 2, 3}, {1, 2}, {1, 5},
	}
	moreDaysMoreUsersStreak := *newStreaksByDay(defaultDateRange)
	moreDaysMoreUsersStreak[1][0] = 2
	moreDaysMoreUsersStreak[1][2] = 1
	moreDaysMoreUsersStreak[1][4] = 1
	moreDaysMoreUsersStreak[3][1] = 1
	moreDaysMoreUsersStreak[5][0] = 1

	testCases := []struct {
		name         string
		givenRecords [][]string
		expected     streaksByDay
	}{
		{"empty date range", [][]string{}, *newStreaksByDay(defaultDateRange)},
		{"single day with single user", createRecords(singleUserActivity), singleUserStreak},
		{"two days same user", createRecords(twoDaysActivity), twoDaysStreak},
		{"two days different users", createRecords(twoDaysDifferentUsersActivity), twoDaysDifferentUsersStreak},
		{"user appears many times", createRecords(twoDaysUserAppearManyTimesActivity), twoDaysDifferentUsersStreak},
		{"many streaks for same user", createRecords(manyStreaksForSameUserActivity), manyStreaksForSameUser},
		{"more days more users", createRecords(moreDaysMoreUsersActivity), moreDaysMoreUsersStreak},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			recordsChannel := make(chan []string, 10)
			userRetention := *New()
			go func() {
				defer close(recordsChannel)
				for _, givenRecord := range tc.givenRecords {
					recordsChannel <- givenRecord
				}
			}()
			userRetention = userRetention.ProcessRecords(recordsChannel)
			assert.Equal(t, tc.expected, userRetention.streaksByDay)
		})
	}
}

// usersByDay is a slice with the users active per day. It is sorted ascendant.
func createRecords(usersByDay [][]models.UserID) [][]string {
	records := make([][]string, 0)

	startDate := time.Unix(startDateTimestamp, 0).UTC()
	for index, userIDs := range usersByDay {
		currentDay := startDate.AddDate(0, 0, index).Unix()

		for _, userID := range userIDs {
			// We add a number just to ensure the hour, minute and seconds are not useful and we don't take it into account
			stringTimestamp := strconv.FormatInt(currentDay+20000, 10)
			stringUserID := strconv.FormatInt(int64(userID), 10)
			records = append(records, []string{stringTimestamp, stringUserID})
		}
	}

	return records
}

// every day will have the same streak in each possible day
// ex.: day 2 can not have a streak on day 1
func generateStreaksWithRandomValue(dateRange, randomStreak uint) streaksByDay {
	streaksByDay := streaksByDay{}
	for i := 1; i <= int(dateRange); i++ {
		streaks := make([]uint, dateRange)
		for j := i - 1; j < int(dateRange); j++ {
			streaks[j] = randomStreak
		}

		streaksByDay[i] = streaks
	}

	return streaksByDay
}
