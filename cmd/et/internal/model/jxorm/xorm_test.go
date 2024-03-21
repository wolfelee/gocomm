package jxorm

import (
	"testing"
)

func TestEngine_GetAllTables(t *testing.T) {
	e, err := NewEngine("root:zhy.1996@tcp(a.zhaohaiyu.com:3306)/information_schema")
	if err != nil {
		t.Fatal(err)
		return
	}

	tables, err := e.GetAllTables("test")
	if err != nil {
		t.Fatal(err)
		return
	}
	t.Log(tables)

	index, err := e.FindIndex("test", "user", "update_time")
	if err != nil {
		t.Log(err)
		return
	}
	t.Log(index)

	//cd, err := e.FindColumns("test", "user")
	//if err != nil {
	//	t.Fatal(err)
	//	return
	//}
	//t.Log(cd.Db)
	//t.Log(cd.Table)
	//for _, column := range cd.Columns {
	//	t.Log(column.DbColumn, column.Index)
	//}

	//table,err := cd.Convert()
	//if err != nil {
	//	t.Fatal(err)
	//	return
	//}
	//t.Logf("%#v",table)

}
