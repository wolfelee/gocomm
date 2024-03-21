package jxorm

import (
	"fmt"
	"strings"
)

var commonMysqlDataTypeMap = map[string]string{
	// For consistency, all integer types are converted to int64
	// number
	"bool":      "int",
	"boolean":   "bool",
	"tinyint":   "int",
	"smallint":  "int",
	"mediumint": "int",
	"int":       "int",
	"integer":   "int",
	"bigint":    "int",
	"float":     "float64",
	"double":    "float64",
	"decimal":   "float64",
	// date&time
	"date":      "time.Time",
	"datetime":  "time.Time",
	"timestamp": "time.Time",
	"time":      "string",
	"year":      "int",
	// string
	"char":       "string",
	"varchar":    "string",
	"binary":     "string",
	"varbinary":  "string",
	"tinytext":   "string",
	"text":       "string",
	"mediumtext": "string",
	"longtext":   "string",
	"enum":       "string",
	"set":        "string",
	"json":       "string",
}

// ConvertDataToType converts mysql column type into golang type
func ConvertDataToType(dataBaseType string, isDefaultNull bool) (string, error) {
	tp, ok := commonMysqlDataTypeMap[strings.ToLower(dataBaseType)]
	if !ok {
		return "", fmt.Errorf("unexpected database type: %s", dataBaseType)
	}

	return mayConvertNullType(tp, isDefaultNull), nil
}

func mayConvertNullType(goDataType string, isDefaultNull bool) string {
	if !isDefaultNull {
		return goDataType
	}

	switch goDataType {
	case "int64":
		return "int"
	case "int32":
		return "int32"
	case "float64":
		return "float64"
	case "bool":
		return "bool"
	case "string":
		return "string"
	case "time.Time":
		return "time.Time"
	default:
		return goDataType
	}
}
