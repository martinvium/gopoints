package main

import (
	"fmt"
	"github.com/gocql/gocql"
)

type Playthrough struct {
	session   *gocql.Session
	userId    string
	timestamp gocql.UUID
	points    int
}

func NewPlaythrough(session *gocql.Session, userId string, points int) *Playthrough {
	return &Playthrough{session, userId, gocql.TimeUUID(), points}
}

func (self *Playthrough) Valid() bool {
	return true
}

func (self *Playthrough) Save() {
	sql := `INSERT INTO playthroughs (user_id, timestamp, points) VALUES (?, ?, ?)`
	err := self.session.Query(sql,
		self.userId,
		self.timestamp,
		self.points).Exec()

	if err != nil {
		fmt.Printf("Failed to insert playthrough: %s\n", err.Error())
	}
}
