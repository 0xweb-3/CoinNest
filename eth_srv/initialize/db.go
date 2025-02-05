package initialize

import (
	"fmt"
	"github.com/0xweb-3/CoinNest/eth_srv/global"
	"go.uber.org/zap"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
	"log"
	"os"
	"time"
)

func InitDB() {
	// 参考 https://github.com/go-sql-driver/mysql#dsn-data-source-name 获取详情
	//dsn := "root:xinbingliang@tcp(192.168.21.2:3310)/fishline?charset=utf8mb4&parseTime=True&loc=Local"
	cnf := global.ServerConfig.Mysql
	dsn := fmt.Sprintf(
		"%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		cnf.Username,
		cnf.Password,
		cnf.Host,
		cnf.Port,
		cnf.DbName,
	)

	zap.S().Debug(dsn)
	// 设置全局的logger，这个logger在我们执行每个sql语句的时候会打印每一行sql

	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // 日志输出的位置
		logger.Config{
			SlowThreshold: time.Second, // 慢sql的阀值
			LogLevel:      logger.Info, // Log level ；Silent、Error、Warn、Info；info 表示所有sql都会打印
			//IgnoreRecordNotFoundError: true,          // 忽略记录器的 ErrRecordNotFound 错误
			//ParameterizedQueries:      true,          // 不要在 SQL 日志中包含参数
			Colorful: true, // 是否禁用彩色打印
		},
	)
	var err error
	global.DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: newLogger, //设置全局的日志级别
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true, //去除表明后的s
		},
	})
	if err != nil {
		panic(err)
	}
}
