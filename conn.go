package go_mocket

import (
	"context"
	"database/sql/driver"
	"errors"
	"strings"
	"sync"
)

type FakeConn struct {
	db     *FakeDB
	currTx *FakeTx // Transaction pointer
	mu     sync.Mutex
	bad    bool
}

func (c *FakeConn) isBad() bool {
	return false
}

func (c *FakeConn) Begin() (driver.Tx, error) {
	if c.isBad() {
		return nil, driver.ErrBadConn
	}
	if c.currTx != nil {
		return nil, errors.New("already in a transaction")
	}
	c.currTx = &FakeTx{c: c}
	return c.currTx, nil
}

func (c *FakeConn) Close() (err error) {
	c.db = nil
	return nil
}

func (c *FakeConn) Exec(query string, args []driver.Value) (driver.Result, error) {
	panic("ExecContext was not called.")
}

func (c *FakeConn) ExecContext(ctx context.Context, query string, args []driver.NamedValue) (driver.Result, error) {
	return nil, driver.ErrSkip
}

func (c *FakeConn) Query(query string, args []driver.Value) (driver.Rows, error) {
	panic("QueryContext was not called.")
}

// We do
func (c *FakeConn) QueryContext(ctx context.Context, query string, args []driver.NamedValue) (driver.Rows, error) {
	return nil, driver.ErrSkip
}

// Should not be called
func (c *FakeConn) Prepare(query string) (driver.Stmt, error) {
	panic("use Prepare")
}

func (c *FakeConn) PrepareContext(ctx context.Context, query string) (driver.Stmt, error) {
	var firstStmt = &FakeStmt{q: query, connection: c}          // Create statement
	firstStmt.placeholders = len(strings.Split(query, "?")) - 1 // Checking how many placeholders do we have
	queryParts := strings.Split(query, " ")                     // By First statement define the query type
	firstStmt.command = strings.ToUpper(queryParts[0])
	return firstStmt, nil
}
