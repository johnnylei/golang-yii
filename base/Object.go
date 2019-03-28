package base

import (
	"icode.baidu.com/baidu/searchbox/wechat_api/libraries/utils"
	"reflect"
	"strconv"
)

type Object struct {
	Instance interface{}
	event *Event
}

func (_self *Object) GetEvent() EventInterface  {
	if _self.event == nil {
		_self.event = &Event{}
	}

	return _self.event
}

func (_self *Object) Init() {

}

func On(name string, object EventContainer, callback func(EventInterface, EventContainer))  {
	event := object.GetEvent()
	event.On(name, callback)
}

func Trigger(name string, object EventContainer)  {
	event := object.GetEvent()
	event.Trigger(name, object)
}

func HasField(object interface{}, fieldName string) bool  {
	ValueOfObject := reflect.ValueOf(object)

	// Check if the passed interface is a pointer
	if ValueOfObject.Type().Kind() != reflect.Ptr {
		// Create a new type of Iface's Type, so we have a pointer to work with
		ValueOfObject = reflect.New(reflect.TypeOf(object))
	}

	field := ValueOfObject.Elem().FieldByName(fieldName)
	if field.IsValid() {
		return true
	}

	return false
}

func HasMethod(object interface{}, methodName string) bool  {
	ValueOfObject := reflect.ValueOf(object)

	// Check if the passed interface is a pointer
	if ValueOfObject.Type().Kind() != reflect.Ptr {
		// Create a new type of Iface's Type, so we have a pointer to work with
		ValueOfObject = reflect.New(reflect.TypeOf(object))
	}

	method := ValueOfObject.Elem().MethodByName(methodName)
	if method.IsValid() {
		return true
	}

	return false
}

const (
	START_LOAD_EVENT = "START_LOAD_EVENT"
	END_LOAD_EVENT = "END_LOAD_EVENT"
)

func Load(instance interface{}, data map[string]interface{}) {
	if reflect.TypeOf(instance).Kind() != reflect.Ptr {
		panic("instance should be pointer")
	}

	eventContainer, eventOk := instance.(EventContainer)
	if eventOk {
		eventContainer.GetEvent().SetData(data)
		Trigger(START_LOAD_EVENT, eventContainer)

		defer func(eventContainer EventContainer) {
			Trigger(END_LOAD_EVENT, eventContainer)
		}(eventContainer)
	}

	_data, ok := utils.DeepCopy(data).(map[string]interface{})
	if !ok {
		panic("Copy failed")
	}
	reflectValue := reflect.ValueOf(instance).Elem()
	count := reflectValue.Type().NumField()
	for i := 0; i < count; i++ {
		if len(_data) == 0 {
			break
		}

		field := reflectValue.Field(i)
		if !field.CanSet() {
			continue
		}

		fieldStruct := reflectValue.Type().Field(i)
		key := fieldStruct.Name
		value, ok := _data[key]
		if !ok {
			key = fieldStruct.Tag.Get("field")
			if key == "" {
				continue
			}

			value, ok = _data[key]
			if !ok {
				continue
			}
		}
		delete(_data, key)

		switch field.Kind() {
		case reflect.Int:
			_value, ok := value.(int64)
			if !ok {
				__value, ok := value.(string)
				if !ok {
					break
				}

				_value, err := strconv.ParseInt(__value, 10, 64)
				if err != nil {
					break
				}
				field.SetInt(_value)
				break
			}
			field.SetInt(_value)
		case reflect.String:
			_value, ok := value.(string)
			if ok {
				field.SetString(_value)
			}
		case reflect.Bool:
			_value, ok := value.(bool)
			if !ok {
				__value, ok := value.(string)
				if !ok {
					break
				}

				_value, err := strconv.ParseBool(__value)
				if err != nil {
					break
				}
				field.SetBool(_value)
				break
			}
			field.SetBool(_value)
		case reflect.Float32:
			_value, ok := value.(float64)
			if !ok {
				__value, ok := value.(string)
				if !ok {
					break
				}

				_value, err := strconv.ParseFloat(__value, 64)
				if err != nil {
					break
				}
				field.SetFloat(_value)
				break
			}
			field.SetFloat(_value)
		case reflect.Float64:
			_value, ok := value.(float64)
			if !ok {
				__value, ok := value.(string)
				if !ok {
					break
				}

				_value, err := strconv.ParseFloat(__value, 64)
				if err != nil {
					break
				}
				field.SetFloat(_value)
				break
			}
			field.SetFloat(_value)
		case reflect.Complex64:
			_value, ok := value.(complex128)
			if ok {
				field.SetComplex(_value)
			}
		case reflect.Complex128:
			_value, ok := value.(complex128)
			if ok {
				field.SetComplex(_value)
			}
		}
	}
}