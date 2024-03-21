package jxorm

import (
	"fmt"
	"log"
	"sort"
	"strings"
)

const (
	unmanaged = iota
	untyped
	intType
	int64Type
	uintType
	uint64Type
	stringType
)

var Placeholder PlaceholderType

type (
	// GoTable describes a mysql table
	GoTable struct {
		Name        Strings
		PrimaryKey  Primary
		UniqueIndex map[string][]*Field
		NormalIndex map[string][]*Field
		Fields      []*Field
	}

	// Primary describes a primary key
	Primary struct {
		Field
		AutoIncrement bool
	}

	// Field describes a table field
	Field struct {
		Name            Strings
		DataBaseType    string
		DataType        string
		Comment         string
		SeqInIndex      int
		OrdinalPosition int
		MakeType        string
	}

	// KeyType types alias of int
	KeyType int

	// PlaceholderType represents a placeholder type.
	PlaceholderType = struct{}

	// Set is not thread-safe, for concurrent use, make sure to use it with synchronization. Set struct {
	Set struct {
		data map[interface{}]PlaceholderType
		tp   int
	}
)

// NewSet returns a managed Set, can only put the values with the same type.
func NewSet() *Set {
	return &Set{
		data: make(map[interface{}]PlaceholderType),
		tp:   untyped,
	}
}

// Contains checks if i is in s.
func (s *Set) Contains(i interface{}) bool {
	if len(s.data) == 0 {
		return false
	}

	s.validate(i)
	_, ok := s.data[i]
	return ok
}

// AddUint64 adds uint64 values ii into s.
func (s *Set) AddUint64(ii ...uint64) {
	for _, each := range ii {
		s.add(each)
	}
}

func (s *Set) add(i interface{}) {
	switch s.tp {
	case unmanaged:
		// do nothing
	case untyped:
		s.setType(i)
	default:
		s.validate(i)
	}
	s.data[i] = Placeholder
}

func (s *Set) setType(i interface{}) {
	// s.tp can only be untyped here
	switch i.(type) {
	case int:
		s.tp = intType
	case int64:
		s.tp = int64Type
	case uint:
		s.tp = uintType
	case uint64:
		s.tp = uint64Type
	case string:
		s.tp = stringType
	}
}

// AddStr adds string values ss into s.
func (s *Set) AddStr(ss ...string) {
	for _, each := range ss {
		s.add(each)
	}
}

func (s *Set) validate(i interface{}) {
	if s.tp == unmanaged {
		return
	}

	switch i.(type) {
	case int:
		if s.tp != intType {
			log.Printf("Error: element is int, but set contains elements with type %d", s.tp)
		}
	case int64:
		if s.tp != int64Type {
			log.Printf("Error: element is int64, but set contains elements with type %d", s.tp)
		}
	case uint:
		if s.tp != uintType {
			log.Printf("Error: element is uint, but set contains elements with type %d", s.tp)
		}
	case uint64:
		if s.tp != uint64Type {
			log.Printf("Error: element is uint64, but set contains elements with type %d", s.tp)
		}
	case string:
		if s.tp != stringType {
			log.Printf("Error: element is string, but set contains elements with type %d", s.tp)
		}
	}
}

