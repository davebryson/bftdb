package bftdb

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidator(t *testing.T) {
	assert := assert.New(t)
	tests := []string{
		"insert into hello(name) values('dave')",
		"drop table hello",
		"select * from hello",
		"insert into hello value(1)",
	}

	v, e := ValidateSql(tests[0])
	assert.Nil(e)
	assert.Equal(OTHER, v)

	v, e = ValidateSql(tests[1])
	assert.Error(e)
	assert.Equal(DROP, v)

	v, e = ValidateSql(tests[2])
	assert.Nil(e)
	assert.Equal(SELECT, v)

	_, e = ValidateSql(tests[3])
	assert.Error(e)
}

func TestConn(t *testing.T) {
	assert := assert.New(t)
	db, e := NewDb()
	assert.Nil(e)

	s := Statement("CREATE TABLE sample (id INTEGER PRIMARY KEY, name TEXT, data INTEGER)")
	e = db.Write(s)
	assert.Nil(e)

	s = Statement("insert into sample(id, name, data) values(1, 'dave', 10)")
	e = db.Write(s)
	assert.Nil(e)

	s = Statement("insert into sample(id, name, data) values(2, 'bob', 12)")
	e = db.Write(s)
	assert.Nil(e)

	r, e := db.Read("select * from sample")
	assert.Nil(e)
	assert.Nil(e)
	assert.Equal(2, len(r.Values))

	r, e = db.Read("select name from sample")
	assert.Nil(e)
	assert.Nil(e)
	assert.Equal(1, len(r.Columns))
}
