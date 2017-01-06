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
	Image   Image
	ImageID uint
}

type Image struct {
	gorm.Model
	Name string
	Url  string
}

func main() {
	db := prepare()
	defer db.Close()

	var user User
	db.First(&user, 1).Related(&user.Image)
	fmt.Printf("名前    : %s\n", user.Name)
	fmt.Printf("アイコン:\n")
	fmt.Printf("  名称: %s\n", user.Image.Name)
	fmt.Printf("  URL : %s\n", user.Image.Url)
}

func prepare() *gorm.DB {
	db, err := gorm.Open("mysql", dsn)
	if err != nil {
		panic("Failed to connect database")
	}
	db.DropTableIfExists(&User{}, &Image{})

	db.AutoMigrate(&Image{})
	db.AutoMigrate(&User{})

	var members = []map[string]string{
		{"Name": "ミク", "ImageFile": "miku.jpg"},
		{"Name": "マホ", "ImageFile": "maho.jpg"},
		{"Name": "コヒメ", "ImageFile": "kohime.jpg"},
	}

	for _, member := range members {
		image := Image{
			Name: member["Name"] + "アイコン",
			Url:  "https://image.example.com/" + member["ImageFile"],
		}
		db.Create(&image)

		user := User{
			Name:    member["Name"],
			ImageID: image.ID,
		}
		db.Create(&user)
	}

	return db
}
