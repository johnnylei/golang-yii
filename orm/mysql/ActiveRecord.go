package mysql

import (
	"github.com/johnnylei/golang-yii/base"
	"github.com/johnnylei/golang-yii/orm"
	"reflect"
)

type ActiveRecord struct {
	base.Object
	Table string
	// 存储字段初始值
	OldFields map[string]interface{}
}

func (_self *ActiveRecord) Init() {
	base.On(base.END_LOAD_EVENT, _self, func(event base.EventInterface, object base.EventContainer) {
		instance := reflect.ValueOf(object).Elem().FieldByName("ActiveRecord")
		if !instance.IsValid() {
			panic("object should with field ActiveRecord")
		}

		record, ok := instance.Interface().(*ActiveRecord)
		if !ok {
			panic("object should with field ActiveRecord")
		}

		// OldFields，以便于后面update的时候，比较哪些数据产生了变化
		if len(record.OldFields) == 0 {
			record.OldFields = event.GetData()
		}
	})
}

func (_self *ActiveRecord) Find() orm.ActiveQueryInterface {
	activeQuery := NewDefaultActiveQuery()
	activeQuery.table = _self.Table
	return activeQuery
}

func (_self *ActiveRecord) Delete() bool {
	activeRecordDelete := NewDefaultActiveRecordDelete(_self)
	return activeRecordDelete.Delete()
}

func (_self *ActiveRecord) Insert() (int64, bool) {
	activeRecordInsert := NewDefaultActiveRecordInsert(_self)
	return activeRecordInsert.Insert()
}

func (_self *ActiveRecord) Update() (int64, bool, error)  {
	activeRecordInsert := NewDefaultActiveRecordUpdate(_self)
	return activeRecordInsert.Update()
}