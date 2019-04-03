package model

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql"
)

type Users struct {
	Id       int
	Username string
	Roles_id int
}

func DBInit() *sql.DB {
	db, err := sql.Open("mysql", "root:123456@tcp(192.168.2.100:3306)/test?charset=utf8")
	if err != nil {
		log.Fatal(err)
	}
	return db
}

var db = DBInit()

func Insert() {
	stmt, _ := db.Prepare("insert into users set username=?, roles_id=?")
	res, _ := stmt.Exec("zhang", 1)
	id, _ := res.LastInsertId()
	fmt.Println(id)
}

func (_ Users) SelectAll() []Users {
	rows, _ := db.Query("select * from users")
	users := Users{}
	var userlist []Users

	for rows.Next() {
		err := rows.Scan(&users.Id, &users.Username, &users.Roles_id)
		if err != nil {
			log.Fatal(err)
		}
		userlist = append(userlist, users)
	}
	return userlist
}

func (u Users) SelectOne() Users {
	id := u.Id
	user := Users{}
	err := db.QueryRow("select * from users where id = ?", id).Scan(&user.Id, &user.Username, &user.Roles_id)
	if err != nil {
		log.Fatal(err)
	}
	return user
}
