package bftdb

import (
	"database/sql"
	"fmt"

	sqlite3 "github.com/mattn/go-sqlite3"
	"github.com/xwb1989/sqlparser"
)

const (
	SELECT = 0
	DROP   = 1
	OTHER  = 2
	BAD    = 3
)

// ValidateSql checks the incoming statements  for correctness.
// 'drop' statements are not allowed
func ValidateSql(val string) (int, error) {
	stmt, err := sqlparser.Parse(val)
	if err != nil {
		return BAD, fmt.Errorf("SQL Error in %s", val)
	}
	switch smt := stmt.(type) {
	case *sqlparser.DDL:
		if smt.Action == "drop" {
			return DROP, fmt.Errorf("DROP statements are not allowed")
		}
	case *sqlparser.Select:
		return SELECT, nil
	}
	return OTHER, nil
}

type QueryResult struct {
	Columns []string        `json:"columns,omitempty"`
	Values  [][]interface{} `json:"values,omitempty"`
}

type DbWrapper struct {
	db *sql.DB
}

// NewDb create new inmemory DB
func NewDb() (*DbWrapper, error) {
	sql.Register("bftdb", &sqlite3.SQLiteDriver{})
	db, err := sql.Open("bftdb", "")
	if err != nil {
		return nil, err
	}

	return &DbWrapper{db}, nil
}

func (self *DbWrapper) Close() {
	self.db.Close()
}

func (self *DbWrapper) Write(stmt Statement) error {
	val := stmt.String()
	fmt.Printf("statment %s\n", val)
	if _, err := ValidateSql(val); err != nil {
		return err
	}
	_, e := self.db.Exec(val, nil)
	return e
}

func (self *DbWrapper) Read(query string) (*QueryResult, error) {
	if _, err := ValidateSql(query); err != nil {
		return nil, err
	}
	result := &QueryResult{}
	rs, err := self.db.Query(query, nil)
	defer func() {
		if rs != nil {
			rs.Close()
		}
	}()
	if err != nil {
		return nil, err
	}

	colNames, err := rs.Columns()
	if err != nil {
		return nil, err
	}
	result.Columns = colNames

	// Make a slice for the values
	values := make([]interface{}, len(result.Columns))
	scanArgs := make([]interface{}, len(values))
	for i := range values {
		scanArgs[i] = &values[i]
	}

	for rs.Next() {
		rowValues := make([]interface{}, 0)
		err = rs.Scan(scanArgs...)
		if err != nil {
			return nil, err
		}
		for _, value := range values {
			switch value.(type) {
			case nil:
			case []byte:
				rowValues = append(rowValues, string(value.([]byte)))
			default:
				rowValues = append(rowValues, value)
			}
		}
		result.Values = append(result.Values, rowValues)
	}

	return result, nil
}
