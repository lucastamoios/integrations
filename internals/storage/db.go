package storage

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/nleof/goyesql"
)

type Database interface {
	Select(dest interface{}, query goyesql.Tag, args ...interface{}) error
	Exec(query goyesql.Tag, args ...interface{}) (int64, error)
}

type Postgres struct {
	DB *sqlx.DB
	queries goyesql.Queries
}

func NewPostgresDatabase(dbName, queryFile string) (*Postgres, error) {
	dbSource := fmt.Sprintf("dbname=%s sslmode=disable", dbName)
	db, err := sqlx.Open("postgres", dbSource)
	if err != nil {
		return nil, err
	}
	queries, err := goyesql.ParseFile(queryFile)
	if err != nil {
		return nil, err
	}
	return &Postgres{db, queries}, nil
}

func (db *Postgres) Select(dest interface{}, query goyesql.Tag, args ...interface{}) error{
	return db.DB.Select(dest, db.queries[query], args)
}

func (db *Postgres) Exec(query goyesql.Tag, args ...interface{}) (int64, error){
	r, err := db.DB.Exec(db.queries[query], args)
	if err != nil {
		return 0, err
	}
	lid, err := r.LastInsertId()
	return lid, err
}