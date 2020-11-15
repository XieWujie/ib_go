package imio

import (
	"../db"
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"net/http"
	"strconv"
	"time"
)

func sendVerify(m db.Message) *AppError {
	sendMsgTo(m, m.Destination)
	return nil
}

func agreeAdd(m db.Message) *AppError {
	member := []int{m.SendFrom, m.Destination}
	memStr, _ := json.Marshal(member)
	cov := db.Conversation{ConversationType: db.COV_TYPE_SINGLE, ConversationMember: string(memStr)}
	_ = cov.Save()
	user := &db.User{UserId: m.SendFrom}
	_ = user.Get()
	user.AddF(m.Destination, cov.ConversationId)
	friend := &db.User{UserId: m.Destination}
	_ = friend.Get()
	friend.AddF(m.SendFrom, cov.ConversationId)
	_ = m.Save()
	sendMsgTo(m, m.Destination)
	return nil
}

func sendChat(m db.Message) *AppError {
	cov := db.Conversation{
		ConversationId: m.Destination,
	}
	_ = cov.Get()
	var member []int
	_ = json.Unmarshal([]byte(cov.ConversationMember), &member)
	for _, v := range member {
		sendMsgTo(m, v)
	}
	return nil
}

func sendMsgTo(message db.Message, to int) {
	ws, exit := wsConnAll[to]
	if !exit {
		return
	}
	user := db.User{UserId: message.SendFrom}
	user.Get()
	m := make(map[string]interface{})
	m["username"] = user.Username
	m["avatar"] = user.Avatar
	m["description"] = user.Description
	m["userId"] = user.UserId
	wrap := make(map[string]interface{})
	wrap["message"] = message
	wrap["user"] = user
	rec, _ := json.Marshal(wrap)
	_ = ws.wsWrite(websocket.TextMessage, rec)
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
	m.Destination = -1
	_ = json.Unmarshal(msg, &m)
	if m.SendFrom == -1 || m.MessageType == -1 {
		return &AppError{statusCode: 400, message: "from 和 messageType 不能为空"}
	}
	if m.CreateAt == 0 {
		m.CreateAt = time.Now().Unix()
	}
	var err *AppError
	_ = m.Save()
	switch m.MessageType {
	case db.SendVerifyAdd:
		err = sendVerify(m)
		break
	case db.VerifyAgree:
		err = agreeAdd(m)
		break
	case db.TEXT:
		err = sendChat(m)
		break
	case db.WRITE:
		err = sendChat(m)
		break
	}
	return err
}
