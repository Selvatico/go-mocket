package gomocket

/*
Package go_mocket is way to mock DB for GORM and just sql package usage
First you need to activate package somewhere in your tests code like this
	package main

	import (
    	"database/sql"
    	mocket "github.com/Selvatico/go-mocket"
    	"github.com/jinzhu/gorm"
	)

	var DB *gorm.DB

	func SetupTests() {
    	mocket.Catcher.Register()
    	// GORM
    	db, err := gorm.Open("fake_test", "connection_string") // Could be any connection string
    	DB = db

		// OR
    	// Regular sql package usage
    	db, err := sql.Open(driver, source)
	}

	func GetUsers(db *sql.DB) []map[string]interface {} {
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
			//fmt.Printf("%s is %d\n", name, age)
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

Somewhere in you tests:
	package main

	import (
		"log"
		"testing"
		"database/sql"
		mocket "github.com/Selvatico/go-mocket"
	)

	func TestResponses (t *testing.T) {
		mocket.Catcher.Register()
		db, _ := sql.Open("fake_test", "connection_string") // Could be any connection string
		DB = db

		t.Run("Simple SELECT caught by query", func(t *testing.T) {
			mocket.Catcher.Logging = true
			commonReply := []map[string]interface{}{{"name": "FirstLast", "age": "30"}}
			mocket.Catcher.Reset().NewMock().WithQuery(`SELECT name FROM users WHERE`).WithReply(commonReply)
			result := GetUsers(DB)
			if len(result) != 1 {
				t.Errorf("Returned sets is not equal to 1. Received %d", len(result))
			}
			if result[0]["age"] != "30" {
				t.Errorf("Age is not equal. Got %v", result[0]["age"])
			}
		})
	}

For more information and use cases please check: https://github.com/Selvatico/go-mocket

*/
