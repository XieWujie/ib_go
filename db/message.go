package db

import (
	"database/sql"
	"fmt"
)

type messageType int32

const VerifyMessage messageType = 3
const TEXT messageType = 11
const WRITE = 12

type Message struct {
	MessageId      int         `json:"messageId" xorm:"pk autoincr"`
	MessageType    messageType `json:"messageType"`
	SendFrom       int         `json:"sendFrom"`
	ConversationId int         `json:"conversationId"`
	Content        string      `json:"content"`
	CreateAt       int64       `json:"createAt" xorm:"updated"`
	Readed         bool        `json:"readed"`
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
	_, err := engine.InsertOne(m)
	return err
}
func (m *Message) Get() error {
	_, err := engine.ID(m.MessageId).Get(m)
	return err
}

func MessageFind(ty int, destination int, before int64) []map[string]interface{} {
	var result []Message
	_ = engine.Where("conversation_id=? and create_at<?", destination, before).Find(&result)
	var list = make([]map[string]interface{}, len(result))
	for i, v := range result {
		var user = User{UserId: v.SendFrom}
		user.Get()
		var m = make(map[string]interface{})
		m["avatar"] = user.Avatar
		m["userId"] = user.UserId
		m["username"] = user.Username
		m["description"] = user.Description
		m["description"] = user.Description
		var r = make(map[string]interface{})
		r["user"] = m
		r["message"] = v
		list[i] = r
	}
	return list
}

func (m *Message) Update() error {
	_, err := engine.ID(m).Update(m)
	return err
}

func FindRecentMessage(conversations string) []map[string]interface{} {
	var result []Message
	var sql = "message_id in (select message_id from message where create_at in (select MAX(create_at) from message group by conversation_id)) and conversation_id in %s"
	sql = fmt.Sprintf(sql, conversations)
	engine.Where(sql).Find(&result)
	var list = make([]map[string]interface{}, len(result))
	for i, v := range result {
		var user = User{UserId: v.SendFrom}
		user.Get()
		var m = make(map[string]interface{})
		m["avatar"] = user.Avatar
		m["userId"] = user.UserId
		m["username"] = user.Username
		m["description"] = user.Description
		m["description"] = user.Description
		var r = make(map[string]interface{})
		r["user"] = m
		r["message"] = v
		list[i] = r
	}
	return list
}
