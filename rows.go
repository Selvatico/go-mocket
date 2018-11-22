package gomocket

import (
	"database/sql"
	"database/sql/driver"
	"errors"
	"io"
	"reflect"
	"time"
)

// RowsCursor is implementation of Rows sql interface
type RowsCursor struct {
	cols    []string
	colType [][]string
	posSet  int
	posRow  int
	rows    [][]*row
	closed  bool

	// errPos and err are for making Next return early with error.
	errPos int
	err    error

	bytesClone map[*byte][]byte
}

type row struct {
	cols []interface{} // must be same size as its table colname + coltype
}

// Close closes the rows iterator.
func (rc *RowsCursor) Close() error {
	if !rc.closed {
		for _, bs := range rc.bytesClone {
			bs[0] = 255 // first byte corrupted
		}
	}
	rc.closed = true
	return nil
}

// Columns returns the names of the columns.
func (rc *RowsCursor) Columns() []string {
	return rc.cols
}

// ColumnTypeScanType may be implemented by Rows. It should return
// the value type that can be used to scan types into.
func (rc *RowsCursor) ColumnTypeScanType(index int) reflect.Type {
	return colTypeToReflectType(rc.colType[rc.posSet][index])
}

// Next is called to populate the next row of data into
// the provided slice.
func (rc *RowsCursor) Next(accumulator []driver.Value) error {
	if rc.closed {
		return errors.New("fake_db_driver: cursor is closed")
	}
	rc.posRow++
	if rc.posRow == rc.errPos {
		return rc.err
	}
	if rc.posRow >= len(rc.rows[rc.posSet]) {
		return io.EOF // per interface spec
	}
	for i, v := range rc.rows[rc.posSet][rc.posRow].cols {
		accumulator[i] = v
		if bs, ok := v.([]byte); ok {
			if rc.bytesClone == nil {
				rc.bytesClone = make(map[*byte][]byte)
			}
			clone, ok := rc.bytesClone[&bs[0]]
			if !ok {
				clone = make([]byte, len(bs))
				copy(clone, bs)
				rc.bytesClone[&bs[0]] = clone
			}
			accumulator[i] = clone
		}
	}
	return nil
}

// HasNextResultSet is called at the end of the current result set and
// reports whether there is another result set after the current one.
func (rc *RowsCursor) HasNextResultSet() bool {
	return rc.posSet < len(rc.rows)-1
}

// NextResultSet advances the driver to the next result set even
// if there are remaining rows in the current result set.
func (rc *RowsCursor) NextResultSet() error {
	if rc.HasNextResultSet() {
		rc.posSet++
		rc.posRow = -1
		return nil
	}
	return io.EOF // Per interface spec.
}

func colTypeToReflectType(typ string) reflect.Type {
	switch typ {
	case "bool":
		return reflect.TypeOf(false)
	case "nullbool":
		return reflect.TypeOf(sql.NullBool{})
	case "int32":
		return reflect.TypeOf(int32(0))
	case "string":
		return reflect.TypeOf("")
	case "nullstring":
		return reflect.TypeOf(sql.NullString{})
	case "int64":
		return reflect.TypeOf(int64(0))
	case "nullint64":
		return reflect.TypeOf(sql.NullInt64{})
	case "float64":
		return reflect.TypeOf(float64(0))
	case "nullfloat64":
		return reflect.TypeOf(sql.NullFloat64{})
	case "datetime":
		return reflect.TypeOf(time.Time{})
	}
	panic("invalid fakedb column type of " + typ)
}
