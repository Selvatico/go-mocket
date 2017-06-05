package go_mocket

import (
	"database/sql"
	"log"
	"testing"
)

var DB *sql.DB

func GetUsers(db *sql.DB) []map[string]interface{} {
	var res []map[string]interface{}
	age := 27
	rows, err := db.Query("SELECT name FROM users WHERE age=?", age)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	for rows.Next() {
		var name string
		var age string
		if err := rows.Scan(&name, &age); err != nil {
			log.Fatal(err)
		}
		row := make(map[string]interface{})
		row["name"] = name
		row["age"] = age
		res = append(res, row)
	}
	if err := rows.Err(); err != nil {
		log.Fatal(err)
	}
	return res
}

func TestResponses(t *testing.T) {
	sql.Register("fake_test", FakeDriver{})
	db, _ := sql.Open("fake_test", "connection_string") // Could be any connection string
	DB = db

	t.Run("Simple SELECT caught by query", func(t *testing.T) {
		Catcher.Logging = false
		commonReply := []map[string]interface{}{{"name": "FirstLast", "age": "30"}}
		Catcher.Reset().NewMock().WithQuery(`SELECT name FROM users WHERE`).WithReply(commonReply)
		result := GetUsers(DB)
		if len(result) != 1 {
			t.Errorf("Returned sets is not equal to 1. Received %d", len(result))
		}
		if result[0]["age"] != "30" {
			t.Errorf("Age is not equal. Got %v", result[0]["age"])
		}
	})
}
