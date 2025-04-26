package database

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
)

type MDB struct {
	db *sql.DB
}

func NewMDB(dataSourceName string) (*MDB, error) {
	db, err := sql.Open("mysql", dataSourceName)
	if err != nil {
		return nil, err
	}
	return &MDB{db: db}, nil
}
