package orm

type ActiveRecordInterface interface {
	Find() ActiveQueryInterface
	Delete()
	Insert()
	Update()
}
