package db

import (
	"bytes"
	"database/sql"
	"log"
)

type User struct {
	UserId string `json:"userid"`
	Username string `json:"username"`
	Password string `json:"password"`
	Create string `json:"create"`
	Friends string `json:"friends"`
	Avatar string `json:"avatar"`
	Description string `json:"description"`
}

func (user * User)Get()error{
	if user.Username == "" && user.UserId == ""{
		return DbError{message:"username 和 userid不能都为空"}
	}
	var username string
	var userid  string
	var password string
	var friends string
	var avatar interface{}
	var description string
	var row *sql.Rows
	if len(user.Username)> 0 {
		stmt,err := Db().Prepare("SELECT * from user where username=?")
		if err != nil{
			return err
		}
		row,err = stmt.Query(user.Username)
		if err != nil{
			return err
		}
	}else {
		stmt,err := Db().Prepare("SELECT * from user where userid=?")
		if err != nil{
			return err
		}
		row,err = stmt.Query(user.UserId)
		if err != nil{
			return err
		}
	}
	row.Next()
	err := row.Scan(&userid,&username,&password,&avatar,&friends,&description)
	user.UserId = userid
	user.Username = username
	user.Password = password
	user.Friends = friends
	if avatar != nil{
		user.Avatar = avatar.(string)
	}
	user.Description = description
	if err != nil{
		return nil
	}
	return nil
}



func (user *User)Update(key string,value string)error  {
	var buffer bytes.Buffer
	buffer.WriteString("Update user set ")
	buffer.WriteString(key)
	buffer.WriteString("=")
	buffer.WriteString(value)
	buffer.WriteString("where userid='")
	buffer.WriteString(user.UserId)
	buffer.WriteString("'")
	log.Println(buffer.String())
	_, err := Db().Exec(buffer.String())
	log.Println(err)
	return err
}

func (user *User)UpdateArray(key string,value string)error  {
	var buffer bytes.Buffer
	buffer.WriteString("Update user set")
	buffer.WriteString(key)
	buffer.WriteString("=")
	buffer.WriteString(value)
	buffer.WriteString("where id=")
	buffer.WriteString(user.UserId)
	_, err := Db().Exec(buffer.String())
	return err
}



func (user User)Save()  error{
	stmt,err := Db().Prepare("INSERT ignore user set userid=?,username=?,password=?")
	if err != nil {
		return err
	}
	_,err = stmt.Exec(user.UserId,user.Username,user.Password)
	return err
}



const userTable  = "create table if not exists user (userid varchar(64) primary key ,username varchar(12) not null ,password varchar(16),avatar text,friends text,description text); "

func createUserTable(db *sql.DB)  {
	_,err := db.Exec(userTable)
	if err != nil {
		log.Fatal(err)
	}
}

