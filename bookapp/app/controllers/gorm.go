package controllers

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	"log"
	"myapp/app/models"
)

var (
	Db *gorm.DB
)

func InitDB() {
	var err error
	log.Println("Connecting to Database")
	Db, err = gorm.Open("mysql", "root:@/bookapp?parseTime=true&loc=Asia%2FTokyo")
	if err != nil {
		panic(err)
	}

	if Db.HasTable("books") == false {
		Db.AutoMigrate(&models.Book{})
		books := []models.Book{
			{Name: "Go Web Programming", Author: "Sau Sheong Chang", Price: 4782},
			{Name: "パターン認識と機械学習 上", Author: "C.M. ビショップ", Price: 7020},
			{Name: "行列プログラマー", Author: "Philip N. Klein", Price: 5832},
		}
		for _, book := range books {
			Db.Create(&book)
		}
	}
}
