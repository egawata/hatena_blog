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
	Name     string
	Articles []Article
}

type Article struct {
	gorm.Model
	Title      string
	Body       string
	PostUserID uint
}

func main() {
	db := prepare()
	defer db.Close()

	for _, id := range []uint{1, 2, 3} {
		var user User
		db.First(&user, id).Related(&user.Articles, "PostUserID")
		fmt.Printf("%s さんの投稿一覧\n", user.Name)
		for _, article := range user.Articles {
			fmt.Println("--------------------------------")
			fmt.Printf("タイトル: %s\n", article.Title)
			fmt.Printf("本文    : %s\n", article.Body)
		}
		fmt.Printf("--------------------------------\n\n")
	}
}

func prepare() *gorm.DB {
	db, err := gorm.Open("mysql", dsn)
	if err != nil {
		panic("Failed to connect database")
	}
	db.DropTableIfExists(&User{}, &Article{})

	db.AutoMigrate(&User{})
	db.AutoMigrate(&Article{})

	prepareUsers(db)
	prepareArticles(db)

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

func prepareArticles(db *gorm.DB) {
	var articles = []map[string]interface{}{
		{
			"Title":  "Type Switch",
			"Body":   "interface{}型で宣言された変数の値の具体的な型により処理を分岐させたい、というケースが考えられる。",
			"UserID": uint(1),
		},
		{
			"Title":  "GORMを使ってみた",
			"Body":   "goで使えるO/R Mapperは結構いろいろあるみたいで、とりあえず今日1日 gorpとGORMを使ってみた。",
			"UserID": uint(1),
		},
		{
			"Title":  "goでmysqlを使うメモ",
			"Body":   "前回の日記が2014年11月3日…だいぶ放置してたなor",
			"UserID": uint(2),
		},
		{
			"Title":  "nnetを使ってコード解析作ってみた",
			"Body":   "以前このブログにも書いた通り、今年の夏に Coursera の Machine Learning の講座を受けた。",
			"UserID": uint(2),
		},
		{
			"Title":  "rgl の plot3d で3Dグラフを描く",
			"Body":   "R上で、OpenGLを使って高度な3Dグラフを描くためのパッケージ。",
			"UserID": uint(3),
		},
		{
			"Title":  "Coursera の Machine Learning",
			"Body":   "Coursera というオンライン学習サイトに Machine Learning という講座がある。",
			"UserID": uint(3),
		},
	}

	for _, articleData := range articles {
		article := Article{
			Title:      articleData["Title"].(string),
			Body:       articleData["Body"].(string),
			PostUserID: articleData["UserID"].(uint),
		}
		db.Create(&article)
	}
}
