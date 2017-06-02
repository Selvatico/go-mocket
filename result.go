package go_mocket

import (
	"database/sql/driver"
)

type FakeResult struct {
	insertID     int64
	rowsAffected int64
}

func NewFakeResult(insertId int64, rowsAffected int64) driver.Result {
	return &FakeResult{insertId, rowsAffected}
}

func (fr *FakeResult) LastInsertId() (int64, error) {
	return fr.insertID, nil
}

func (fr *FakeResult) RowsAffected() (int64, error) {
	return fr.rowsAffected, nil
}
