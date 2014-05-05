package main

import (
	"fmt"
	"github.com/gocql/gocql"
)

type User struct {
	session   *gocql.Session
	UserId    string
	MaxPoints int
}

func NewUser(session *gocql.Session, userId string, maxPoints int) *User {
	return &User{session, userId, maxPoints}
}

func FindUser(session *gocql.Session, userId string) (*User, error) {
	user := new(User)
	user.session = session

	sql := `SELECT user_id, maxpoints FROM users where user_id = ?`
	query := session.Query(sql, userId)
	if err := query.Scan(&user.UserId, &user.MaxPoints); err != nil {
		return nil, err
	}

	return user, nil
}

func findAllFriendsByUserId(session *gocql.Session, userId string) *gocql.Iter {
	sql := `SELECT friend_id FROM friends WHERE user_id = ?`
	return session.Query(sql, userId).Iter()
}

func UpdateFriendsMaxPoints(session *gocql.Session) {
	var user User
	user.session = session

	sql := `SELECT user_id, maxpoints FROM users` // AND updated_at > LAST_RUN
	iter := session.Query(sql).Iter()
	for iter.Scan(&user.UserId, &user.MaxPoints) {
		fmt.Printf("Updating max points for user: %s\n", user.UserId)
		go user.updateFriendsMaxPoints()
	}

	if err := iter.Close(); err != nil {
		panic(err)
	}
}

// friendsIn := strings.Join(friends, ",")
func (self *User) updateFriendsMaxPoints() {
	var friendId string

	iter := self.FriendsIter()
	for iter.Scan(nil, &friendId, nil) {
		sql := `UPDATE friends SET maxpoints = ? WHERE user_id = ? AND friend_id = ?`
		if err := self.session.Query(sql, self.MaxPoints, friendId, self.UserId).Exec(); err != nil {
			panic(err)
		}
	}
}

func (self *User) FriendsIter() *gocql.Iter {
	sql := `SELECT user_id, friend_id, maxpoints FROM friends WHERE user_id = ?`
	return self.session.Query(sql, self.UserId).Iter()
}

func (self *User) Valid() bool {
	return true
}

func (self *User) UpdateMaxPointsIfLarger(points int) {
	if points <= self.MaxPoints {
		return
	}

	sql := `UPDATE users SET maxpoints = ? WHERE user_id = ?`
	if err := self.session.Query(sql, points, self.UserId).Exec(); err != nil {
		fmt.Printf("Failed to save user: %s\n", err.Error())
	}
}

func (self *User) Save() {
	sql := `INSERT INTO users (user_id, maxpoints) VALUES (?, ?, ?)`
	err := self.session.Query(sql, self.UserId, self.MaxPoints).Exec()
	if err != nil {
		fmt.Printf("Failed to save user: %s\n", err.Error())
	}
}
