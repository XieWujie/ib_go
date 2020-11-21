package db

import (
	"fmt"
)

type User struct {
	UserId      int            `json:"userId" xorm:"pk autoincr"`
	Username    string         `json:"username"`
	Password    string         `json:"password"`
	Avatar      string         `json:"avatar"`
	Description string         `json:"description"`
	Friends     []RelationShip `json:"friends" xorm:"json"`
	Rooms       []int          `json:"conversations" xorm:"json"`
}

type RelationShip struct {
	UserId         int `json:"user_id"`
	ConversationId int `json:"conversation_id"`
}

func (user *User) Get() error {
	_, err := engine.Get(user)
	return err
}

func FindUserByIds(ids string) []map[string]interface{} {
	var result []User
	_ = engine.Where("userId in ?", ids).Find(&result)

	list := make([]map[string]interface{}, len(result))
	for i, user := range result {
		m := make(map[string]interface{})
		m["username"] = user.Username
		m["avatar"] = user.Avatar
		m["userId"] = user.UserId
		m["description"] = user.Description
		list[i] = m
	}
	return list
}

func (user *User) Save() error {
	_, err := engine.InsertOne(user)
	return err
}

func (user *User) Update() error {
	_, err := engine.ID(user.UserId).Update(user)
	return err
}

func UserSearch(key string) []map[string]interface{} {
	var result []User
	var rows []map[string]interface{}
	err := engine.Where("username like ?", "%"+key+"%").Find(&result)
	if err != nil {
		fmt.Println(err)
		return rows
	}
	for _, user := range result {
		m := make(map[string]interface{})
		m["username"] = user.Username
		m["avatar"] = user.Avatar
		m["description"] = user.Description
		m["userId"] = user.UserId
		rows = append(rows, m)
	}
	return rows
}
