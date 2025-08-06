package database

import (
	"fmt"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"org.thinkinai.com/recruit-center/pkg/config"
)

// Init 初始化数据库连接
func Init(cfg *config.DB) (*gorm.DB, error) {
	// 构建 DSN
	dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable TimeZone=Asia/Shanghai",
		cfg.Host,
		cfg.Port,
		cfg.Username,
		cfg.Password,
		cfg.Database,
	)

	// 配置 GORM
	gormConfig := &gorm.Config{
		// 设置日志级别
		Logger: logger.Default.LogMode(logger.Info),
		// 禁用默认事务以提高性能
		SkipDefaultTransaction: true,
		// 禁用外键约束
		DisableForeignKeyConstraintWhenMigrating: true,
	}

	// 建立连接
	db, err := gorm.Open(postgres.Open(dsn), gormConfig)
	if err != nil {
		return nil, fmt.Errorf("连接数据库失败: %w", err)
	}

	// 获取底层 *sql.DB 对象
	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("获取 *sql.DB 失败: %w", err)
	}

	// 配置连接池
	sqlDB.SetMaxIdleConns(10)           // 设置空闲连接池中的最大连接数
	sqlDB.SetMaxOpenConns(100)          // 设置打开数据库连接的最大数量
	sqlDB.SetConnMaxLifetime(time.Hour) // 设置连接可复用的最大时间

	return db, nil
}

// // AutoMigrate 自动迁移数据库表结构
// func AutoMigrate(db *gorm.DB) error {
// 	return db.AutoMigrate(
// 		&dao.Job{},
// 		&dao.JobApply{},
// 		// 添加其他需要迁移的模型
// 	)
// }

// Close 关闭数据库连接
func Close(db *gorm.DB) error {
	sqlDB, err := db.DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}
