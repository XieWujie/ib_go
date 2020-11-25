package imio

import (
	"../db"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
)

func getMembers(w http.ResponseWriter, r *http.Request) *AppError {
	var conversationId, _ = strconv.Atoi(r.URL.Query().Get("conversationId"))
	conversation := db.Conversation{ConversationId: conversationId}
	conversation.Get()
	list := make([]map[string]interface{}, len(conversation.Members))
	for i, v := range conversation.Members {
		user := db.User{UserId: v.UserId}
		_ = user.Get()
		m := make(map[string]interface{})
		m["username"] = user.Username
		m["avatar"] = user.Avatar
		m["userId"] = user.UserId
		m["nickname"] = v.NickName
		list[i] = m
	}
	sendOkWithData(w, list)
	return nil
}

func requestMessageList(w http.ResponseWriter, r *http.Request) *AppError {
	fmt.Println("请求消息列表")
	q := r.URL.Query()
	conversationId, _ := strconv.Atoi(q.Get("conversationId"))
	messageType, _ := strconv.Atoi(q.Get("messageType"))
	createAt := q.Get("before")
	before := int64(^uint64(0) >> 1)
	if len(createAt) != 0 {
		before, _ = strconv.ParseInt(createAt, 10, 64)
	}
	list := db.MessageFind(messageType, conversationId, before)
	sendOkWithData(w, list)
	return nil
}

func getRecentMessage(w http.ResponseWriter, r *http.Request) *AppError {
	var userId = r.URL.Query().Get("userId")
	id, _ := strconv.Atoi(userId)
	user := db.User{UserId: id}
	_ = user.Get()
	var list = make([]int, len(user.Friends)+len(user.Rooms))
	for i, v := range user.Rooms {
		list[i] = v
	}
	var roomPre = len(user.Rooms)
	for i, v := range user.Friends {
		list[i+roomPre] = v.ConversationId
	}
	var b, _ = json.Marshal(list)
	var str = string(b)
	str = strings.Replace(str, "[", "(", 1)
	str = strings.Replace(str, "]", ")", 1)
	result := db.FindRecentMessage(str)
	sendOkWithData(w, result)
	return nil
}
