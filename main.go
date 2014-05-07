package main

import (
	"github.com/gocql/gocql"
	"net/http"
)

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

// see actions.go for handlers
func main() {
	session := NewCassandraSession("gopoints")
	defer session.Close()

	http.HandleFunc("/playthroughs", func(w http.ResponseWriter, r *http.Request) {
		CreateAction(w, r, session)
	})

	http.HandleFunc("/topfriends", func(w http.ResponseWriter, r *http.Request) {
		TopFriendsAction(w, r, session)
	})

	http.HandleFunc("/wannabecron", func(w http.ResponseWriter, r *http.Request) {
		WannabeCronAction(w, r, session)
	})

	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		panic("ListenAndServe: " + err.Error())
	}
}
