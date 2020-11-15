package db

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"strconv"
)

type User struct {
	UserId      int    `json:"userId"`
	Username    string `json:"username"`
	Password    string `json:"password"`
	Avatar      string `json:"avatar"`
	Description string `json:"description"`
	Friends     string `json:"friends"`
	Rooms       string `json:"conversations"`
}

type RelationShip struct {
	UserId         int `json:"user_id"`
	ConversationId int `json:"conversation_id"`
}

func (user *User) Get() error {
	tx := Db()
	var row *sql.Rows
	if len(user.Username) > 0 {
		stmt, err := tx.Prepare("SELECT * from user where username=?")
		if err != nil {
			return err
		}
		row, err = stmt.Query(user.Username)
		if err != nil {
			return err
		}
	} else {
		stmt, err := tx.Prepare("SELECT * from user where userid=?")
		if err != nil {
			return err
		}
		row, err = stmt.Query(user.UserId)
		if err != nil {
			return err
		}
	}
	row.Next()
	err := row.Scan(&user.UserId, &user.Username, &user.Password, &user.Avatar, &user.Description, &user.Friends, &user.Rooms)
	if err != nil {
		return nil
	}
	_ = tx.Commit()
	return nil
}

func (user User) Save() error {
	tx := Db()
	stmt, err := tx.Prepare("INSERT ignore user set username=?,password=?,avatar=?,description=?,friends=?,rooms=?")
	if err != nil {
		println(string(err.Error()))
		return err
	}
	_, err = stmt.Exec(&user.Username, &user.Password, &user.Avatar, &user.Description, &user.Friends, &user.Rooms)
	id, _ := tx.Exec("SELECT LAST_INSERT_ID()")
	idn, _ := id.LastInsertId()
	user.UserId = int(idn)
	_ = tx.Commit()
	return err
}

func (user *User) AddF(userid int, covId int) {
	rel := RelationShip{
		UserId:         userid,
		ConversationId: covId,
	}
	var fs []RelationShip
	json.Unmarshal([]byte(user.Friends), &fs)
	fu := append(fs, rel)
	rec, _ := json.Marshal(fu)
	tx := Db()
	result, err := tx.Exec("UPDATE user SET friends=? WHERE userid=?", string(rec), user.UserId)
	_ = tx.Commit()
	v, e := result.RowsAffected()
	if e == nil {
		fmt.Println("更新了user记录" + strconv.FormatInt(v, 10))
	}
	if err != nil {
		fmt.Println(err)
	}
}

const userTable = "create table if not exists user (" +
	"userid integer primary key auto_increment," +
	"username varchar(12) not null ," +
	"password varchar(16)," +
	"avatar text," +
	"description text," +
	"friends text," +
	"rooms text); "

func createUserTable(db *sql.DB) {
	_, err := db.Exec(userTable)
	if err != nil {
		log.Fatal(err)
	}
}

const searchSql = "select * from user where username LIKE ?"

func UserSearch(key string) []map[string]interface{} {
	tx := Db()
	st, _ := tx.Prepare(searchSql)
	row, err := st.Query("%" + key + "%")
	rows := make([]map[string]interface{}, 0, 10)
	if err != nil {
		fmt.Println(err)
		return rows
	}
	i := 0
	for row.Next() && i < 10 {
		user := User{}
		_ = row.Scan(&user.UserId, &user.Username, &user.Password, &user.Avatar, &user.Description, &user.Friends, &user.Rooms)
		m := make(map[string]interface{})
		m["username"] = user.Username
		m["avatar"] = user.Avatar
		m["description"] = user.Description
		m["userId"] = user.UserId
		rows = append(rows, m)
		i++
	}
	_ = tx.Commit()
	return rows
}
