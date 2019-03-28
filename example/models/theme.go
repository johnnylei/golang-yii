package models

import (
	"github.com/johnnylei/golang-yii/orm/mysql"
)

type Theme struct {
	*mysql.ActiveRecord

	Id int `field:"id" primary`
	Name string `field:"name"`
	Sex string `field:"sex"`
}

func NewDefaultTheme() *Theme  {
	theme := &Theme{
		ActiveRecord:&mysql.ActiveRecord{},
	}
	theme.Table = "input_skinthemes"
	theme.Instance = theme
	theme.Init()
	return theme
}
