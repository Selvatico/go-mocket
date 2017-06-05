package go_mocket

import (
	"database/sql/driver"
)

// FakeResult implementation of sql Result interface
type FakeResult struct {
	insertID     int64
	rowsAffected int64
}

// NewFakeResult returns result interface instance
func NewFakeResult(insertId int64, rowsAffected int64) driver.Result {
	return &FakeResult{insertId, rowsAffected}
}

// LastInsertId required to give sql package ability get ID of inserted record
func (fr *FakeResult) LastInsertId() (int64, error) {
	return fr.insertID, nil
}

//  RowsAffected returns the number of rows affected
func (fr *FakeResult) RowsAffected() (int64, error) {
	return fr.rowsAffected, nil
}
