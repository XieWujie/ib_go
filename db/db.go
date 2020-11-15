package db

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"log"
)

var db *sql.DB

func init() {
	db := createDatabase()
	createUserTable(db)
	createMessageTable(db)
	createConversationTable(db)
}

type Update interface {
	Update(key string, value string) error
}

type UpdateArray interface {
	UpdateArray(key string, value string) error
}

func Db() *sql.Tx {
	if db == nil {
		db = createDatabase()
	}
	tx, _ := db.Begin()
	return tx
}

func UpdateString(table string, key string, value string) {
	sql := "Update " + table + " Set " + key + " = " + value
	tx := Db()
	res, err := tx.Exec(sql)
	_ = tx.Commit()
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(res.RowsAffected())
	}
}

type Save interface {
	Save() error
}

type Get interface {
	Get() error
}

type DbError struct {
	message string
}

func (e DbError) Error() string {
	return e.message
}

func createDatabase() *sql.DB {
	path := "root:123456@tcp(127.0.0.1:3306)/imblog?charset=utf8"
	db, err := sql.Open("mysql", path)
	if err != nil {
		log.Fatal(err)
	}
	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}
	return db
}
