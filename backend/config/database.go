package config

import (
	"fmt"
	"log"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// DB 全局数据库连接对象
var DB *gorm.DB

// 数据库配置结构体
type DBConfig struct {
	Host     string // 数据库主机地址
	Port     string // 数据库端口
	User     string // 数据库用户名
	Password string // 数据库密码
	DBName   string // 数据库名称
}

// InitDB 初始化数据库连接
// 参数: config - 数据库配置信息
// 返回: 错误信息（如果有）
func InitDB(config DBConfig) error {
	// 构建MySQL连接字符串（DSN - Data Source Name）
	// 格式: 用户名:密码@tcp(主机:端口)/数据库名?参数
	// parseTime=True: 自动将数据库中的datetime类型转换为Go的time.Time
	// loc=Local: 使用本地时区
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		config.User,
		config.Password,
		config.Host,
		config.Port,
		config.DBName,
	)

	// 使用GORM打开MySQL数据库连接
	var err error
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		// 连接失败，记录错误并返回
		log.Printf("数据库连接失败: %v", err)
		return err
	}

	// 连接成功
	log.Println("数据库连接成功")
	return nil
}
