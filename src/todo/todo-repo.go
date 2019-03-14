package todo

import (
	"encoding/json"
	"github.com/gocql/gocql"
	"log"
	"net/http"
)

const (
	LOG_ERROR = "Error al guardar por: "
)
const (
	SELECT = "SELECT id, text FROM todo"
	INSERT = "INSERT INTO todo (id, text) VALUES (?, ?)"
)

func GetTodo(writer http.ResponseWriter, request *http.Request, session *gocql.Session) {
	todo := findAll(session)
	json.NewEncoder(writer).Encode(&todo)
}

func PostTodo(writer http.ResponseWriter, request *http.Request, session *gocql.Session) {
	var t Todo
	json.NewDecoder(request.Body).Decode(&t)
	save(session, &t)
}

func findAll(session *gocql.Session) []Todo {
	var ts []Todo
	var t Todo
	it := session.Query(SELECT).Iter()
	for it.Scan(&t.ID, &t.Name) {
		ts = append(ts, t)
	}
	if err := it.Close(); err != nil {
		log.Println(LOG_ERROR, err)
	}
	return ts
}

func save(session *gocql.Session, todo *Todo) {
	if err := session.Query(INSERT,
		gocql.TimeUUID(), todo.Name).Exec(); err != nil {
		log.Println(LOG_ERROR, err)
	}
}
