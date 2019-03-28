package mysql

import (
	"fmt"
	"reflect"
	"strings"
)

type ActiveRecordDelete struct {
	record *ActiveRecord
}

func NewDefaultActiveRecordDelete(record *ActiveRecord) *ActiveRecordInsert {
	return &ActiveRecordInsert{
		record:record,
	}
}

func (_self *ActiveRecordInsert) Delete() bool {
	sql, data := _self.buildSql()
	_, err := NewDefaultCommand().Exec(sql, data)
	if err != nil {
		panic(err)
	}
	return true
}

func (_self *ActiveRecordInsert) buildSql() (string, map[string]interface{}) {
	where, data :=  _self.buildWhereByPrimary()
	sql := fmt.Sprintf("DELETE FROM %s WHERE %s", _self.record.Table, where)
	return sql, data
}

func (_self *ActiveRecordInsert) buildWhereByPrimary() (string, map[string]interface{})  {
	instance := _self.record.Instance
	reflectValue := reflect.ValueOf(instance).Elem()
	count := reflectValue.Type().NumField()
	for i := 0; i < count; i++ {
		field := reflectValue.Type().Field(i)
		fieldName,ok := field.Tag.Lookup("field")
		if !ok || !strings.Contains(string(field.Tag), "primary") {
			continue
		}

		value := reflectValue.Field(i)
		switch field.Type.Kind() {
		case reflect.Int:
			return fmt.Sprintf("%s={{%s}}", fieldName, fieldName), map[string]interface{}{
				fieldName:value.Int(),
			}
		case reflect.String:
			return fmt.Sprintf("%s={{%s}}", fieldName, fieldName), map[string]interface{}{
				fieldName:value.String(),
			}
		case reflect.Float32:
			return fmt.Sprintf("%s={{%s}}", fieldName, fieldName), map[string]interface{}{
				fieldName:value.Float(),
			}
		case reflect.Float64:
			return fmt.Sprintf("%s={{%s}}", fieldName, fieldName), map[string]interface{}{
				fieldName:value.Float(),
			}
		}
	}

	panic("without primary key, could not be delete!!!")
}
