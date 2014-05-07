package main

import (
	"sort"
)

type FriendJSONRepr struct {
	FriendId  string
	MaxPoints int
}

type Friend struct {
	UserId    string
	FriendId  string
	MaxPoints int
}

func FindTopFriends(user *User) []FriendJSONRepr {
	friends := make([]FriendJSONRepr, 0)

	friend := FriendJSONRepr{}
	iter := user.FriendsIter()
	for iter.Scan(nil, &friend.FriendId, &friend.MaxPoints) {
		friends = append(friends, friend)
	}

	if err := iter.Close(); err != nil {
		panic(err)
	}

	sort.Sort(ByMaxPoints(friends))

	return friends
}

type ByMaxPoints []FriendJSONRepr

func (a ByMaxPoints) Len() int           { return len(a) }
func (a ByMaxPoints) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByMaxPoints) Less(i, j int) bool { return a[i].MaxPoints > a[j].MaxPoints }
