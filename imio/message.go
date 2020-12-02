package imio

import (
	"../db"
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"io"
	"log"
	"net/http"
	"os"
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
	_ = user.Get()
	var m = make(map[string]interface{})
	m["avatar"] = user.Avatar
	m["userId"] = user.UserId
	m["username"] = user.Username
	m["description"] = user.Description
	var target = make(map[string]interface{})
	target["message"] = message
	target["user"] = m
	if message.FromType == db.MessageFromRoom {
		room := db.Room{ConversationId: message.ConversationId}
		room.Get()
		target["room"] = room
	}
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

func handleFileMsg(w http.ResponseWriter, r *http.Request) *AppError {
	_ = r.ParseMultipartForm(32 << 20)
	var message = r.FormValue("message")
	var m db.Message
	_ = json.Unmarshal([]byte(message), &m)
	file, handler, err := r.FormFile("file")
	if err != nil {
		log.Println(err)
		return &AppError{statusCode: 400, error: err}
	}
	defer file.Close()
	var filePath = createFilePath(handler.Filename)
	f, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE, 0666)
	defer f.Close()
	_, err = io.Copy(f, file)
	if err != nil {
		log.Println(err)
		return &AppError{statusCode: 400, error: err}
	}
	m.Content = createUrl(handler.Filename)
	var e = messageDispatch(m)
	if e != nil {
		return e
	}
	sendOkWithData(w, m)
	return nil
}

func HandleMsg(msg []byte) *AppError {
	var m db.Message
	_ = json.Unmarshal(msg, &m)
	return messageDispatch(m)
}

func messageDispatch(m db.Message) *AppError {
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
	case db.IMAGE:
		err = sendChat(m)
		break
	}
	return err
}
