package db

import (
	//"github.com/go-xorm/core"
	"fmt"
	"mint-platform/platform/utils"

	_ "github.com/go-sql-driver/mysql"
	"github.com/go-xorm/xorm"
)

var engine *xorm.Engine

var engineMap map[string]*xorm.Engine

func NewEngineMap() {
	engineMap = make(map[string]*xorm.Engine)
}

func Init(dbname string) {
	var err error

	dbConfig, err := utils.LoadConfig("database_test", dbname)
	//	dbConfig, err := utils.LoadConfig("database", dbname)
	if err != nil {
		utils.WriteLog("log_config", "logconfig err,err:", err)
	}
	dbinfo := fmt.Sprintf("%s:%s@tcp(%s:3306)/%s?charset=%s", dbConfig["username"], dbConfig["password"], dbConfig["hostname"], dbConfig["database"], dbConfig["charset"])
	//engine, err = ConnectDB(dbinfo, dbname)
	err = ConnectDB(dbinfo, dbname)
	if err != nil {
		utils.WriteLog("log_db", "connect db err, err:", err)
	}
}

/*
//ConnectDB ...
func ConnectDB(dataSourceName string) (engin *xorm.Engine, err error) {
	if engine == nil {
		engine, err = xorm.NewEngine("mysql", dataSourceName)
		//("postgres", "user=test password=test dbname=test sslmode=disable")
		//"postgres", "postgres://postgres:beta@localhost:5432/public"
		if err != nil {
			utils.WriteLog("log_db", err)
			return nil, err
		}
		engine.Ping()
		engine.ShowSQL(true)
        //fmt.Println(engineMap)
		//		engine.Logger().SetLevel(core.LOG_DEBUG)
		//		prefixMapper := core.NewPrefixMapper(core.GonicMapper{}, "beta_")
		//		engine.SetTableMapper(prefixMapper)
	}
	return  engine, nil
}
*/

//ConnectDB ...
func ConnectDB(dataSourceName string, dbname string) (err error) {
	if engineMap[dbname] == nil {
		engine, err = xorm.NewEngine("mysql", dataSourceName)
		//("postgres", "user=test password=test dbname=test sslmode=disable")
		//"postgres", "postgres://postgres:beta@localhost:5432/public"
		if err != nil {
			utils.WriteLog("log_db", err)
			return err
		}
		engine.Ping()
		//engine.ShowSQL(true)
		engineMap[dbname] = engine
		//fmt.Println(engineMap)
		//		engine.Logger().SetLevel(core.LOG_DEBUG)
		//		prefixMapper := core.NewPrefixMapper(core.GonicMapper{}, "beta_")
		//		engine.SetTableMapper(prefixMapper)
	}
	return nil
}

//GetDB ...
func GetDB(dbname string) *xorm.Engine {
	return engineMap[dbname]
}
