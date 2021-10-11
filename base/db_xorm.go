package base

import (
	//_ "github.com/go-sql-driver/mysql"
	"github.com/lib/pq"
	"github.com/go-xorm/xorm"
	"github.com/mszhangyi/infra"
	"github.com/sirupsen/logrus"
	"time"
)

var orm *xorm.Engine

func OrmDatabase() *xorm.Engine {
	Check(orm)
	return orm
}

type DatabaseStarter struct {
	infra.BaseStarter
}

func (s *DatabaseStarter) Setup() {
	//数据库配置
	//engine, err := xorm.NewEngine("mysql", props.DataSource)
	engine, err := xorm.NewEngine("postgres", props.DataSource)
	if err != nil {
		logrus.Panicf("数据库orm：", err)
	}
	logrus.Debug("db_orm 开始启动了", engine.Ping())
	//设置连接复用时间
	engine.SetConnMaxLifetime(30 * time.Second)
	engine.SetMaxOpenConns(5000)
	engine.SetMaxIdleConns(100)
	orm = engine
}
