package mysql

import (
	"icode.baidu.com/baidu/gdp/gdp/store/db"
	"sync"
)

var connector *db.DB
var lock sync.Mutex

func GetConnector() *db.DB  {
	if connector != nil {
		return connector
	}

	lock.Lock()
	defer lock.Unlock()
	connector = db.GetConn("clusterxdb")
	return connector
}

type Command struct {
	connector *db.DB
}

func NewDefaultCommand() *Command {
	return &Command{
		connector:GetConnector(),
	}
}

func (_self *Command) Query(sql string, parameters map[string]interface{}, fetchOne bool) []map[string]string  {
	rows, err := _self.connector.Query(sql, parameters)
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	data := make([]map[string]string, 0)
	cols, err := rows.Columns()
	if err != nil {
		panic(err)
	}
	columnsLen := len(cols)

	for rows.Next() {
		columnPoints := make([]interface{}, columnsLen)
		columnValues := make([]string, columnsLen)
		rowData := make(map[string]string, columnsLen)

		for i, _ := range columnValues {
			columnPoints[i] = &columnValues[i]
		}

		err = rows.Scan(columnPoints...)
		if err != nil {
			panic(err)
		}
		for i, columnName := range cols {
			rowData[columnName] = columnValues[i]
		}

		data = append(data, rowData)
		if fetchOne {
			break
		}
	}

	return data
}

func (_self *Command) Insert(tableName string, data db.Data) (int64, error)  {
	return _self.connector.Insert(tableName, data)
}

func (_self *Command) Exec(sql string, data db.Data) (db.Result, error)  {
	return _self.connector.Exec(sql, data)
}
