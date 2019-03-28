package orm

type ActiveQueryInterface interface {
	One() map[string]string
	All() []map[string]string

	Alias(string) ActiveQueryInterface
	Field(interface{}) ActiveQueryInterface
	Where(interface{}) ActiveQueryInterface
	Count() ActiveQueryInterface
	Limit(int, int) ActiveQueryInterface
	GroupBy([]string) ActiveQueryInterface
	Having(map[string]interface{}) ActiveQueryInterface
	OrderBy(map[string]string) ActiveQueryInterface
	LeftJoin(string, string) ActiveQueryInterface
	InnerJoin(string, string) ActiveQueryInterface
	RightJoin(string, string) ActiveQueryInterface
}