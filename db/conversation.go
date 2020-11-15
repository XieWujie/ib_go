package db

import (
	"database/sql"
	"fmt"
)

const COV_TYPE_SINGLE = 0
const COV_TYPE_ROOM = 1

type Conversation struct {
	ConversationId     int    `json:"conversation_id"`
	ConversationName   string `json:"conversation_name"`
	ConversationMember string `json:"conversation_member"`
	ConversationType   int    `json:"conversation_type"`
}

const conversationTable = "create table if not exists conversation(" +
	"conversationId integer primary key auto_increment," +
	"ConversationName varchar(16) not null," +
	"conversationMembers text," +
	"conversationType integer);"

func createConversationTable(sq *sql.DB) {
	if _, err := sq.Exec(conversationTable); err != nil {
		fmt.Println(err)
	}
}

func (c *Conversation) Save() error {
	tx := Db()
	stmt, _ := tx.Prepare("INSERT ignore conversation set conversationId=?,conversationName=?,conversationMembers=?,conversationType=?")
	result, _ := stmt.Exec(&c.ConversationId, &c.ConversationName, &c.ConversationMember, &c.ConversationType)
	idn, _ := result.LastInsertId()
	c.ConversationId = int(idn)
	err := tx.Commit()
	return err
}
func (c *Conversation) Get() error {
	tx := Db()
	stm, _ := tx.Prepare("Select * from conversation where ConversationId=?")
	row, _ := stm.Query(c.ConversationId)
	row.Next()
	_ = row.Scan(&c.ConversationId, &c.ConversationName, &c.ConversationMember, &c.ConversationType)
	err := tx.Commit()
	return err
}
