package db

import (
	"database/sql"
	"fmt"
)

type messageType int32

const SendVerifyAdd messageType = 3
const VerifyAgree messageType = 4
const TEXT messageType = 11
const WRITE = 12

type Message struct {
	MessageId   int         `json:"messageId"`
	MessageType messageType `json:"messageType"`
	SendFrom    int         `json:"sendFrom"`
	Destination int         `json:"destination"`
	Content     string      `json:"content"`
	CreateAt    int64       `json:"createAt"`
	Readed      bool        `json:"readed"`
}

const messageTable = "create table  if not exists message(" +
	"messageId integer primary key auto_increment," +
	"messageType integer," +
	"sendFrom integer," +
	"destination integer," +
	"content text," +
	"createAt long," +
	"readed bool);"

func createMessageTable(sq *sql.DB) {
	if _, err := sq.Exec(messageTable); err != nil {
		fmt.Println(err)
	}
}

func (m *Message) Save() error {
	tx := Db()
	stmt, err := tx.Prepare("INSERT ignore message set messageId=?,messageType=?,sendFrom=?,destination=?,content=?,createAt=?,readed=?")
	if err != nil {
		fmt.Println(err)
		return err
	}
	result, err := stmt.Exec(&m.MessageId, &m.MessageType, &m.SendFrom, &m.Destination, &m.Content, &m.CreateAt, &m.Readed)
	ids, _ := result.LastInsertId()
	m.MessageId = int(ids)
	_ = tx.Commit()
	if err != nil {
		fmt.Println(err)
	}
	return err
}
func (m *Message) Get(messageId string) error {
	tx := Db()
	stm, _ := tx.Prepare("Select * from message where messageId=?")
	row, err := stm.Query(messageId)
	if err != nil {
		fmt.Println(err)
		return err
	}
	row.Next()
	err = row.Scan(&m.MessageId, &m.MessageType, &m.SendFrom, &m.Destination, &m.Content, &m.CreateAt, &m.Readed)

	_ = tx.Commit()
	return err
}

func MessageFind(ty int, destination int, before int64) []map[string]interface{} {
	tx := Db()
	var sq string
	if ty >= 10 {
		sq = "Select * from message where destination=? and messageType>? and createAt<?"
	} else {
		sq = "Select * from message where destination=? and messageType=? and createAt<?"
	}
	stm, _ := tx.Prepare(sq)
	row, err := stm.Query(destination, ty, before)
	var list []map[string]interface{}
	if err != nil {
		fmt.Println(err)
		return list
	}
	for row.Next() {
		m := Message{}
		l := make(map[string]interface{})
		err = row.Scan(&m.MessageId, &m.MessageType, &m.SendFrom,
			&m.Destination, &m.Content, &m.CreateAt, &m.Readed)
		user := User{UserId: m.SendFrom}
		user.Get()
		l["user"] = user
		l["message"] = m
		list = append(list, l)
	}
	_ = tx.Commit()
	return list
}
