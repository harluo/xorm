package config

import (
	"github.com/goexl/db"
	"github.com/harluo/config"
)

type DB struct {
	Database `json:",squash" validate:"required_without=Databases"`

	// 其它数据库
	Databases map[string]*Database `json:"databases,omitempty" validate:"required_without=Database"`
}

func newDB(getter config.Getter) (db *DB, typ db.Type, err error) {
	db = new(DB)
	err = getter.Get(&struct {
		DB *DB `json:"db,omitempty" validate:"required"`
	}{
		DB: db,
	})
	if nil == err { // 暴露数据库类型，供其它项目使用
		typ = db.Type
	}

	return
}
