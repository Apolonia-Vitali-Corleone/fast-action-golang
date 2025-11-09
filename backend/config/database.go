package config

import (
	"fmt"
	"log"
	"time"

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

	// 获取底层的sql.DB对象以配置连接池
	sqlDB, err := DB.DB()
	if err != nil {
		log.Printf("获取数据库实例失败: %v", err)
		return err
	}

	// ========== 连接池配置 ==========
	// SetMaxIdleConns 设置空闲连接池中的最大连接数
	// 保持一定数量的空闲连接，避免频繁创建和销毁连接
	sqlDB.SetMaxIdleConns(10)

	// SetMaxOpenConns 设置打开数据库连接的最大数量
	// 限制最大连接数，防止数据库过载
	sqlDB.SetMaxOpenConns(100)

	// SetConnMaxLifetime 设置连接可复用的最大时间
	// 超过这个时间的连接将被关闭，避免使用过期连接
	sqlDB.SetConnMaxLifetime(time.Hour)

	// SetConnMaxIdleTime 设置连接空闲的最大时间
	// 空闲超过这个时间的连接将被关闭，释放资源
	sqlDB.SetConnMaxIdleTime(10 * time.Minute)

	// 连接成功
	log.Println("数据库连接成功，连接池已配置")
	return nil
}
