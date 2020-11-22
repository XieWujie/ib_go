package imio

import (
	"../db"
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
