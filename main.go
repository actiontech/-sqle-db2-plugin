package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/actiontech/sqle/sqle/pkg/params"
	"strings"

	"github.com/actiontech/sqle/sqle/driver"
	adaptor "github.com/actiontech/sqle/sqle/pkg/driver"
	_ "github.com/ibmdb/go_ibm_db"
)

var database_name string

type DB2Dialector struct{}

func (d *DB2Dialector) Dialect(dsn *driver.DSN) (string, string) {
	user, database := dsn.AdditionalParams.GetParam("test_name"), dsn.AdditionalParams.GetParam("database")
	us, da := "", ""
	if user != nil {
		us = user.Value
	}
	if database != nil {
		da = database.Value
	}

	con := fmt.Sprintf("HOSTNAME=%v;DATABASE=%v;PORT=%v;UID=%v;PWD=%v", dsn.Host, da, dsn.Port, us, dsn.Password)
	return "go_ibm_db", con
}

func (d *DB2Dialector) String() string {
	return "DB2-New"
}

func (d *DB2Dialector) ShowDatabaseSQL() string {
	return fmt.Sprintf("select '%v' as name from SYSIBM.SYSDUMMY1", database_name)
}

var version string
var printVersion = flag.Bool("version", false, "Print version & exit")

func main() {
	flag.Parse()

	if *printVersion {
		fmt.Println(version)
		return
	}

	plugin := adaptor.NewAdaptor(&DB2Dialector{})

	ruleDQL1 := &driver.Rule{
		Name:     "db2_dql_1",
		Desc:     "禁止不带where条件的查询",
		Category: "SQL规范",
		Level:    driver.RuleLevelError,
	}
	ruleDQL1Handler := func(ctx context.Context, rule *driver.Rule, sql string) (string, error) {
		lowerSql := strings.ToLower(sql)
		if strings.Contains(lowerSql, "select") &&
			!strings.Contains(lowerSql, "where") {
			return rule.Desc, nil
		}
		return "", nil
	}
	plugin.AddRule(ruleDQL1, ruleDQL1Handler)

	plugin.AddAdditionalParams(&params.Param{
		Key:   "test_name",
		Value: "",
		Desc:  "用户名(测试凑数用)",
		Type:  "string",
	})

	plugin.AddAdditionalParams(&params.Param{
		Key:   "database",
		Value: "",
		Desc:  "数据库名",
		Type:  "string",
	})

	plugin.Serve()
}
