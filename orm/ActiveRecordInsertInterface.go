package orm

type ActiveRecordInsertInterface interface {
	Insert() (int64, bool)
}
