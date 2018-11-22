[![GoDoc](https://godoc.org/github.com/Selvatico/go-mocket?status.svg)](https://godoc.org/github.com/Selvatico/go-mocket)  [![Build Status](https://travis-ci.org/Selvatico/go-mocket.svg?branch=master)](https://travis-ci.org/Selvatico/go-mocket) [![Go Report Card](https://goreportcard.com/badge/github.com/Selvatico/go-mocket)](https://goreportcard.com/report/github.com/Selvatico/go-mocket)

# Go-Mocket – Go GORM & SQL mocking library

Go-Mocket is a library inspired by [DATA-DOG/go-sqlmock](https://github.com/DATA-DOG/go-sqlmock).

As an inspiration library, it is the implementation of [sql/driver](https://godoc.org/database/sql/driver) interface but at the same time it follows a different approach and only has a similar API.
This library helps to mock any DB connection with [jinzhu/gorm](https://github.com/jinzhu/gorm), as it was the created to serve this purpose.

List of features in the library:

* Mock `INSERT`, `UPDATE`, `SELECT`, `DELETE`
* Support for transactions
* 2 API's to use - `chaining` and via specifying a whole mock object
* Matching by prepared statements arguments
* You don't require to change anything inside your code to start using this library
* Ability to trigger exceptions
* Attach callbacks to mocked responses to add an additional check or modify a response

**NOTE**, Please be aware that driver catches SQL without DB specifics. Generation of queries is done by *SQL* package

## Install

```
go get github.com/Selvatico/go-mocket
```

## Usage

There are two possible ways to use `mocket`:

* Chaining API
* Specifying `FakeResponse` object with all fields manually. Could be useful for cases when mocks stored separately as the list of FakeResponses. 

### Enabling driver

Somewhere in your code, do this to set up a tests

```go
import (
    "database/sql"
    mocket "github.com/Selvatico/go-mocket"
    "github.com/jinzhu/gorm"
)

func SetupTests() {
    mocket.Catcher.Register()
    // GORM
    db, err := gorm.Open(mocket.DriverName, "any_string") // Could be any connection string
    app.DB = db // Assumption that it will be used everywhere the same
    //OR 
    // Regular sql package usage
    db, err := sql.Open(mocket.DriverName, "any_string")
}
```

Now, if you use a singleton instance of DB, it will use a mocked connection everywhere.

### Chain usage

#### Example of mocking by pattern

```go
import mocket "github.com/Selvatico/go-mocket"
import "net/http/httptest"

func TestHandler(t *testing.T) {
    request := httptest.NewRequest("POST", "/application", nil)
    recorder := httptest.NewRecorder()

    GlobalMock := mocket.Catcher
    GlobalMock.Logging = true // log mocket behavior

    // field names here mapped to the database schema
    commonReply := []map[string]interface{}{{"some_id": "2", "field": "value"}}
    // Mock only by query pattern
    GlobalMock.NewMock().WithQuery(`"campaigns".name IS NULL AND (("uuid" = test_uuid))`).WithReply(commonReply)
    Post(recorder, request) // call handler

    r := recorder.Result()
    body, _ := ioutil.ReadAll(r.Body)

    // some assertion about results
    //...
}

```

## Documentation


For more documentation please consult the [Documentation](https://github.com/Selvatico/go-mocket/blob/master/DOCUMENTATION.md)

## License

MIT License

Copyright (c) 2017 Seredenko Dmitry

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
