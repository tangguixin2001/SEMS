package mysql

import (
	"database/sql"
	"log"
)

const (
	ROOT = "root"
	PWD  = "123456"
	ADDR = "127.0.0.1:3306"
	DB   = "sems"
)

func NewConn() *sql.DB {
	db, err := sql.Open("mysql", ROOT+":"+PWD+"@tcp("+ADDR+")/"+DB+"?charset=utf8")
	if err != nil {
		log.Println(err)
	}
	return db
}
