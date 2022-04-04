package models

type Set map[int64]struct{}

func (set *Set) Add(userID int64) {
	(*set)[userID] = struct{}{}
}
