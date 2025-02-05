package main

import (
	"github.com/0xweb-3/CoinNest/eth_srv/model"
	"gorm.io/gorm/schema"
	"log"
	"os"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func OpenDB() (*gorm.DB, error) {
	// 参考 https://github.com/go-sql-driver/mysql#dsn-data-source-name 获取详情
	dsn := "root:xin1234567890@tcp(127.0.0.1:3320)/coin_nest?charset=utf8mb4&parseTime=True&loc=Local"
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

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: newLogger, //设置全局的日志级别
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true, //去除表明后的s
		},
	})
	if err != nil {
		panic(err)
	}

	return db, err
}

func main() {
	db, err := OpenDB()
	if err != nil {
		panic(err)
	}
	// 迁移生成表
	err = db.AutoMigrate(&model.User{})
	if err != nil {
		panic(err)
	}

	for i := 0; i < 2; i++ {
		//now := time.Now()
		db.Create(&model.User{})
	}

	//fmt.Println(db)

	//var user model.User
	//db.First(&user)
	//
	//db.Model(&user).Update("Password", "xin")
	//fmt.Println(user.Password)

	//fmt.Println(crypto.CompareHash(user.Password, "123456"))

}
