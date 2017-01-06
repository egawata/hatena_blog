package main

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
)

const (
	dsn = "root@/association?parseTime=true&loc=Asia%2fTokyo"
)

type User struct {
	gorm.Model
	Name    string
	Profile Profile
}

type Profile struct {
	gorm.Model
	BirthDay string
	Hobby    string
	UserID   uint
}

func main() {
	db := prepare()
	defer db.Close()

	var user User
	db.First(&user, 2).Related(&user.Profile)
	fmt.Printf("名前    : %s\n", user.Name)
	fmt.Printf("誕生日  : %s\n", user.Profile.BirthDay)
	fmt.Printf("趣味    : %s\n", user.Profile.Hobby)
}

func prepare() *gorm.DB {
	db, err := gorm.Open("mysql", dsn)
	if err != nil {
		panic("Failed to connect database")
	}
	db.DropTableIfExists(&User{}, &Profile{})

	db.AutoMigrate(&Profile{})
	db.AutoMigrate(&User{})

	var members = []map[string]string{
		{"Name": "ミク", "BirthDay": "10/19", "Hobby": "にゃんにゃんダンス"},
		{"Name": "マホ", "BirthDay": "01/08", "Hobby": "漫画"},
		{"Name": "コヒメ", "BirthDay": "11/24", "Hobby": "茶道"},
	}

	for _, member := range members {
		user := User{
			Name: member["Name"],
		}
		db.Create(&user)

		profile := Profile{
			BirthDay: member["BirthDay"],
			Hobby:    member["Hobby"],
			UserID:   user.ID,
		}
		db.Create(&profile)
	}

	return db
}
