package mysql

import (
	"fmt"
	"github.com/johnnylei/golang-yii/base"
	"github.com/johnnylei/golang-yii/orm"
	"github.com/johnnylei/golang-yii/utils"
	"math/rand"
	"strconv"
	"strings"
)

const LEFT_JOIN  = "LEFT JOIN"
const RIGHT_JOIN  = "RIGHT JOIN"
const INNER_JOIN  = "INNER JOIN"
const DESC  = "DESC"
const ASC  = "ASC"
const START_FETCH_ONE = "START_FETCH_ONE"
const END_FETCH_ONE = "END_FETCH_ONE"
const START_FETCH_ALL = "START_FETCH_ALL"
const END_FETCH_ALL = "END_FETCH_ALL"

type ActiveQuery struct {
	base.Object
	table string
	orderBy map[string]string
	where interface{}
	field interface{}
	start int
	rowsLength int
	groupBy []string
	having map[string]interface{}
	alias string
	sqlTemplate string
	join []map[string]string
	whereParams map[string]interface{}
}

func NewDefaultActiveQuery() *ActiveQuery  {
	return &ActiveQuery{
		field:"*",
		start:0,
		alias:"t1",
		whereParams:make(map[string]interface{}),
	}
}

func (activeQuery *ActiveQuery) Alias(alias string) orm.ActiveQueryInterface  {
	activeQuery.alias = alias
	return activeQuery
}

func (activeQuery *ActiveQuery) One() map[string]string  {
	base.Trigger(START_FETCH_ONE, activeQuery)
	activeQuery.buildSql()
	command := NewDefaultCommand()
	data := command.Query(activeQuery.sqlTemplate, activeQuery.whereParams, true)
	base.Trigger(END_FETCH_ONE, activeQuery)
	if len(data) == 0 {
		return map[string]string{}
	}

	return data[0]
}

func (activeQuery *ActiveQuery) All() []map[string]string {
	base.Trigger(START_FETCH_ALL, activeQuery)
	activeQuery.buildSql()
	data := NewDefaultCommand().Query(activeQuery.sqlTemplate, activeQuery.whereParams, false)
	base.Trigger(END_FETCH_ALL, activeQuery)
	return data
}

func (activeQuery *ActiveQuery) OrderBy(orderBy map[string]string) orm.ActiveQueryInterface  {
	activeQuery.orderBy = orderBy
	return activeQuery
}

func (activeQuery *ActiveQuery) Field(field interface{}) orm.ActiveQueryInterface  {
	activeQuery.field = field
	return activeQuery
}

func (activeQuery *ActiveQuery) Where(where interface{}) orm.ActiveQueryInterface  {
	activeQuery.where = where
	return activeQuery
}

func (activeQuery *ActiveQuery) Count() orm.ActiveQueryInterface  {
	return activeQuery
}

func (activeQuery *ActiveQuery) Limit(start int, rowsLength int) orm.ActiveQueryInterface  {
	activeQuery.start = start
	activeQuery.rowsLength = rowsLength
	return activeQuery
}

func (activeQuery *ActiveQuery) GroupBy(groupBy []string) orm.ActiveQueryInterface  {
	activeQuery.groupBy = groupBy
	return activeQuery
}

func (activeQuery *ActiveQuery) Having(having map[string]interface{}) orm.ActiveQueryInterface  {
	activeQuery.having = having
	return activeQuery
}

func (activeQuery *ActiveQuery) LeftJoin(table string, on string) orm.ActiveQueryInterface  {
	activeQuery.join = append(activeQuery.join, map[string]string{
		"join":LEFT_JOIN,
		"table":table,
		"on":on,
	})
	return activeQuery
}

func (activeQuery *ActiveQuery) RightJoin(table string, on string) orm.ActiveQueryInterface  {
	activeQuery.join = append(activeQuery.join, map[string]string{
		"join":RIGHT_JOIN,
		"table":table,
		"on":on,
	})
	return activeQuery
}

