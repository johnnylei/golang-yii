package example

import (
	"fmt"
	"github.com/johnnylei/golang-yii/example/models"
	"github.com/johnnylei/golang-yii/orm/mysql"
)

//orm 支持的query主要包括，string, 数组，Map,和回调函数
// string
func testQuery00()  {
	theme := models.NewDefaultTheme()

	// id=1
	data:= theme.Find().
		Field("*").
		Where("id=1").
		One()
	fmt.Println(data)
}

// 数组
func testQuery01()  {
	theme := models.NewDefaultTheme()

	// id=1 and tile like 'zz'
	data:= theme.Find().
		Field("*").
		Where([]interface{}{
			"id=1",
			"tile like 'zz'",
		}).
		One()
	fmt.Println(data)

	// id=1 or tile like 'zz'
	data = theme.Find().
		Field("*").
		Where([]interface{}{
			"or",
			"id=1",
			"tile like 'zz'",
		}).
		One()
	fmt.Println(data)

	// 嵌套式,支持无限次嵌套
	// []interface{}{
	//   	"=" #标识操作符号
	// 		"field" #字段
	// 		value #字段比较的值
	// }
	// id=1 or name like "zz"
	data = theme.Find().
		Field("*").
		Where([]interface{}{
			"or",
			[]interface{}{
				"=",
				"id",
				1,
			},
			[]interface{}{
				"like",
				"name",
				"zz",
			},
		}).
		One()
	fmt.Println(data)

	// id=1 or name like "zz" (id=1 or name like "zz")
	data = theme.Find().
		Field("*").
		Where([]interface{}{
			"or",
			[]interface{}{
				"=",
				"id",
				1,
			},
			[]interface{}{
				"like",
				"name",
				"zz",
			},
			[]interface{}{
				"or",
				[]interface{}{
					"=",
					"id",
					1,
				},
				[]interface{}{
					"like",
					"name",
					"zz",
				},
			},
		}).
		One()
	fmt.Println(data)
}

// map
func testQuery02()  {
	theme := models.NewDefaultTheme()

	// id=1 and name='xx'
	data:= theme.Find().
		Field("*").
		Where(map[string]interface{}{
			"id":1,
			"name":"xx",
		}).
		One()
	fmt.Println(data)

	// id=1 or name='xx'
	data = theme.Find().
		Field("*").
		Where([]interface{}{
			"or",
			map[string]interface{}{
				"id":1,
			},
			map[string]interface{}{
				"name":"zz",
			},
		}).
		One()
	fmt.Println(data)
}

// function
func testQuery03() {
	theme := models.NewDefaultTheme()

	// id=1 and name like 'xx'
	data := theme.Find().
		Field("*").
		Where(func(query *mysql.ActiveQuery) string {
			return "id=1 and name like 'zzz'"
		}).
		One()
	fmt.Println(data)

	// 函数结合数组
	data = theme.Find().
		Field("*").
		Where([]interface{}{
			"like",
			"name",
			func(query *mysql.ActiveQuery) string {
				return "zz"
			},
		}).
		One()
	fmt.Println(data)

	// 函数结合map, name='zz'
	data = theme.Find().
		Field("*").
		Where(map[string]interface{}{
			"name":func(query *mysql.ActiveQuery) string {
				return "zz"
			},
		}).
		One()
	fmt.Println(data)
}

// 复杂的query
func TestQuery()  {
	theme := models.NewDefaultTheme()

	data := theme.Find().
		Field("*").
		Alias("t1").
		Where([]interface{}{
			"or",
			func(query *mysql.ActiveQuery) string {
				return "name='johnny'"
			},
			func(query *mysql.ActiveQuery) string {
				return "id=1"
			},
			map[string]interface{}{
				"id":1,
				"name":"johnny",
				"sex": func(query *mysql.ActiveQuery) string{
					return "xxxx"
				},
			},
			[]interface{}{
				"=",
				"name",
				"johnny",
			},
			[]interface{}{
				"like",
				"sex",
				func(query *mysql.ActiveQuery) string{
					return "xxxx"
				},
			},
			[]interface{}{
				"or",
				func(query *mysql.ActiveQuery) string {
					return "name='johnny'"
				},
				func(query *mysql.ActiveQuery) string {
					return "id=1"
				},
				map[string]interface{}{
					"id":1,
					"name":"johnny",
					"sex": func(query *mysql.ActiveQuery) string{
						return "xxxx"
					},
				},
				[]interface{}{
					"=",
					"name",
					"johnny",
				},
				[]interface{}{
					"like",
					"sex",
					func(query *mysql.ActiveQuery) string{
						return "xxxx"
					},
				},
			},
		}).
		LeftJoin("user t2", "t1.`user_id`=t2.`id`").
		InnerJoin("user t2", "t1.`user_id`=t2.`id`").
		RightJoin("user t2", "t1.`user_id`=t2.`id`").
		OrderBy(map[string]string{
			"id":mysql.DESC,
			"name":mysql.ASC,
		}).
		Limit(0, 100).
		GroupBy([]string{
			"id",
			"name",
		}).
		One()
	fmt.Println(data)
}

