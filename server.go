package main

import (
	"encoding/json"
	"github.com/gocql/gocql"
	"net/http"
)

type PlaythroughMessage struct {
	UserId string
	Points int
}

// TODO error if content type is not json
func createAction(w http.ResponseWriter, r *http.Request, session *gocql.Session) {
	var message PlaythroughMessage
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&message); err != nil {
		http.Error(w, "Failed to parse json from body", http.StatusBadRequest)
		return
	}

	playthrough := NewPlaythrough(session, message.UserId, message.Points)
	if playthrough.valid() {
		go playthrough.save()
	} else {
		http.Error(w, "Playthrough was not valid", http.StatusBadRequest)
	}
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

	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		panic("ListenAndServe: " + err.Error())
	}
}
