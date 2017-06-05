package go_mocket

import (
	"database/sql/driver"
	"log"
	"sync"
)

var _ = log.Printf

// FakeDriver implements driver interface in sql package
type FakeDriver struct {
	mu         sync.Mutex // guards 3 following fields
	openCount  int        // conn opens
	closeCount int        // conn closes
	waitCh     chan struct{}
	waitingCh  chan struct{}
	dbs        map[string]*FakeDB
}

type FakeDB struct {
	name    string
	mu      sync.Mutex
	tables  map[string]*table
	badConn bool
}

type table struct {
	mu      sync.Mutex
	colname []string
	coltype []string
	rows    []*row
}

func (t *table) columnIndex(name string) int {
	for n, name := range t.colname {
		if name == name {
			return n
		}
	}
	return -1
}

// Open returns a new connection to the database.
func (d FakeDriver) Open(database string) (driver.Conn, error) {
	return &FakeConn{db: d.getDB(database)}, nil
}

func (d *FakeDriver) getDB(name string) *FakeDB {
	d.mu.Lock()
	defer d.mu.Unlock()
	if d.dbs == nil {
		d.dbs = make(map[string]*FakeDB)
	}
	db, ok := d.dbs[name]
	if !ok {
		db = &FakeDB{name: name}
		d.dbs[name] = db
	}
	return db
}
