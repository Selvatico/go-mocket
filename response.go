package go_mocket

import (
	"database/sql"
	"database/sql/driver"
	"fmt"
	"log"
	"reflect"
	"strings"
)

const (
	DRIVER_NAME = "MOCK_FAKE_DRIVER"
)

//Catcher is global instance of Catcher used for attaching all mocks to connection
var Catcher *MockCatcher

//MockCatcher is global entity to save all mocks aka FakeResponses
type MockCatcher struct {
	Mocks                []*FakeResponse // Slice of all mocks
	Logging              bool            // Do we need to log what we catching?
	PanicOnEmptyResponse bool            // If not response matches - do we need to panic?
}

// Register safely register FakeDriver
func (mc *MockCatcher) Register() {
	driversList := sql.Drivers()
	for _, name := range driversList {
		if name == DRIVER_NAME {
			return
		}
	}
	sql.Register(DRIVER_NAME, FakeDriver{})
}

// Attach several mocks to MockCather. Could be useful to attach mocks from some factories of mocks
func (mc *MockCatcher) Attach(fr []*FakeResponse) {
	mc.Mocks = append(mc.Mocks, fr...)
}

//FindResponse finds suitable response by provided
func (mc *MockCatcher) FindResponse(query string, args []driver.NamedValue) *FakeResponse {
	if mc.Logging {
		log.Printf("mock_catcher: check query: %s", query)
	}

	for _, resp := range mc.Mocks {
		if resp.IsMatch(query, args) {
			resp.MarkAsTriggered()
			return resp
		}
	}

	if mc.PanicOnEmptyResponse {
		panic(fmt.Sprintf("No responses matches query %s ", query))
	}

	// Let's have always dummy version of response
	return &FakeResponse{
		Response:   make([]map[string]interface{}, 0),
		Exceptions: &Exceptions{},
	}
}

//NewMock creates new FakeResponse and return for chains of attachments
func (mc *MockCatcher) NewMock() *FakeResponse {
	fr := &FakeResponse{Exceptions: &Exceptions{}, Response: make([]map[string]interface{}, 0)}
	mc.Mocks = append(mc.Mocks, fr)
	return fr
}

//Reset removes all Mocks to start process again
func (mc *MockCatcher) Reset() *MockCatcher {
	mc.Mocks = make([]*FakeResponse, 0)
	return mc
}

//Exceptions represents	 possible exceptions during query executions
type Exceptions struct {
	HookQueryBadConnection func() bool
	HookExecBadConnection  func() bool
}

//FakeResponse represents mock of response with holding all required values to return mocked response
type FakeResponse struct {
	Pattern      string                            // SQL query pattern to match with
	Args         []interface{}                     // List args to be matched with
	Response     []map[string]interface{}          // Array of rows to be parsed as result
	Once         bool                              // To trigger only once
	Triggered    bool                              // If it was triggered at least once
	Callback     func(string, []driver.NamedValue) // Callback to execute when response triggered
	RowsAffected int64                             // Defines affected rows count
	LastInsertId int64                             // ID to be returned for INSERT queries
	Error        error                             // Any type of error which could happen dur
	*Exceptions
}

// Returns true either when nothing to compare or deep equal check passed
func (fr *FakeResponse) isArgsMatch(args []driver.NamedValue) bool {
	arguments := make([]interface{}, len(args))
	if len(args) > 0 {
		for index, arg := range args {
			arguments[index] = arg.Value
		}
	}
	return fr.Args == nil || reflect.DeepEqual(fr.Args, arguments)
}

func (fr *FakeResponse) isQueryMatch(query string) bool {
	return fr.Pattern == "" || strings.Contains(query, fr.Pattern)
}

// IsMatch checks if both query and args matcher's return true and if this is Once mock
func (fr *FakeResponse) IsMatch(query string, args []driver.NamedValue) bool {
	if fr.Once && fr.Triggered {
		return false
	}
	return fr.isQueryMatch(query) && fr.isArgsMatch(args)
}

//MarkAsTriggered marks response as executed. For one time catches it will not make this possible to execute anymore
func (fr *FakeResponse) MarkAsTriggered() {
	fr.Triggered = true
}

//WithQuery adds SQL query pattern to match for
func (fr *FakeResponse) WithQuery(query string) *FakeResponse {
	fr.Pattern = query
	return fr
}

// WithArgs attaches Args check for prepared statements
func (fr *FakeResponse) WithArgs(vars ...interface{}) *FakeResponse {
	if len(vars) > 0 {
		fr.Args = make([]interface{}, len(vars))
		for index, v := range vars {
			fr.Args[index] = v
		}
	}
	return fr
}

// WithReply adds to chain and assign some parts of response
func (fr *FakeResponse) WithReply(response []map[string]interface{}) *FakeResponse {
	fr.Response = response
	return fr
}

// OneTime sets current mock to be triggered only once
func (fr *FakeResponse) OneTime() *FakeResponse {
	fr.Once = true
	return fr
}

// WithExecException says that if mock attached to non-SELECT query we need to trigger error there
func (fr *FakeResponse) WithExecException() *FakeResponse {
	fr.Exceptions.HookExecBadConnection = func() bool {
		return true
	}
	return fr
}

// WithQueryException adds to SELECT mocks triggering of error
func (fr *FakeResponse) WithQueryException() *FakeResponse {
	fr.Exceptions.HookQueryBadConnection = func() bool {
		return true
	}
	return fr
}

// WithCallback adds callback to be executed during matching
func (fr *FakeResponse) WithCallback(f func(string, []driver.NamedValue)) *FakeResponse {
	fr.Callback = f
	return fr
}

// WithRowsNum specifies how many records to consider as affected
func (fr *FakeResponse) WithRowsNum(num int64) *FakeResponse {
	fr.RowsAffected = num
	return fr
}

// WithId sets ID to be considered as insert ID for INSERT statements
func (fr *FakeResponse) WithId(id int64) *FakeResponse {
	fr.LastInsertId = id
	return fr
}

// WithError sets Error to FakeResponse struct to have it available on any statements executed
// example: WithError(sql.ErrNoRows)
func (fr *FakeResponse) WithError(err error) *FakeResponse {
	fr.Error = err
	return fr
}

func init() {
	Catcher = &MockCatcher{}
}
