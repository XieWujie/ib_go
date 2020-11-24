package imio

import (
	"../db"
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"strconv"
)

func agreeAdd(m db.Verify) int {

	member := make([]db.MemberInfo, 2)
	members := append(member, db.MemberInfo{UserId: m.UserFrom}, db.MemberInfo{UserId: m.UserTo})
	cov := db.Conversation{ConversationType: db.COV_TYPE_SINGLE, Members: members}
	_ = cov.Save()
	user := &db.User{UserId: m.UserFrom}
	_ = user.Get()
	user.Friends = append(user.Friends, db.RelationShip{UserId: m.UserTo, ConversationId: cov.ConversationId})
	_ = user.Update()
	friend := &db.User{UserId: m.UserTo}
	_ = friend.Get()
	friend.Friends = append(friend.Friends, db.RelationShip{UserId: m.UserFrom, ConversationId: cov.ConversationId})
	_ = friend.Update()
	return cov.ConversationId
}

func sendChat(m db.Message) *AppError {
	cov := db.Conversation{
		ConversationId: m.ConversationId,
	}
	_ = cov.Get()
	var member = cov.Members
	for _, v := range member {
		sendMsgTo(m, v.UserId)
	}
	return nil
}

func sendMsgTo(message db.Message, to int) {
	ws, exit := wsConnAll[to]
	user := db.User{UserId: message.SendFrom}
	user.Get()
	var m = make(map[string]interface{})
	m["avatar"] = user.Avatar
	m["userId"] = user.UserId
	m["username"] = user.Username
	m["description"] = user.Description
	var target = make(map[string]interface{})
	target["message"] = message
	target["user"] = m
	if !exit {
		return
	}
	rec, _ := json.Marshal(target)
	fmt.Println("sendTo:", user.Username, string(rec))
	err := ws.wsWrite(websocket.TextMessage, rec)
	if err != nil {
		log.Println(err)
	}
}

func requestMessageList(w http.ResponseWriter, r *http.Request) *AppError {
	fmt.Println("请求消息列表")
	q := r.URL.Query()
	destination, _ := strconv.Atoi(q.Get("destination"))
	messageType, _ := strconv.Atoi(q.Get("messageType"))
	createAt := q.Get("before")
	before := int64(^uint64(0) >> 1)
	if len(createAt) != 0 {
		before, _ = strconv.ParseInt(createAt, 10, 64)
	}
	list := db.MessageFind(messageType, destination, before)
	sendOkWithData(w, list)
	return nil
}

func HandleMsg(msg []byte) *AppError {
	var m db.Message
	m.SendFrom = -1
	m.ConversationId = -1
	_ = json.Unmarshal(msg, &m)
	if m.SendFrom == -1 || m.MessageType == -1 {
		return &AppError{statusCode: 400, message: "from 和 messageType 不能为空"}
	}
	var err *AppError
	if m.MessageId == 0 {
		_ = m.Save()
	} else {
		_ = m.Update()
		return nil
	}
	switch m.MessageType {
	case db.TEXT:
		err = sendChat(m)
		break
	case db.WRITE:
		err = sendChat(m)
		break
	}
	return err
}
