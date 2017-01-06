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
	Name      string
	FavMusics []Music `gorm:"many2many:user_fav_musics;"`
}

type Music struct {
	gorm.Model
	Title string
}

func main() {
	db := prepare()
	defer db.Close()

	for _, id := range []uint{1, 2, 3} {
		var user User
		db.First(&user, id).Related(&user.FavMusics, "FavMusics")
		fmt.Println("---------------------------")
		fmt.Printf("%s さんのお気に入り曲\n", user.Name)
		for _, music := range user.FavMusics {
			fmt.Printf("曲名: %s\n", music.Title)
		}
	}
	fmt.Println("---------------------------")
}

func prepare() *gorm.DB {
	db, err := gorm.Open("mysql", dsn)
	if err != nil {
		panic("Failed to connect database")
	}
	db.LogMode(true)
	db.DropTableIfExists(&User{}, &Music{}, "user_fav_musics")

	db.AutoMigrate(&User{})
	db.AutoMigrate(&Music{})

	prepareUsers(db)
	prepareMusics(db)
	prepareFavMusics(db)

	return db
}

func prepareUsers(db *gorm.DB) {
	var usernames = []string{"Mario", "ヤス", "たかし"}

	for _, username := range usernames {
		user := User{
			Name: username,
		}
		db.Create(&user)
	}
}

func prepareMusics(db *gorm.DB) {
	var titles = []string{
		"放課後ロマンス",
		"私Loveなオトメ",
		"ニーハイエゴイスト",
		"禁断無敵のダーリン",
		"S・M・L",
	}

	for _, title := range titles {
		music := Music{
			Title: title,
		}
		db.Create(&music)
	}
}

func prepareFavMusics(db *gorm.DB) {
	//  お気に入り楽曲
	//  {user.ID, music.ID} の組
	var favorites = [][]uint{
		{1, 1}, {1, 2}, {1, 4},
		{2, 2}, {2, 3}, {2, 5},
		{3, 2}, {3, 5},
	}
	for _, favorite := range favorites {
		var user User
		user.ID = favorite[0]

		var music Music
		music.ID = favorite[1]

		db.Model(&user).Association("FavMusics").Append(&music)
	}
}
