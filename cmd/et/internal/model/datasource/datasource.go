package datasource

import (
	"errors"
	"fmt"
	"github.com/go-sql-driver/mysql"
	"github.com/spf13/cobra"
	"github.com/wolfelee/gocomm/cmd/et/internal/base"
	"github.com/wolfelee/gocomm/cmd/et/internal/model/gen"
	jxorm2 "github.com/wolfelee/gocomm/cmd/et/internal/model/jxorm"
	"github.com/wolfelee/gocomm/cmd/et/internal/utils/config"
	"path/filepath"
	"strings"
)

// CmdDatasource is connect mysql to generate code
var CmdDatasource = &cobra.Command{
	Use:   "datasource",
	Short: "datasource is connect mysql to generate code",
	Long:  "datasource is connect mysql to generate code. Example: et model datasource --url={datasource} --table={patterns}  --dir={dir} --cache=true",
	Run:   run,
}

var (
	url   string
	table string
	dir   string
	cache bool
	style string
)

func init() {
	CmdDatasource.Flags().StringVar(&url, "url", "", "mysql url")
	CmdDatasource.Flags().StringVar(&table, "table", "", "mysql tables:have access to '*'")
	CmdDatasource.Flags().StringVar(&dir, "dir", "", "generate model's dir")
	CmdDatasource.Flags().StringVar(&style, "style", "", "generate model's style")
	CmdDatasource.Flags().BoolVar(&cache, "cache", false, "Use cache?")
}

func run(cmd *cobra.Command, args []string) {
	if url == "" {
		fmt.Println("请使用 --url 输入mysql连接datasource")
	}
	if table == "" {
		fmt.Println("请使用 --table 输入mysql表")
	}
	if dir == "" {
		fmt.Println("请使用 --dir 输入生成代码的文件夹")
	}

	url = strings.TrimSpace(url)
	table = strings.TrimSpace(table)
	dir = strings.TrimSpace(dir)

	cfg, err := config.NewConfig(style)
	if err != nil {
		fmt.Println("生成config错误!")
		return
	}
	err = fromDataSource(url, table, dir, cfg, cache)
	if err != nil {
		fmt.Println("生成文件错误")
	}

	shells := [][]string{
		{"go", "mod", "tidy"},
		{"gofmt", "-w", "."},
	}

	err = base.FormatCode("./", shells...)
	if err != nil {
		fmt.Println("tidy or format filed", err)
	}
}

func fromDataSource(url, pattern, dir string, cfg *config.Config, cache bool) error {
	dsn, err := mysql.ParseDSN(url)
	if err != nil {
		return err
	}
	databaseSource := strings.TrimSuffix(url, "/"+dsn.DBName) + "/information_schema"

	eng, err := jxorm2.NewEngine(databaseSource)
	if err != nil {
		return err
	}

	tables, err := eng.GetAllTables(dsn.DBName)
	if err != nil {
		return err
	}

	matchTables := make(map[string]*jxorm2.Table)

	for _, item := range tables {
		match, err := filepath.Match(pattern, item)
		if err != nil {
			return err
		}

		if !match {
			continue
		}

		columnData, err := eng.FindColumns(dsn.DBName, item)
		if err != nil {
			return err
		}

		table, err := columnData.Convert()
		if err != nil {
			return err
		}

		matchTables[item] = table
	}

	if len(matchTables) == 0 {
		return errors.New("no tables matched")
	}

	generator, err := gen.NewDefaultGenerator(dir, cfg)
	if err != nil {
		return err
	}

	return generator.StartFromInformationSchema(matchTables, cache)
}
