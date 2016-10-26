package models

import (
	"github.com/jinzhu/gorm"
)

type Book struct {
	gorm.Model
	Name   string
	Author string
	Price  uint
}
