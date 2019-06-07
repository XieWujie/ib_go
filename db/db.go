package db

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"log"
)

var db *sql.DB

func init() {
	db := createDatabase()
	createUserTable(db)
}

type Update interface {
	Update(key string,value string)error
}

type UpdateArray interface {
	UpdateArray(key string,value string)error
}


func Db()*sql.DB  {
	if db == nil {
		db = createDatabase()
	}
	return db
}

type Save interface {
	Save()error
}

type Get interface {
	Get()error
}

type DbError struct {
	message string
}

func (e DbError)Error()string  {
	return e.message
}


func createDatabase() * sql.DB{
	path := "root:mysqlyyyy@tcp(127.0.0.1:3306)/imblog?charset=utf8"
	db,err := sql.Open("mysql",path)
	if err != nil {
		log.Fatal(err)
	}
	err =db.Ping()
	if err != nil {
		log.Fatal(err)
	}
	return db
}