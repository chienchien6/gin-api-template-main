/**
 * @Description mysql迁移
 **/
package core

import (
	"RCSP/global"
	"RCSP/model"
	"fmt"

	"gorm.io/gorm"
)

// 设置表信息
func setTableOption(tableComment string) *gorm.DB {
	value := fmt.Sprintf("ENGINE=InnoDB COMMENT='%s'", tableComment)
	return global.GvaMysqlClient.Set("gorm:table_options", value)
}

// 用户相关表
func userTable() {
	// 用户账号表
	_ = setTableOption("用户表").AutoMigrate(&model.User{})
	// 用户信息表
	_ = setTableOption("用户信息表").AutoMigrate(&model.UserInfo{})
}

// 数据表迁移
func AutoMigrate() {
	// 创建用户相关表
	userTable()
}
