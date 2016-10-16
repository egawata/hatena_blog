package main

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"log"
)

type Member struct {
	Id        int
	Name      string
	Birthday  string
	Bloodtype string
	Hobby     string
}

var Db *sql.DB

func registerMembers() (err error) {
	members := []Member{
		{Name: "ミク", Birthday: "10/19", Bloodtype: "AB", Hobby: "ショッピング"},
		{Name: "マホ", Birthday: "1/8", Bloodtype: "AB", Hobby: "漫画"},
		{Name: "コヒメ", Birthday: "11/24", Bloodtype: "O", Hobby: "ゲーム"},
	}

	stmt, err := Db.Prepare(`
		INSERT INTO member (name, birthday, blood_type, hobby) 
		VALUES (?, ?, ?, ?)
	`)
	if err != nil {
		return
	}
	defer stmt.Close()

	for _, m := range members {
		var ret sql.Result
		ret, err = stmt.Exec(m.Name, m.Birthday, m.Bloodtype, m.Hobby)
		if err != nil {
			return
		}

		var id int64
		id, err = ret.LastInsertId()
		if err != nil {
			return
		}
		m.Id = int(id)
		fmt.Printf("Inserted. ID = %d, Name = %s\n", m.Id, m.Name)
	}

	return
}

func getAllMembers() (members []Member, err error) {
	rows, err := Db.Query("SELECT id, name, birthday, blood_type, hobby FROM member")
	if err != nil {
		return
	}

	for rows.Next() {
		m := Member{}
		err = rows.Scan(&m.Id, &m.Name, &m.Birthday, &m.Bloodtype, &m.Hobby)
		if err != nil {
			return
		}
		members = append(members, m)
	}

	return
}

func getMember(id int) (member Member, err error) {
	member = Member{}

	err = Db.QueryRow(`
		SELECT id, name, birthday, blood_type, hobby 
		  FROM member
		 WHERE id = ?
	`, id).Scan(&member.Id, &member.Name, &member.Birthday, &member.Bloodtype, &member.Hobby)
	return
}

func main() {
	var err error
	Db, err = sql.Open("mysql", "user1:password1@tcp(127.0.0.1:3306)/mydatabase?charset=utf8")
	if err != nil {
		log.Fatalln(err)
	}
	defer Db.Close()

	err = registerMembers()
	if err != nil {
		log.Fatalln(err)
	}

	members, err := getAllMembers()
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println("Get all members:")
	fmt.Println(members)

	member, err := getMember(members[1].Id)
	if err == sql.ErrNoRows {
		log.Fatalln("そのIDのメンバーは存在しません")
	} else if err != nil {
		log.Fatalln(err)
	}
	fmt.Println("Get one member")
	fmt.Println(member)
}
