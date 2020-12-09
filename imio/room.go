package imio

import (
	"../db"
	"encoding/json"
	"net/http"
	"strconv"
)

type roomCreate struct {
	Room    db.Room         `json:"room"`
	Members []db.MemberInfo `json:"members"`
}

type roomMemberAdd struct {
	ConversationId int             `json:"conversationId"`
	Members        []db.MemberInfo `json:"members"`
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
		user.Rooms = append(user.Rooms, db.RoomsMeta{ConversationId: conversation.ConversationId, Notify: true})
		_ = user.Update()
	}
	_ = en.Room.Save()
	sendOkWithData(w, en.Room)
	return nil
}

func addRoomMember(w http.ResponseWriter, r *http.Request) *AppError {
	en := roomMemberAdd{}
	_ = json.NewDecoder(r.Body).Decode(&en)
	var conversation = db.Conversation{ConversationId: en.ConversationId}
	_ = conversation.Get()
	for _, v := range en.Members {
		conversation.Members = append(conversation.Members, v)
	}
	_ = conversation.Update()
	sendOkWithData(w, conversation)
	return nil
}

type roomQuit struct {
	UserId         int `json:"userId"`
	ConversationId int `json:"conversationId"`
}

func quitRoom(w http.ResponseWriter, r *http.Request) *AppError {
	en := roomQuit{}
	_ = json.NewDecoder(r.Body).Decode(&en)
	var conversation = db.Conversation{ConversationId: en.ConversationId}
	_ = conversation.Get()
	var members = conversation.Members
	var newMembers = make([]db.MemberInfo, 0)
	for _, v := range members {
		if v.UserId != en.UserId {
			newMembers = append(newMembers, v)
		}
	}
	conversation.Members = newMembers
	_ = conversation.Update()
	var user = db.User{UserId: en.UserId}
	user.Get()
	var newRooms = make([]db.RoomsMeta, 0)
	for _, v := range user.Rooms {
		if v.ConversationId != en.ConversationId {
			newRooms = append(newRooms, v)
		}
	}
	user.Rooms = newRooms
	_ = user.Update()
	sendOkWithData(w, conversation)
	return nil
}

func getRoom(w http.ResponseWriter, r *http.Request) *AppError {
	var userId, _ = strconv.Atoi(r.URL.Query().Get("userId"))
	user := db.User{UserId: userId}
	_ = user.Get()
	var rooms = user.Rooms
	var list = make([]map[string]interface{}, len(rooms))
	for i, v := range rooms {
		room := db.Room{ConversationId: v.ConversationId}
		room.Get()
		var m = make(map[string]interface{})
		m["conversationId"] = room.ConversationId
		m["roomAvatar"] = room.RoomAvatar
		m["roomMasterId"] = room.RoomMasterId
		m["roomName"] = room.RoomName
		m["notify"] = v.Notify
		m["background"] = v.Background
		list[i] = m
	}
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
