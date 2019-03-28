package mysql

import (
	"reflect"
)

type ActiveRecordInsert struct {
	record *ActiveRecord
}

func NewDefaultActiveRecordInsert(record *ActiveRecord) *ActiveRecordInsert {
	return &ActiveRecordInsert{
		record:record,
	}
}

func (_self *ActiveRecordInsert) Insert() (int64, bool) {
	command := NewDefaultCommand()
	data := _self.GetFields()
	affected, err := command.Insert(_self.record.Table, data)
	if err != nil {
		panic(err)
	}

	if affected > 0 {
		return affected, true
	}

	return 0, false
}

func (_self *ActiveRecordInsert) GetFields() map[string]interface{} {
	data := make(map[string]interface{})
	instance := _self.record.Instance
	reflectValue := reflect.ValueOf(instance).Elem()
	count := reflectValue.Type().NumField()
	for i := 0; i < count; i++ {
		field := reflectValue.Type().Field(i)
		fieldName := field.Tag.Get("field")
		if fieldName == "" {
			continue
		}

		value := reflectValue.Field(i)
		switch field.Type.Kind() {
		case reflect.Int:
			data[fieldName] = value.Int()
		case reflect.Float64:
			data[fieldName] = value.Float()
		case reflect.Float32:
			data[fieldName] = value.Float()
		case reflect.Bool:
			data[fieldName] = value.Bool()
		case reflect.String:
			data[fieldName] = value.String()
		}
	}

	return data
}
