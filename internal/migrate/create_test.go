// Copyright (C) 2019 Storj Labs, Inc.
// See LICENSE for copying information.

package migrate_test

import (
	"database/sql"
	"strconv"
	"testing"

	_ "github.com/lib/pq"
	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"storj.io/storj/internal/dbutil/pgutil"
	"storj.io/storj/internal/dbutil/pgutil/pgtest"
	"storj.io/storj/internal/migrate"
)

func TestCreate_Sqlite(t *testing.T) {
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Fatal(err)
	}
	defer func() { assert.NoError(t, db.Close()) }()

	// should create table
	err = migrate.Create("example", &sqliteDB{db, "CREATE TABLE example_table (id text)"})
	assert.NoError(t, err)

	// shouldn't create a new table
	err = migrate.Create("example", &sqliteDB{db, "CREATE TABLE example_table (id text)"})
	assert.NoError(t, err)

	// should fail, because schema changed
	err = migrate.Create("example", &sqliteDB{db, "CREATE TABLE example_table (id text, version int)"})
	assert.Error(t, err)

	// should fail, because of trying to CREATE TABLE with same name
	err = migrate.Create("conflict", &sqliteDB{db, "CREATE TABLE example_table (id text, version int)"})
	assert.Error(t, err)
}

func TestCreate_Postgres(t *testing.T) {
	if *pgtest.ConnStr == "" {
		t.Skipf("postgres flag missing, example:\n-postgres-test-db=%s", pgtest.DefaultConnStr)
	}

	schema := "create-" + pgutil.CreateRandomTestingSchemaName(8)

	db, err := sql.Open("postgres", pgutil.ConnstrWithSchema(*pgtest.ConnStr, schema))
	if err != nil {
		t.Fatal(err)
	}
	defer func() { assert.NoError(t, db.Close()) }()

	require.NoError(t, pgutil.CreateSchema(db, schema))
	defer func() { assert.NoError(t, pgutil.DropSchema(db, schema)) }()

	// should create table
	err = migrate.Create("example", &postgresDB{db, "CREATE TABLE example_table (id text)"})
	assert.NoError(t, err)

	// shouldn't create a new table
	err = migrate.Create("example", &postgresDB{db, "CREATE TABLE example_table (id text)"})
	assert.NoError(t, err)

	// should fail, because schema changed
	err = migrate.Create("example", &postgresDB{db, "CREATE TABLE example_table (id text, version integer)"})
	assert.Error(t, err)

	// should fail, because of trying to CREATE TABLE with same name
	err = migrate.Create("conflict", &postgresDB{db, "CREATE TABLE example_table (id text, version integer)"})
	assert.Error(t, err)
}

type sqliteDB struct {
	*sql.DB
	schema string
}

func (db *sqliteDB) Rebind(s string) string { return s }
func (db *sqliteDB) Schema() string         { return db.schema }

type postgresDB struct {
	*sql.DB
	schema string
}

func (db *postgresDB) Rebind(sql string) string {
	out := make([]byte, 0, len(sql)+10)

	j := 1
	for i := 0; i < len(sql); i++ {
		ch := sql[i]
		if ch != '?' {
			out = append(out, ch)
			continue
		}

		out = append(out, '$')
		out = append(out, strconv.Itoa(j)...)
		j++
	}

	return string(out)
}
func (db *postgresDB) Schema() string { return db.schema }