func (activeQuery *ActiveQuery) InnerJoin(table string, on string) orm.ActiveQueryInterface  {
	activeQuery.join = append(activeQuery.join, map[string]string{
		"join":INNER_JOIN,
		"table":table,
		"on":on,
	})
	return activeQuery
}

func (activeQuery *ActiveQuery) buildSql() string  {
	if activeQuery.table == "" {
		panic("table name should not be null")
	}

	// build field
	field := "*"
	if activeQuery.field != nil {
		field = activeQuery.buildField()
	}

	// build with alias
	activeQuery.sqlTemplate = fmt.Sprintf("SELECT %s FROM `%s`", field, activeQuery.table)
	if activeQuery.alias != "" {
		activeQuery.sqlTemplate = fmt.Sprintf("%s AS %s", activeQuery.sqlTemplate, activeQuery.alias)
	}

	// build join
	if activeQuery.join != nil {
		for _, item := range activeQuery.join {
			activeQuery.sqlTemplate = fmt.Sprintf(
				"%s %s %s ON %s",
				activeQuery.sqlTemplate,
				item["join"],
				item["table"],
				item["on"],
			)
		}
	}

	// build condition
	if activeQuery.where != nil {
		condition := activeQuery.buildCondition(activeQuery.where)
		if condition != "" {
			_condition := []rune(condition)
			condition = string(_condition[1:len(_condition) - 1])
			activeQuery.sqlTemplate = fmt.Sprintf("%s WHERE %s", activeQuery.sqlTemplate, condition)
		}
	}

	// build group by
	if activeQuery.groupBy != nil {
		activeQuery.sqlTemplate = fmt.Sprintf(
			"%s GROUP BY %s",
			activeQuery.sqlTemplate,
			activeQuery.buildGroupBy())

		// build having
		// waiting to do
	}

	// build order
	if activeQuery.orderBy != nil {
		order := activeQuery.buildOrderBy()
		if order != "" {
			activeQuery.sqlTemplate = fmt.Sprintf("%s ORDER BY %s", activeQuery.sqlTemplate, order)
		}
	}

	// build limit
	if activeQuery.rowsLength != 0 {
		activeQuery.sqlTemplate = fmt.Sprintf(
			"%s LIMIT %d, %d",
			activeQuery.sqlTemplate,
			activeQuery.start,
			activeQuery.rowsLength)
	}

	return activeQuery.sqlTemplate
}

func (activeQuery *ActiveQuery) buildField() string  {
	if field, ok := activeQuery.field.(string); ok {
		return field
	}

	if field, ok := activeQuery.field.(func(*ActiveQuery) string); ok {
		return field(activeQuery)
	}

	if field, ok := activeQuery.field.([]string); ok {
		return strings.Join(field, ",")
	}

	panic("invalid field")
}

