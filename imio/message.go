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
	"strconv"
)



func sendChat(m db.Message) *AppError {
	if m.MessageId == 0 {
		_ = m.Save()
	} else {
		_ = m.Update()
	}
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

func sendAgreeFriendMessage(from int, conversationId int, to int) {
	var message = db.Message{
		MessageType:    db.AgreeFriend,
		SendFrom:       from,
		ConversationId: conversationId,
		FromType:       db.MessageFromFriend,
		CiteMessageId:  -1,
	}
	_ = message.Save()
	sendMsgTo(message, to)
}

func sendVerifyMessage(from int, to int) {
	var message = db.Message{
		MessageType:    db.VerifyMessage,
		SendFrom:       from,
		ConversationId: -1,
		FromType:       db.MessageFromFriend,
		CiteMessageId:  -1,
	}
	_ = message.Save()
	sendMsgTo(message, to)
}

func sendMsgTo(message db.Message, to int) {
	ws, exit := wsConnAll[to]
	if !exit {
		return
	}
	var toUser = db.User{UserId:to}
	_ = toUser.Get()
	if message.FromType == db.MessageFromRoom {
		for _,v := range toUser.Rooms{
			if message.ConversationId == v.ConversationId{
				message.Notify = v.Notify
				break
			}
		}
	}else if message.FromType == db.MessageFromFriend {
		for _,v := range toUser.Friends{
			if message.ConversationId == v.ConversationId{
				message.Notify = v.Notify
				break
			}
		}
	}else {
		message.Notify = true
	}
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
	case db.WITHDRAW:
		err = withdraw(m)
		break
	case db.RECORD:
		err = sendChat(m)
		break
	}
	return err
}

func withdraw(m db.Message) *AppError {
	var id, _ = strconv.ParseInt(m.Content, 10, 32)
	var message = db.Message{MessageId: int(id)}
	_ = message.Get()
	message.MessageType = db.WITHDRAW
	message.Content = "消息撤回"
	return sendChat(message)
}
