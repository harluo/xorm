package config

import (
	"github.com/goexl/db"
	"github.com/harluo/xorm/internal/config/internal"
)

type Server struct {
	// 数据库类型
	Type db.Type `default:"sqlite3" json:"type,omitempty" validate:"required,oneof=mysql sqlite sqlite3 mssql oracle postgres postgresql"` // nolint:lll

	// 主机
	Host string `json:"host,omitempty" validate:"required,hostname|ip"`
	// 端口
	Port int `default:"3306" json:"port,omitempty" validate:"required,max=65535"`
	// 用户名
	Username string `json:"username,omitempty"`
	// 密码
	Password string `json:"password,omitempty"`
	// 连接协议
	Protocol string `default:"tcp" json:"protocol,omitempty" validate:"required,oneof=tcp udp"`

	// 连接池配置
	Connection Connection `json:"connection,omitempty"`
	// 参数配置
	Sqlite internal.Sqlite `json:"sqlite,omitempty"`
}