func (activeQuery *ActiveQuery) buildCondition(where interface{}) string  {
	// 为string时的处理方式
	if where, ok := where.(string); ok {
		return fmt.Sprintf("(%s)", where)
	}

	// 为函数时的处理方式
	if handler, ok := where.(func(*ActiveQuery) string); ok {
		ret := handler(activeQuery)
		return fmt.Sprintf("(%s)", ret)
	}

	// 为map时的处理方式
	if where, ok := where.(map[string]interface{}); ok {
		var ret []string
		for key, value := range where {
			// map[string]string
			if value, ok := value.(string); ok {
				_key := getWhereParamsKey(activeQuery.whereParams, key)
				ret = append(ret, fmt.Sprintf("(%s={{%s}})", key, _key))
				activeQuery.whereParams[_key] = value
				continue
			}

			if value, ok := value.(int); ok {
				_key := getWhereParamsKey(activeQuery.whereParams, key)
				ret = append(ret, fmt.Sprintf("(%s={{%s}})", key, _key))
				activeQuery.whereParams[_key] = value
				continue
			}

			if value, ok := value.(float32); ok {
				_key := getWhereParamsKey(activeQuery.whereParams, key)
				ret = append(ret, fmt.Sprintf("(%s={{%s}})", key, _key))
				activeQuery.whereParams[_key] = value
				continue
			}

			if value, ok := value.(float64); ok {
				_key := getWhereParamsKey(activeQuery.whereParams, key)
				ret = append(ret, fmt.Sprintf("(%s={{%s}})", key, _key))
				activeQuery.whereParams[_key] = value
				continue
			}

			// callback
			if handler, ok := value.(func(*ActiveQuery) string); ok {
				_key := getWhereParamsKey(activeQuery.whereParams, key)
				ret = append(ret, fmt.Sprintf("(%s={{%s}})", key, _key))
				activeQuery.whereParams[_key] = handler(activeQuery)
				continue
			}

			panic("invalid condition item in while condition is map")
		}

		_ret := strings.Join(ret, " and ")
		return fmt.Sprintf("(%s)", _ret)
	}

	// 为数组时的处理方式
	if where, ok := where.([]interface{}); ok {
		var ret []string
		length := len(where)
		if length < 1 {
			return ""
		}

		operator := "and"
		if conditionOperator, ok := where[0].(string); ok && utils.InArray(conditionOperator, []interface{}{"or", "and"}) {
			operator = conditionOperator
			where = where[1:]
		}

		for _, item := range where {
			// nil
			if item == nil {
				continue
			}

			// string
			if item, ok := item.(string); ok {
				ret = append(ret, fmt.Sprintf("(%s)", item))
				continue
			}

			// map
			if item, ok := item.(map[string]interface{}); ok {
				ret = append(ret, activeQuery.buildCondition(item))
				continue
			}

			// callback
			if item, ok := item.(func(*ActiveQuery) string); ok {
				ret = append(ret, activeQuery.buildCondition(item))
				continue
			}

			// []interface{}
			if item, ok := item.([]interface{}); ok {
				// item[0] in {"or", "and"}
				item_0, ok := item[0].(string)
				if ok && utils.InArray(item_0, []interface{}{"or", "and"}) {
					ret = append(ret, activeQuery.buildCondition(item))
					continue
				}

				if len(item) < 3 {
					panic("invalid condition item while condition is array")
				}

				item_1, ok := item[1].(string)
				if !ok {
					panic("item 1 should be string")
				}

				if !utils.InArray(item_0, []interface{}{
					"=",
					"<>",
					"!=",
					">=",
					"<=",
					">",
					"<",
					"like",
					"LIKE",
					"not like",
					"NOT LIKE",
					"IN",
					"in",
					"NOT IN",
					"not in",
				}) {
					panic("invalid item 0")
				}

				_key := getWhereParamsKey(activeQuery.whereParams, item_1)
				// item_0 is operator, "=", ">", "<", "like", "in" ...
				// string
				if item_2, ok := item[2].(string); ok {
					if item_2 == "" {
						continue
					}

					ret = append(ret, fmt.Sprintf("(%s %s {{%s}})", item_1, item_0, _key))
					activeQuery.whereParams[_key] = item_2
					continue
				}

				// int
				if item_2, ok := item[2].(int); ok {
					ret = append(ret, fmt.Sprintf("(%s %s {{%s}})", item_1, item_0, _key))
					activeQuery.whereParams[_key] = item_2
					continue
				}

				// float
				if item_2, ok := item[2].(float64); ok {
					ret = append(ret, fmt.Sprintf("(%s %s {{%s}})", item_1, item_0, _key))
					activeQuery.whereParams[_key] = item_2
					continue
				}

				if item_2, ok := item[2].(float32); ok {
					ret = append(ret, fmt.Sprintf("(%s %s {{%s}})", item_1, item_0, _key))
					activeQuery.whereParams[_key] = item_2
					continue
				}

				if item_2, ok := item[2].(func(*ActiveQuery) string); ok {
					value := item_2(activeQuery)
					if value == "" {
						continue
					}

					ret = append(ret, fmt.Sprintf("(%s %s {{%s}})", item_1, item_0, _key))
					activeQuery.whereParams[_key] = value
					continue
				}

				// item_0 "in", "not in", so item[2] is array
				if !utils.InArray(item_0, []interface{}{"in", "not in"}) {
					panic("invalid item 0")
				}

				if item_2, ok := item[2].([]string); ok {
					item_2 := strings.Join(item_2, ",")
					ret = append(ret, fmt.Sprintf("(%s %s ({{%s}}))", item_1, item_0, _key))
					activeQuery.whereParams[_key] = item_2
					continue
				}

				if item_2, ok := item[2].([]int); ok {
					item_2 := utils.NumberArrayToString(item_2, ",")
					ret = append(ret, fmt.Sprintf("(%s %s ({{%s}}))", item_1, item_0, _key))
					activeQuery.whereParams[_key] = item_2
					continue
				}

				if item_2, ok := item[2].([]float32); ok {
					item_2 := utils.NumberArrayToString(item_2, ",")
					ret = append(ret, fmt.Sprintf("(%s %s ({{%s}}))", item_1, item_0, _key))
					activeQuery.whereParams[_key] = item_2
					continue
				}

				if item_2, ok := item[2].([]float64); ok {
					item_2 := utils.NumberArrayToString(item_2, ",")
					ret = append(ret, fmt.Sprintf("(%s %s ({{%s}}))", item_1, item_0, _key))
					activeQuery.whereParams[_key] = item_2
					continue
				}
			}

			// ![]interface{}
			panic("invalid condition item while condition nor of nil, string, array, map, function")
		}

		_ret := strings.Join(ret, fmt.Sprintf(" %s ", operator))
		return fmt.Sprintf("(%s)", _ret)
	}

	panic("invalid where condition")
}

