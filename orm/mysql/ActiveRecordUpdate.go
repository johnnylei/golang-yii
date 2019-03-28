package mysql

import (
	"errors"
	"fmt"
	"icode.baidu.com/baidu/searchbox/wechat_api/libraries/utils"
	"reflect"
	"strings"
)

type ActiveRecordUpdate struct {
	record *ActiveRecord
}

func NewDefaultActiveRecordUpdate(record *ActiveRecord) *ActiveRecordUpdate {
	return &ActiveRecordUpdate{
		record:record,
	}
}

func (_self *ActiveRecordUpdate) Update() (int64, bool, error)  {
	sql, data, err := _self.buildSql()
	if err != nil {
		return 0, false, err
	}
	// 没有变化的时候
	if data == nil {
		return 0, true, nil
	}

	result, err := NewDefaultCommand().Exec(sql, data)
	if err != nil {
		return 0, false, err
	}

	affected, err := result.RowsAffected()
	if err != nil {
		return 0, false, err
	}

	return affected, true, nil
}

func (_self *ActiveRecordUpdate) buildSql() (string, map[string]interface{}, error)  {
	where, primary, primaryData, err := _self.buildWhereByPrimary()
	if err != nil {
		return "", nil, err
	}

	updateContent, data := _self.buildUpdateContent()
	if data == nil {
		return "", nil, nil
	}

	data[primary] = primaryData[primary]

	sql := fmt.Sprintf("UPDATE `%s` SET %s WHERE %s", _self.record.Table, updateContent, where)
	return sql, data, nil
}

func (_self *ActiveRecordUpdate) buildUpdateContent() (string, map[string]interface{}) {
	instance := _self.record.Instance
	reflectValue := reflect.ValueOf(instance).Elem()

	oldFields := _self.record.OldFields
	count := reflectValue.Type().NumField()
	changed := make([]string, 0)
	data := make(map[string]interface{})
	for i := 0; i < count; i++ {
		field := reflectValue.Type().Field(i)
		fieldName,ok := field.Tag.Lookup("field")
		if !ok || strings.Contains(string(field.Tag), "primary") {
			continue
		}

		// 如果没有变化,则不写入update 内容里面
		value := reflectValue.Field(i)
		OldValue, ok := oldFields[fieldName]
		if ok && OldValue == value.Interface() {
			continue
		}

		switch field.Type.Kind() {
		case reflect.Int:
			_oldValue := utils.ConvertInt64(OldValue)
			_value := value.Int()
			if ok && _oldValue == _value {
				break
			}
			changed = append(changed, fmt.Sprintf("%s={{%s}}", fieldName, fieldName))
			data[fieldName] = _value
		case reflect.String:
			_oldValue, ok := OldValue.(string)
			_value := value.String()
			if ok && _oldValue == _value {
				break
			}
			changed = append(changed, fmt.Sprintf("%s={{%s}}", fieldName, fieldName))
			data[fieldName] = _value
		case reflect.Float32:
			_oldValue := utils.ConvertFloat64(OldValue)
			_value := value.Float()
			if ok && _oldValue == _value {
				break
			}
			changed = append(changed, fmt.Sprintf("%s={{%s}}", fieldName, fieldName))
			data[fieldName] = _value
		case reflect.Float64:
			_oldValue := utils.ConvertFloat64(OldValue)
			_value := value.Float()
			if ok && _oldValue == _value {
				break
			}
			changed = append(changed, fmt.Sprintf("%s={{%s}}", fieldName, fieldName))
			data[fieldName] = _value
		}
	}

	if len(changed) == 0 {
		return "", nil
	}

	return strings.Join(changed, ","), data
}

func (_self *ActiveRecordUpdate) buildWhereByPrimary() (string, string, map[string]interface{}, error)  {
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
			return fmt.Sprintf("%s={{%s}}", fieldName, fieldName), fieldName, map[string]interface{}{
				fieldName:value.Int(),
			}, nil
		case reflect.String:
			return fmt.Sprintf("%s={{%s}}", fieldName, fieldName), fieldName, map[string]interface{}{
				fieldName:value.String(),
			}, nil
		case reflect.Float32:
			return fmt.Sprintf("%s={{%s}}", fieldName, fieldName), fieldName, map[string]interface{}{
				fieldName:value.Float(),
			}, nil
		case reflect.Float64:
			return fmt.Sprintf("%s={{%s}}", fieldName, fieldName), fieldName, map[string]interface{}{
				fieldName:value.Float(),
			}, nil
		}
	}

	return "", "", nil, errors.New("without primary key, could not be delete")
}