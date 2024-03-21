package jxorm

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"sort"
	"xorm.io/xorm"
)

const indexPri = "PRIMARY"

type MysqlEngine interface {
	GetAllTables(database string) ([]string, error)
	FindIndex(db, table, column string) ([]DbIndex, error)
	FindColumns(database, table string) (*ColumnData, error)
}

type engine struct {
	*xorm.Engine
}

type (
	DbColumn struct {
		Name            string `xorm:"COLUMN_NAME"`
		DataType        string `xorm:"DATA_TYPE"`
		Extra           string `xorm:"EXTRA"`
		Comment         string `xorm:"COLUMN_COMMENT"`
		ColumnDefault   string `xorm:"COLUMN_DEFAULT"`
		IsNullAble      string `xorm:"IS_NULLABLE"`
		OrdinalPosition int    `xorm:"ORDINAL_POSITION"`
	}
	DbIndex struct {
		IndexName  string `xorm:"INDEX_NAME"`
		NonUnique  int    `xorm:"NON_UNIQUE"`
		SeqInIndex int    `xorm:"SEQ_IN_INDEX"`
	}
	Column struct {
		*DbColumn
		Index *DbIndex
	}
	ColumnData struct {
		Db      string
		Table   string
		Columns []*Column
	}

	//  Table describes mysql table which contains database name, table name, columns, keys
	Table struct {
		Db      string
		Table   string
		Columns []*Column
		// Primary key not included
		UniqueIndex map[string][]*Column
		PrimaryKey  *Column
		NormalIndex map[string][]*Column
	}
)

func (*DbColumn) TableName() string {
	return "COLUMNS"
}

func (*DbIndex) TableName() string {
	return "STATISTICS"
}

func NewEngine(dataSourceName string) (MysqlEngine, error) {
	e, err := xorm.NewEngine("mysql", dataSourceName)
	if err != nil {
		return nil, err
	}
	return &engine{e}, nil
}

func (e *engine) GetAllTables(database string) ([]string, error) {
	var tables = make([]string, 0)
	m, err := e.Query("SELECT TABLE_NAME FROM INFORMATION_SCHEMA.TABLES WHERE TABLE_SCHEMA = ?", database)
	if err != nil {
		return nil, err
	}
	//fmt.Println(m)
	for _, table := range m {
		tables = append(tables, string(table["TABLE_NAME"]))
	}
	return tables, nil
}

// FindIndex finds index with given db, table and column.
func (e *engine) FindIndex(db, table, column string) ([]DbIndex, error) {
	var reply = make([]DbIndex, 0)
	err := e.Where("TABLE_SCHEMA = ? and TABLE_NAME = ? and COLUMN_NAME = ?", db, table, column).Find(&reply)
	if err != nil {
		return nil, err
	}

	return reply, nil
}

func (e *engine) FindColumns(database, table string) (*ColumnData, error) {
	var dbColumns = make([]DbColumn, 0)
	err := e.Where("TABLE_SCHEMA = ? and TABLE_NAME = ?", database, table).Find(&dbColumns)
	if err != nil {
		return nil, err
	}

	var list []*Column
	for _, item := range dbColumns {
		item := item
		index, err := e.FindIndex(database, table, item.Name)
		if err != nil {
			return nil, err
		}

		if len(index) > 0 {
			for _, i := range index {
				list = append(list, &Column{
					DbColumn: &item,
					Index:    &i,
				})
			}
		} else {
			list = append(list, &Column{
				DbColumn: &item,
			})
		}
	}
	sort.Slice(list, func(i, j int) bool {
		return list[i].OrdinalPosition < list[j].OrdinalPosition
	})

	var columnData ColumnData
	columnData.Db = database
	columnData.Table = table
	columnData.Columns = list
	return &columnData, nil
}

func (c *ColumnData) Convert() (*Table, error) {
	var table = Table{
		Db:          c.Db,
		Table:       c.Table,
		Columns:     c.Columns,
		UniqueIndex: map[string][]*Column{},
		PrimaryKey:  nil,
		NormalIndex: map[string][]*Column{},
	}

	m := make(map[string][]*Column)
	for _, each := range c.Columns {
		if each.Index != nil {
			m[each.Index.IndexName] = append(m[each.Index.IndexName], each)
		}
	}

	primaryColumns := m[indexPri]
	if len(primaryColumns) == 0 {
		return nil, fmt.Errorf("db:%s, table:%s, missing primary key", c.Db, c.Table)
	}

	if len(primaryColumns) > 1 {
		return nil, fmt.Errorf("db:%s, table:%s, joint primary key is not supported", c.Db, c.Table)
	}

	table.PrimaryKey = primaryColumns[0]
	for indexName, columns := range m {
		if indexName == indexPri {
			continue
		}

		for _, one := range columns {
			if one.Index != nil {
				if one.Index.NonUnique == 0 {
					table.UniqueIndex[indexName] = columns
				} else {
					table.NormalIndex[indexName] = columns
				}
			}
		}
	}
	return &table, nil
}
