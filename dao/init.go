package dao

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

var _db *sql.DB

func MySQLInit() {
	username := "root"
	password := "root"
	host := "127.0.0.1"
	port := 8889
	Dbname := "roomino"
	timeout := "10s"
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8&parseTime=True&loc=Local&timeout=%s",
		username, password, host, port, Dbname, timeout,
	)
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		panic("sqllink, error=" + err.Error())
	}

	db.SetMaxOpenConns(100)
	db.SetMaxIdleConns(20)
	db.SetConnMaxLifetime(time.Second * 30)

	_db = db

}

func NewDBClient(ctx context.Context) *sql.DB {
	return _db
}