// ConvertDataType converts mysql data type into golang data type
func ConvertDataType(table *Table) (*GoTable, error) {
	isPrimaryDefaultNull := table.PrimaryKey.ColumnDefault == "" && table.PrimaryKey.IsNullAble == "YES"
	primaryDataType, err := ConvertDataToType(table.PrimaryKey.DataType, isPrimaryDefaultNull)
	if err != nil {
		return nil, err
	}

	var reply GoTable
	reply.UniqueIndex = map[string][]*Field{}
	reply.NormalIndex = map[string][]*Field{}
	reply.Name = From(table.Table)
	seqInIndex := 0
	if table.PrimaryKey.Index != nil {
		seqInIndex = table.PrimaryKey.Index.SeqInIndex
	}

	reply.PrimaryKey = Primary{
		Field: Field{
			Name:            From(table.PrimaryKey.Name),
			DataBaseType:    table.PrimaryKey.DataType,
			DataType:        primaryDataType,
			Comment:         table.PrimaryKey.Comment,
			SeqInIndex:      seqInIndex,
			OrdinalPosition: table.PrimaryKey.OrdinalPosition,
		},
		AutoIncrement: strings.Contains(table.PrimaryKey.Extra, "auto_increment"),
	}

	fieldM := make(map[string]*Field)
	for _, each := range table.Columns {
		isDefaultNull := each.ColumnDefault == "" && each.IsNullAble == "YES"
		dt, err := ConvertDataToType(each.DataType, isDefaultNull)
		if err != nil {
			return nil, err
		}
		columnSeqInIndex := 0
		if each.Index != nil {
			columnSeqInIndex = each.Index.SeqInIndex
		}

		field := &Field{
			Name:            From(each.Name),
			DataBaseType:    each.DataType,
			DataType:        dt,
			Comment:         each.Comment,
			SeqInIndex:      columnSeqInIndex,
			OrdinalPosition: each.OrdinalPosition,
		}
		if strings.Contains(strings.ToLower(each.Extra), "on update current_timestamp") {
			field.MakeType = "updateTime"
		} else if each.ColumnDefault == "CURRENT_TIMESTAMP" {
			field.MakeType = "insertTime"
		}
		fieldM[each.Name] = field
	}

	for _, each := range fieldM {
		reply.Fields = append(reply.Fields, each)
	}
	sort.Slice(reply.Fields, func(i, j int) bool {
		return reply.Fields[i].OrdinalPosition < reply.Fields[j].OrdinalPosition
	})

	uniqueIndexSet := NewSet()
	for indexName, each := range table.UniqueIndex {
		sort.Slice(each, func(i, j int) bool {
			if each[i].Index != nil {
				return each[i].Index.SeqInIndex < each[j].Index.SeqInIndex
			}
			return false
		})

		if len(each) == 1 {
			one := each[0]
			if one.Name == table.PrimaryKey.Name {
				fmt.Printf("table %s: duplicate unique index with primary key, %s", table.Table, one.Name)
				continue
			}
		}

		var list []*Field
		var uniqueJoin []string
		for _, c := range each {
			list = append(list, fieldM[c.Name])
			uniqueJoin = append(uniqueJoin, c.Name)
		}

		uniqueKey := strings.Join(uniqueJoin, ",")
		if uniqueIndexSet.Contains(uniqueKey) {
			fmt.Printf("table %s: duplicate unique index, %s", table.Table, uniqueKey)
			continue
		}

		uniqueIndexSet.AddStr(uniqueKey)
		reply.UniqueIndex[indexName] = list
	}

	normalIndexSet := NewSet()
	for indexName, each := range table.NormalIndex {
		var list []*Field
		var normalJoin []string
		for _, c := range each {
			list = append(list, fieldM[c.Name])
			normalJoin = append(normalJoin, c.Name)
		}

		normalKey := strings.Join(normalJoin, ",")
		if normalIndexSet.Contains(normalKey) {
			log.Printf("table %s: duplicate index, %s", table.Table, normalKey)
			continue
		}

		normalIndexSet.AddStr(normalKey)
		sort.Slice(list, func(i, j int) bool {
			return list[i].SeqInIndex < list[j].SeqInIndex
		})

		reply.NormalIndex[indexName] = list
	}

	return &reply, nil
}

const timeImport = "time.Time"

// ContainsTime returns true if contains golang type time.Time
func (t *GoTable) ContainsTime() bool {
	for _, item := range t.Fields {
		if item.DataType == timeImport {
			return true
		}
	}
	return false
}

// KeysStr returns string keys in s.
func (s *Set) KeysStr() []string {
	var keys []string

	for key := range s.data {
		if strKey, ok := key.(string); !ok {
			continue
		} else {
			keys = append(keys, strKey)
		}
	}

	return keys
}
