package sqlsrv

import (
	"errors"
	"fmt"
	"github.com/go-gourd/gourd/config"
	"github.com/go-gourd/gourd/logger"
	"gorm.io/driver/sqlserver"
	"gorm.io/gorm"
	gormLogger "gorm.io/gorm/logger"
	"time"
)

var dbs = make(map[string]*gorm.DB)

type LogWriter struct{}

func (w LogWriter) Printf(format string, args ...any) {
	logger.Errorf(format, args...)
}

// GetDb 获取数据库连接
func GetDb(name string) (*gorm.DB, error) {

	if _, ok := dbs[name]; ok {
		return dbs[name], nil
	}

	dbConf := config.GetDbConfig()
	if _, ok := dbConf[name]; !ok {
		return nil, errors.New("Database config '" + name + "' does not exist.")
	}

	conf := dbConf[name]

	//判断配置数据库类型
	if conf.Type != "sqlsrv" && conf.Type != "sqlserver" {
		return nil, errors.New("Database config '" + name + "' type is not sqlserver.")
	}

	dsnParam := ""
	if conf.Param != "" {
		dsnParam = "&" + conf.Param
	}
	dsnF := "sqlserver://%s:%s@%s:%d?database=%s%s"
	dsn := fmt.Sprintf(dsnF, conf.User, conf.Pass, conf.Host, conf.Port, conf.Database, dsnParam)

	// 慢日志阈值
	slowLogTime := conf.SlowLogTime
	if slowLogTime == 0 {
		slowLogTime = 60000 //默认1分钟
	}

	newLogger := gormLogger.New(
		LogWriter{},
		gormLogger.Config{
			SlowThreshold:             time.Duration(slowLogTime) * time.Millisecond, // 慢 SQL 阈值
			LogLevel:                  gormLogger.Warn,                               // 日志级别
			IgnoreRecordNotFoundError: true,                                          // 忽略ErrRecordNotFound（记录未找到）错误
			Colorful:                  false,                                         // 禁用彩色打印
		},
	)

	// 连接数据库
	newDb, err := gorm.Open(sqlserver.Open(dsn), &gorm.Config{
		Logger: newLogger,
	})
	if err != nil {
		logger.Error("cannot establish db connection: %w" + err.Error())
		return nil, err
	}

	dbs[name] = newDb

	return newDb, nil
}
