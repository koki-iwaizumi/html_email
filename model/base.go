package model

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/go-xorm/xorm"
)

//データベースエンジン
var Engine *xorm.Engine

const (
	MODEL_DRIVER   = "mysql"
	MODEL_USER     = "*"
	MODEL_PASSWORD = "*"
	MODEL_NAME     = "*"
)

//データベース接続
func Connect() (err error) {

	//データベース接続
	Engine, err = xorm.NewEngine(MODEL_DRIVER, MODEL_USER+":"+MODEL_PASSWORD+"@/"+MODEL_NAME+"?parseTime=true")

	return err
}
