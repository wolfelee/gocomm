package jxorm

import (
	"path/filepath"
	"testing"
)

func TestConvertDataType(t *testing.T) {
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
	//t.Log(tables)
	matchTables := make(map[string]*Table)
	for _, item := range tables {
		match, err := filepath.Match("user", item)
		if err != nil {
			t.Fatal(err)
		}

		if !match {
			continue
		}

		columnData, err := e.FindColumns("test", item)
		if err != nil {
			t.Fatal(err)
		}
		//t.Log(columnData)
		table, err := columnData.Convert()
		if err != nil {
			t.Fatal(err)
		}

		matchTables[item] = table
	}

	if len(matchTables) == 0 {
		t.Fatal("no tables matched")
	}

	//generator, err := gen.NewDefaultGenerator("./")
	//if err != nil {
	//	t.Fatal(err)
	//}
	//t.Log(matchTables["user"].PrimaryKey.DbColumn)
	gt, err := ConvertDataType(matchTables["user"])
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("%#v", &gt.PrimaryKey.Name)

}
