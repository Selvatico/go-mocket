package go_mocket

import (
	"database/sql"
	"log"
	"testing"
)

var DB *sql.DB

func GetUsers(db *sql.DB) []map[string]string {
	var res []map[string]string
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
		row := map[string]string{"name": name, "age": age}
		res = append(res, row)
	}
	if err := rows.Err(); err != nil {
		log.Fatal(err)
	}
	return res
}

func GetUsersWithError(db *sql.DB) error {
	age := 27
	_, err := db.Query("SELECT name FROM users WHERE age=?", age)
	return err
}

func CreateUsersWithError(db *sql.DB) error {
	age := 27
	_, err := db.Query("INSERT INTO users (age) VALUES (?) ", age)
	return err
}

func InsertRecord(db *sql.DB) int64 {
	res, err := db.Exec(`INSERT INTO foo VALUES("bar", ?))`, "value")
	if err != nil {
		return 0
	}
	id, _ := res.LastInsertId()
	return id
}

func TestResponses(t *testing.T) {
	Catcher.Register()
	db, _ := sql.Open(DRIVER_NAME, "connection_string") // Could be any connection string
	DB = db
	commonReply := []map[string]interface{}{{"name": "FirstLast", "age": "30"}}

	t.Run("Simple SELECT caught by query", func(t *testing.T) {
		Catcher.Logging = false
		Catcher.Reset().NewMock().WithQuery(`SELECT name FROM users WHERE`).WithReply(commonReply)
		result := GetUsers(DB)
		if len(result) != 1 {
			t.Errorf("Returned sets is not equal to 1. Received %d", len(result))
		}
		if result[0]["age"] != "30" {
			t.Errorf("Age is not equal. Got %v", result[0]["age"])
		}
	})

	t.Run("Simple SELECT with direct object", func(t *testing.T) {
		t.Run("Not a once", func(t *testing.T) {
			Catcher.Reset()
			Catcher.Attach([]*FakeResponse{
				{
					Pattern:  "SELECT name FROM users WHERE",
					Response: commonReply,
					Once:     false,
				},
			})
			result := GetUsers(DB)
			if len(result) != 1 {
				t.Errorf("Returned sets is not equal to 1. Received %d", len(result))
			}
			if result[0]["age"] != "30" {
				t.Errorf("Age is not equal. Got %v", result[0]["age"])
			}
		})

		t.Run("Once", func(t *testing.T) {
			Catcher.Reset()
			Catcher.Attach([]*FakeResponse{
				{
					Pattern:  "SELECT name FROM users WHERE",
					Response: commonReply,
					Once:     true,
				},
			})
			GetUsers(DB)           // Trigger once to use this mock
			result := GetUsers(DB) // trigger second time to receive empty results
			if len(result) != 0 {
				t.Errorf("Returned sets is not equal to 0. Received %d", len(result))
			}
		})
	})

	t.Run("Catch by arguments", func(t *testing.T) {
		Catcher.Reset().NewMock().WithArgs(int64(27)).WithReply(commonReply)
		result := GetUsers(DB)
		if len(result) != 1 {
			t.Fatalf("Returned sets is not equal to 1. Received %d", len(result))
		}
		if result[0]["age"] != "30" {
			t.Errorf("Age is not equal. Got %v", result[0]["age"])
		}
	})

	t.Run("Exceptions and Errors", func(t *testing.T) {
		t.Run("Fire Query error", func(t *testing.T) {
			Catcher.Reset().NewMock().WithArgs(int64(27)).WithReply(commonReply).WithQueryException()
			err := GetUsersWithError(DB)
			if err == nil {
				t.Fatal("Error not triggered")
			}
		})
		t.Run("Fire Execute error", func(t *testing.T) {
			Catcher.Reset().NewMock().WithQuery("INSERT INTO users (age)").WithQueryException()
			err := CreateUsersWithError(DB)
			if err == nil {
				t.Fatal("Error not triggered")
			}
		})
		t.Run("Fire Execute error", func(t *testing.T) {
			Catcher.Reset().NewMock().WithQuery("INSERT INTO users (age)").WithError(sql.ErrNoRows)
			err := CreateUsersWithError(DB)
			if err == nil || err != sql.ErrNoRows {
				t.Fatal("Error not triggered")
			}
		})
	})

	t.Run("Last insert id", func(t *testing.T) {
		var mockedId int64
		mockedId = 64
		Catcher.Reset().NewMock().WithQuery("INSERT INTO foo").WithId(mockedId)
		returnedId := InsertRecord(DB)
		if returnedId != mockedId {
			t.Fatalf("Last insert id not returned. Expected: [%v] , Got: [%v]", mockedId, returnedId)
		}
	})

	t.Run(`Recognise both ? and $1 Postgres placeholders for raw query`, func(t *testing.T) {
		t.Run("Question mark", func(t *testing.T) {
			testFunc := func(db *sql.DB) string {
				var name string
				err := db.QueryRow(`SELECT * FROM foo a = $1 AND b = $2 AND c = $3`, "value", "value2", "value3").Scan(&name)
				if err != nil {
					t.Fatalf("Test function failed [%v]", err)
					return ""
				}
				return name
			}

			Catcher.Reset().NewMock().WithQuery("SELECT * FROM foo ").WithReply([]map[string]interface{}{{"name": "full_name"}})
			returnedName := testFunc(DB)

			if returnedName != "full_name" {
				t.Fatalf("Returned name mismatches. Expected: [%v] , Got: [%v]", "full_name", returnedName)
			}

		})
	})
}
