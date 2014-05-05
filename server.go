package main

import (
	"encoding/json"
	// "fmt"
	"github.com/gocql/gocql"
	"net/http"
)

type PlaythroughMessage struct {
	UserId string
	Points int
}

type FriendMessage struct {
	FriendId  string
	MaxPoints int
}

// TODO error if content type is not json
func createAction(w http.ResponseWriter, r *http.Request, session *gocql.Session) {
	var message PlaythroughMessage
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&message); err != nil {
		http.Error(w, "Failed to parse json from body", http.StatusBadRequest)
		return
	}

	user, err := FindUser(session, message.UserId)
	if err != nil {
		http.Error(w, "User not found: "+message.UserId, http.StatusNotFound)
		return
	}

	playthrough := NewPlaythrough(session, message.UserId, message.Points)
	if playthrough.Valid() {
		go playthrough.Save()
		go user.UpdateMaxPointsIfLarger(message.Points)
	} else {
		http.Error(w, "Playthrough was not valid", http.StatusBadRequest)
	}
}

// TODO use an array instead of sending multiple seperate entities to the encoder
func topFriendsAction(w http.ResponseWriter, r *http.Request, session *gocql.Session) {
	user, err := FindUser(session, "player_a")
	if err != nil {
		panic(err)
	}

	w.Header().Add("Content-Type", "application/json")
	encoder := json.NewEncoder(w)

	friend := new(FriendMessage)
	iter := user.FriendsIter()
	for iter.Scan(nil, &friend.FriendId, &friend.MaxPoints) {
		encoder.Encode(friend)
	}

	if err := iter.Close(); err != nil {
		panic(err)
	}
}

// imagine this is a cron job instead
func wannabeCronAction(w http.ResponseWriter, r *http.Request, session *gocql.Session) {
	go UpdateFriendsMaxPoints(session)
}

func NewCassandraSession(keyspace string) *gocql.Session {
	cluster := gocql.NewCluster("127.0.0.1")
	cluster.Keyspace = keyspace
	cluster.Consistency = gocql.Quorum
	session, err := cluster.CreateSession()
	if err != nil {
		panic(err)
	}

	return session
}

func main() {
	session := NewCassandraSession("gopoints")
	defer session.Close()

	http.HandleFunc("/playthroughs", func(w http.ResponseWriter, r *http.Request) {
		createAction(w, r, session)
	})

	http.HandleFunc("/topfriends", func(w http.ResponseWriter, r *http.Request) {
		topFriendsAction(w, r, session)
	})

	http.HandleFunc("/wannabecron", func(w http.ResponseWriter, r *http.Request) {
		wannabeCronAction(w, r, session)
	})

	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		panic("ListenAndServe: " + err.Error())
	}
}
