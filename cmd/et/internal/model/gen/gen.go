package gen

import (
	"bytes"
	"fmt"
	jxorm "github.com/wolfelee/gocomm/cmd/et/internal/model/jxorm"
	"github.com/wolfelee/gocomm/cmd/et/internal/model/template"
	"github.com/wolfelee/gocomm/cmd/et/internal/utils"
	"github.com/wolfelee/gocomm/cmd/et/internal/utils/config"
	"github.com/wolfelee/gocomm/cmd/et/internal/utils/format"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

const pwd = "."

type (
	defaultGenerator struct {
		// source string
		dir string
		pkg string
		cfg *config.Config
	}

	// Option defines a function with argument defaultGenerator
	Option func(generator *defaultGenerator)

	code struct {
		importsCode string
		varsCode    string
		typesCode   string
		newCode     string
		insertCode  string
		findCode    []string
		updateCode  string
		deleteCode  string
		cacheExtra  string
		tableName   string
	}
)

func NewDefaultGenerator(dir string, cfg *config.Config, opt ...Option) (*defaultGenerator, error) {
	if dir == "" {
		dir = pwd
	}
	dirAbs, err := filepath.Abs(dir)
	if err != nil {
		return nil, err
	}
	dir = dirAbs
	pkg := filepath.Base(dirAbs)
	err = utils.MkdirIfNotExist(dir)
	if err != nil {
		return nil, err
	}

	generator := &defaultGenerator{dir: dir, pkg: pkg, cfg: cfg}
	var optionList []Option
	optionList = append(optionList, opt...)
	for _, fn := range optionList {
		fn(generator)
	}

	return generator, nil
}

func (g *defaultGenerator) StartFromInformationSchema(tables map[string]*jxorm.Table, withCache bool) error {
	m := make(map[string]string)
	for _, each := range tables {
		table, err := jxorm.ConvertDataType(each)
		if err != nil {
			return err
		}
		code, err := g.genModel(*table, withCache)
		if err != nil {
			return err
		}

		m[table.Name.Source()] = code
	}
	return g.createFile(m)
}

func (g *defaultGenerator) createFile(modelList map[string]string) error {
	dirAbs, err := filepath.Abs(g.dir)
	if err != nil {
		return err
	}

	g.dir = dirAbs
	g.pkg = filepath.Base(dirAbs)
	err = utils.MkdirIfNotExist(dirAbs)
	if err != nil {
		return err
	}

	for tableName, code := range modelList {
		tn := jxorm.From(tableName)
		modelFilename, err := format.FileNamingFormat(g.cfg.NamingFormat, fmt.Sprintf("%s_model", tn.Source()))
		if err != nil {
			return err
		}

		name := modelFilename + ".go"
		filename := filepath.Join(dirAbs, name)
		if utils.FileExists(filename) {
			fmt.Printf("%s already exists, ignored.\n", name)
			continue
		}
		err = ioutil.WriteFile(filename, []byte(code), os.ModePerm)
		if err != nil {
			return err
		}
	}

	fmt.Println("Done.")
	return nil
}

// Table defines mysql table
type Table struct {
	jxorm.GoTable
	PrimaryCacheKey        Key
	UniqueCacheKey         []Key
	ContainsUniqueCacheKey bool
}

func (g *defaultGenerator) genModel(in jxorm.GoTable, withCache bool) (string, error) {
	if len(in.PrimaryKey.Name.Source()) == 0 {
		return "", fmt.Errorf("table %s: missing primary key", in.Name.Source())
	}

	primaryKey, uniqueKey := genCacheKeys(in)

	importsCode, err := genImports(withCache, in.ContainsTime())
	if err != nil {
		return "", err
	}

	var table Table
	table.GoTable = in
	table.PrimaryCacheKey = primaryKey
	table.UniqueCacheKey = uniqueKey
	table.ContainsUniqueCacheKey = len(uniqueKey) > 0

	varsCode, err := genVars(table, withCache)
	if err != nil {
		return "", err
	}

	insertCode, err := genInsert(table, withCache)
	if err != nil {
		return "", err
	}

	var findCode = make([]string, 0)
	findOneCode, err := genFindOne(table, withCache)
	if err != nil {

		return "", err
	}

	findCode = append(findCode, findOneCode)

	updateCode, err := genUpdate(table, withCache)
	if err != nil {
		fmt.Println(err)
		return "", err
	}

	deleteCode, err := genDelete(table, withCache)
	if err != nil {
		return "", err
	}

	typesCode, err := genTypes(table, withCache)
	if err != nil {
		return "", err
	}

	newCode, err := genNew(table, withCache)
	if err != nil {
		return "", err
	}

	tableName, err := genTableName(table)
	if err != nil {
		return "", err
	}

	code := &code{
		importsCode: importsCode,
		varsCode:    varsCode,
		typesCode:   typesCode,
		newCode:     newCode,
		insertCode:  insertCode,
		findCode:    findCode,
		updateCode:  updateCode,
		deleteCode:  deleteCode,
		//cacheExtra:  ret.cacheExtra,
		tableName: tableName,
	}

	output, err := g.executeModel(code)
	if err != nil {
		return "", err
	}

	return output.String(), nil

}

func (g *defaultGenerator) executeModel(code *code) (*bytes.Buffer, error) {
	text, err := utils.LoadTemplate(category, modelTemplateFile, template.Model)
	if err != nil {
		return nil, err
	}
	t := utils.With("model").
		Parse(text).
		GoFmt(true)
	output, err := t.Execute(map[string]interface{}{
		"pkg":         g.pkg,
		"imports":     code.importsCode,
		"vars":        code.varsCode,
		"types":       code.typesCode,
		"new":         code.newCode,
		"insert":      code.insertCode,
		"find":        strings.Join(code.findCode, "\n"),
		"update":      code.updateCode,
		"delete":      code.deleteCode,
		"extraMethod": code.cacheExtra,
		"tableName":   code.tableName,
	})
	if err != nil {
		return nil, err
	}
	return output, nil
}

func wrapWithRawString(v string) string {
	if v == "`" {
		return v
	}

	if !strings.HasPrefix(v, "`") {
		v = "`" + v
	}

	if !strings.HasSuffix(v, "`") {
		v = v + "`"
	} else if len(v) == 1 {
		v = v + "`"
	}

	return v
}
