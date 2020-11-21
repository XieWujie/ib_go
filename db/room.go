package db

import "fmt"

type Room struct {
	ConversationId int    `json:"conversationId" xorm:"pk"`
	RoomName       string `json:"roomName"`
	RoomMasterId   int    `json:"roomMasterId"`
	RoomAvatar     string `json:"roomAvatar"`
}

type MemberInfo struct {
	UserId   int    `json:"userId"`
	NickName string `json:nickName`
}

func (room *Room) Save() error {
	_, err := engine.InsertOne(room)
	return err
}

func (room *Room) Get() error {
	_, err := engine.Get(room)
	return err
}

func FindRoom(roomIds string) []Room {
	var rooms []Room
	sql := fmt.Sprintf("conversation_id in %s", roomIds)
	err := engine.Where(sql).Find(&rooms)
	if err != nil {
		fmt.Println(err)
	}
	return rooms
}
