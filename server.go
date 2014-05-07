package main

import (
	"encoding/json"
	"github.com/gocql/gocql"
	"net/http"
)

// GET /playthroughs
func createAction(w http.ResponseWriter, r *http.Request, session *gocql.Session) {
	if r.Method != "GET" {
		http.NotFound(w, r)
		return
	}

	var message PlaythroughJSONRepr
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

// GET /topfriends
func topFriendsAction(w http.ResponseWriter, r *http.Request, session *gocql.Session) {
	if r.Method != "GET" {
		http.NotFound(w, r)
		return
	}

	user, err := FindUser(session, "player_a") // implement authentication or make it an arg?
	if err != nil {
		http.Error(w, "User not found", http.StatusNotFound)
	}

	w.Header().Add("Content-Type", "application/json")

	friends := FindTopFriends(user)
	encoder := json.NewEncoder(w)
	encoder.Encode(friends)
}

// imagine this is a cron job instead
// POST /wannabecron
func wannabeCronAction(w http.ResponseWriter, r *http.Request, session *gocql.Session) {
	if r.Method != "POST" {
		http.NotFound(w, r)
		return
	}

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
