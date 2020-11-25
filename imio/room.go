package imio

import (
	"../db"
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
)

type roomCreate struct {
	Room    db.Room         `json:"room"`
	Members []db.MemberInfo `json:"members"`
}

func createRoom(w http.ResponseWriter, r *http.Request) *AppError {
	en := new(roomCreate)
	_ = json.NewDecoder(r.Body).Decode(&en)
	conversation := db.Conversation{ConversationType: db.COV_TYPE_ROOM, Members: en.Members}
	_ = conversation.Save()
	en.Room.ConversationId = conversation.ConversationId
	if len(en.Room.RoomName) == 0 {
		for i, v := range en.Members {

			user := db.User{UserId: v.UserId}
			_ = user.Get()
			if i == 1 {
				en.Room.RoomName = user.Username
			} else {
				en.Room.RoomName += "&" + user.Username
			}
			if i > 4 {
				break
			}
		}
	}
	for _, v := range en.Members {
		user := db.User{UserId: v.UserId}
		_ = user.Get()
		user.Rooms = append(user.Rooms, conversation.ConversationId)
		_ = user.Update()
	}
	_ = en.Room.Save()
	sendOkWithData(w, en.Room)
	return nil
}

func getRoom(w http.ResponseWriter, r *http.Request) *AppError {
	var userId, _ = strconv.Atoi(r.URL.Query().Get("userId"))
	user := db.User{UserId: userId}
	_ = user.Get()
	var rooms, _ = json.Marshal(user.Rooms)
	rstr := string(rooms)
	rstr = strings.Replace(rstr, "[", "(", 1)
	rstr = strings.Replace(rstr, "]", ")", -1)
	list := db.FindRoom(rstr)
	sendOkWithData(w, list)
	return nil
}

func roomUpdate(w http.ResponseWriter, r *http.Request) *AppError {
	room := new(db.Room)
	json.NewDecoder(r.Body).Decode(&room)
	_ = room.Update()
	sendOkWithData(w, room)
	return nil
}
