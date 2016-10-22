package main

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	"log"
	"time"
)

const (
	dsn = "user1:password1@tcp(127.0.0.1:3306)/mydb2?parseTime=true&loc=Asia%2FTokyo"
)

type Member struct {
	Id        int64 `gorm:"primary_key"`
	Name      string
	Birthday  string
	BloodType string
	Hobbies   []Hobby
	CreatedAt time.Time
	UpdatedAt time.Time
}

type Hobby struct {
	Id        int64 `gorm:"primary_key"`
	MemberId  int64
	Name      string
	CreatedAt time.Time
	UpdatedAt time.Time
}

//  デフォルトではテーブル名は型名を複数形にしたもの(member -> members)
//  異なるテーブル名を使用したい場合は TableName() でテーブル名を返すようにする
func (m *Member) TableName() string {
	return "member"
}

func (h *Hobby) TableName() string {
	return "hobby"
}

func main() {
	db, err := gorm.Open("mysql", dsn)
	checkErr(err, "Failed to connect database")
	defer db.Close()

	//  レコードの追加
	//  memberだけでなく関係テーブルのレコードも同時に追加している
	members := []Member{
		{Name: "ミク", Birthday: "10/19", BloodType: "AB", Hobbies: []Hobby{{Name: "ブログ"}, {Name: "ショッピング"}}},
		{Name: "マホ", Birthday: "1/8", BloodType: "AB", Hobbies: []Hobby{{Name: "漫画"}, {Name: "ゲーム"}}},
		{Name: "コヒメ", Birthday: "11/24", BloodType: "O", Hobbies: []Hobby{{Name: "ゲーム"}, {Name: "茶道"}}},
	}
	for _, member := range members {
		db.Create(&member)
	}

	fmt.Println("\nすべてのメンバーを取得:\n")
	var allMembers []Member
	db.Find(&allMembers)
	fmt.Println(allMembers)

	fmt.Println("\n\nメンバー1人のみ取得(Hobbiesもあわせて取得):\n")
	var miku Member
	db.Where("name = ?", "ミク").First(&miku)
	db.Model(&miku).Related(&miku.Hobbies)
	fmt.Printf("%#v\n", miku)
}

func checkErr(err error, msg string) {
	if err != nil {
		log.Fatalf("%s : %s\n", err, msg)
	}
}
