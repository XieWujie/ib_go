package imio

import (
	"../db"
	"fmt"
	"net/http"
	"strconv"
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
