package datasource

import (
	"fmt"
	"github.com/wolfelee/gocomm/cmd/et/internal/utils/config"
	"testing"
)

func TestFromDataSource(t *testing.T) {
	url := "web_app:WEB_app~!@`12@tcp(192.168.1.111:3306)/coursemanagerdev"
	table = "test123"
	dir = "./user"
	//cache := true
	cfg, err := config.NewConfig("")
	if err != nil {
		fmt.Println("生成config错误!")
		return
	}

	err = fromDataSource(url, table, dir, cfg, false)
	if err != nil {
		fmt.Println("生成文件错误", err)
	}

}
