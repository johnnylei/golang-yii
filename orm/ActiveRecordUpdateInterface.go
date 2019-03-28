package orm

type ActiveRecordUpdateInterface interface {
	Update() (int64, bool, error)
}