func (activeQuery *ActiveQuery) buildConditionTemplate(where interface{}) string  {
	// 为string时的处理方式
	if where, ok := where.(string); ok {
		return fmt.Sprintf("(%s)", where)
	}

	// 为函数时的处理方式
	if handler, ok := where.(func(*ActiveQuery) string); ok {
		ret := handler(activeQuery)
		return fmt.Sprintf("(%s)", ret)
	}

	// 为map时的处理方式
	if where, ok := where.(map[string]interface{}); ok {
		var ret []string
		for key, value := range where {
			// map[string]string
			if value, ok := value.(string); ok {
				ret = append(ret, fmt.Sprintf("(%s='%s')", key, value))
				continue
			}

			if value, ok := value.(int); ok {
				ret = append(ret, fmt.Sprintf("(%s=%d)", key, value))
				continue
			}

			if value, ok := value.(float32); ok {
				ret = append(ret, fmt.Sprintf("(%s=%f)", key, value))
				continue
			}

			if value, ok := value.(float64); ok {
				ret = append(ret, fmt.Sprintf("(%s=%f)", key, value))
				continue
			}

			// callback
			if handler, ok := value.(func(*ActiveQuery) string); ok {
				ret = append(ret, fmt.Sprintf("(%s=%s)", key, handler(activeQuery)))
				continue
			}

			panic("invalid condition item in while condition is map")
		}

		_ret := strings.Join(ret, " and ")
		return fmt.Sprintf("(%s)", _ret)
	}

	// 为数组时的处理方式
	if where, ok := where.([]interface{}); ok {
		var ret []string
		length := len(where)
		if length < 1 {
			return ""
		}

		operator := "and"
		if conditionOperator, ok := where[0].(string); ok && utils.InArray(conditionOperator, []interface{}{"or", "and"}) {
			operator = conditionOperator
			where = where[1:]
		}

		for _, item := range where {
			// nil
			if item == nil {
				continue
			}

			// string
			if item, ok := item.(string); ok {
				ret = append(ret, fmt.Sprintf("(%s)", item))
				continue
			}

			// map
			if item, ok := item.(map[string]interface{}); ok {
				ret = append(ret, activeQuery.buildCondition(item))
				continue
			}

			// callback
			if item, ok := item.(func(*ActiveQuery) string); ok {
				ret = append(ret, activeQuery.buildCondition(item))
				continue
			}

			// []interface{}
			if item, ok := item.([]interface{}); ok {
				// item[0] in {"or", "and"}
				item_0, ok := item[0].(string)
				if ok && utils.InArray(item_0, []interface{}{"or", "and"}) {
					ret = append(ret, activeQuery.buildCondition(item))
					continue
				}

				if len(item) < 3 {
					panic("invalid condition item while condition is array")
				}

				item_1, ok := item[1].(string)
				if !ok {
					panic("item 1 should be string")
				}

				if !utils.InArray(item_0, []interface{}{
					"=",
					"<>",
					"!=",
					">=",
					"<=",
					">",
					"<",
					"like",
					"LIKE",
					"not like",
					"NOT LIKE",
					"IN",
					"in",
					"NOT IN",
					"not in",
				}) {
					panic("invalid item 0")
				}

				// item_0 is operator, "=", ">", "<", "like", "in" ...
				// string
				if item_2, ok := item[2].(string); ok {
					if item_2 == "" {
						continue
					}

					ret = append(ret, fmt.Sprintf("(%s %s '%s')", item_1, item_0, item_2))
					continue
				}

				// int
				if item_2, ok := item[2].(int); ok {
					ret = append(ret, fmt.Sprintf("(%s %s %d)", item_1, item_0, item_2))
					continue
				}

				// float
				if item_2, ok := item[2].(float64); ok {
					ret = append(ret, fmt.Sprintf("(%s %s %f)", item_1, item_0, item_2))
					continue
				}

				if item_2, ok := item[2].(float32); ok {
					ret = append(ret, fmt.Sprintf("(%s %s %f)", item_1, item_0, item_2))
					continue
				}

				if item_2, ok := item[2].(func(*ActiveQuery) string); ok {
					value := item_2(activeQuery)
					if value == "" {
						continue
					}

					ret = append(ret, fmt.Sprintf("(%s %s '%s')", item_1, item_0, value))
					continue
				}

				// item_0 "in", "not in", so item[2] is array
				if !utils.InArray(item_0, []interface{}{"in", "not in"}) {
					panic("invalid item 0")
				}

				if item_2, ok := item[2].([]string); ok {
					item_2 := strings.Join(item_2, ",")
					ret = append(ret, fmt.Sprintf("(%s %s (%s))", item_1, item_0, item_2))
					continue
				}

				if item_2, ok := item[2].([]int); ok {
					item_2 := utils.NumberArrayToString(item_2, ",")
					ret = append(ret, fmt.Sprintf("(%s %s (%s))", item_1, item_0, item_2))
					continue
				}

				if item_2, ok := item[2].([]float32); ok {
					item_2 := utils.NumberArrayToString(item_2, ",")
					ret = append(ret, fmt.Sprintf("(%s %s (%s))", item_1, item_0, item_2))
					continue
				}

				if item_2, ok := item[2].([]float64); ok {
					item_2 := utils.NumberArrayToString(item_2, ",")
					ret = append(ret, fmt.Sprintf("(%s %s (%s))", item_1, item_0, item_2))
					continue
				}
			}

			// ![]interface{}
			panic("invalid condition item while condition nor of nil, string, array, map, function")
		}

		_ret := strings.Join(ret, fmt.Sprintf(" %s ", operator))
		return fmt.Sprintf("(%s)", _ret)
	}

	panic("invalid where condition")
}

func (activeQuery *ActiveQuery) buildOrderBy() string  {
	var ret string
	for field, order := range activeQuery.orderBy {
		if !utils.InArray(order, []interface{}{ASC, DESC}) {
			panic("invalid order")
		}

		if ret == "" {
			ret = fmt.Sprintf("%s %s", field, order)
			continue
		}

		ret += fmt.Sprintf(",%s %s", field, order)
	}

	return ret
}

func (activeQuery *ActiveQuery) buildGroupBy() string  {
	return strings.Join(activeQuery.groupBy, ",")
}

func (activeQuery *ActiveQuery) GetSqlTemplate() string  {
	return activeQuery.sqlTemplate
}

func getWhereParamsKey(whereParams map[string]interface{}, key string) string  {
	if _, ok := whereParams[key]; !ok {
		return key
	}

	for {
		key += strconv.Itoa(rand.Intn(1000000000))
		if _, ok := whereParams[key]; !ok {
			return key
		}
	}
}