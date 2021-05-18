package article

import (
	"goblog/app/models"
	"goblog/pkg/route"
	"strconv"
)

type Article struct {
	models.BaseModel

	Title string `gorm:"column:title;type:varchar(255)"`
	Body  string `gorm:"column:body;type:varchar(255)"`
}

func (a Article) Link() string {
	return route.Name2URL("articles.show", "id", strconv.FormatInt(int64(a.ID), 10))

}
