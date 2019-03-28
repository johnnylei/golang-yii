package example

import (
	"github.com/johnnylei/golang-yii/base"
	"github.com/johnnylei/golang-yii/example/models"
)

func TestInsert() {
	theme := models.NewDefaultTheme()
	data := map[string]interface{}{
		"name":"johnny",
		"sex":"male",
	}
	base.Load(theme, data)
	println(theme.Insert())
}

func TestUpdate()  {
	theme := models.NewDefaultTheme()

	data:= theme.Find().
		Field("*").
		Where([]interface{}{
			"=",
			"id",
			1,
		}).
		One()
	_data := make(map[string]interface{})
	for key, value := range data {
		_data[key] = interface{}(value)
	}
	base.Load(theme, _data)
	theme.Sex = "female"
	theme.Name = "xxxx"
	println(theme.Update())
}
