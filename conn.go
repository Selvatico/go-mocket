package go_mocket

import (
	"context"
	"database/sql/driver"
	"errors"
	"log"
	"regexp"
	"strings"
	"sync"
)

// FakeConn implements connection
type FakeConn struct {
	db     *FakeDB
	currTx *FakeTx // Transaction pointer
	mu     sync.Mutex
	bad    bool
}

func (c *FakeConn) isBad() bool {
	return false
}

// Begin starts and returns a new transaction.
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

// Exec is deprecated
func (c *FakeConn) Exec(query string, args []driver.Value) (driver.Result, error) {
	panic("ExecContext was not called.")
}

// ExecContext is optional to implement and it returns skip
func (c *FakeConn) ExecContext(ctx context.Context, query string, args []driver.NamedValue) (driver.Result, error) {
	return nil, driver.ErrSkip
}

// Query is deprecated
func (c *FakeConn) Query(query string, args []driver.Value) (driver.Rows, error) {
	panic("QueryContext was not called.")
}

// QueryContext is optional
func (c *FakeConn) QueryContext(ctx context.Context, query string, args []driver.NamedValue) (driver.Rows, error) {
	return nil, driver.ErrSkip
}

// Prepare is optional
func (c *FakeConn) Prepare(query string) (driver.Stmt, error) {
	panic("use Prepare")
}

// PrepareContext returns a prepared statement, bound to this connection.
// context is for the preparation of the statement,
// it must not store the context within the statement itself.
func (c *FakeConn) PrepareContext(ctx context.Context, query string) (driver.Stmt, error) {
	var firstStmt = &FakeStmt{q: query, connection: c}
	// Checking how many placeholders do we have
	if strings.Contains(query, "$1") {
		r, err := regexp.Compile(`[$]\d+`)
		if err != nil {
			log.Fatalf(`Cant't compile regexp with err [%v]`, err)
		}
		firstStmt.placeholders = len(strings.Split(r.ReplaceAllString(query, `$$$`), "$$")) - 1 // Postgres notation
	} else {
		firstStmt.placeholders = len(strings.Split(query, "?")) - 1 // Postgres notation
	}

	queryParts := strings.Split(query, " ") // By First statement define the query type
	firstStmt.command = strings.ToUpper(queryParts[0])
	return firstStmt, nil
}
