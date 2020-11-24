package db

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"xorm.io/xorm"
)

var engine *xorm.Engine

func init() {
	var err error
	engine, err = xorm.NewEngine("mysql", "root:123456@tcp(127.0.0.1:3306)/vlog?charset=utf8")
	if err != nil {
		fmt.Println(err)
	}
	err = nil
	err = engine.Sync(new(User), new(Verify), new(Room), new(Message), new(Conversation))
	if err != nil {
		fmt.Println(err)
	}
}

type Update interface {
	Update(key string, value string) error
}

type UpdateArray interface {
	UpdateArray(key string, value string) error
}

type Save interface {
	Save() error
}

type Get interface {
	Get() error
}
