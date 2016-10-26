package controllers

import (
	"github.com/revel/revel"
	"myapp/app/models"
)

type Book struct {
	App
}

func (c Book) Get(id int64) revel.Result {
	var book models.Book
	Db.First(&book, id)
	c.RenderArgs["book"] = book

	return c.Render()
}

func (c Book) List() revel.Result {
	var books []models.Book
	Db.Find(&books)
	c.RenderArgs["books"] = books

	return c.Render()
}
