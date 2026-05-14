package config

import (
	"fmt"
	"strings"

	"github.com/goexl/db"
	"github.com/goexl/exception"
	"github.com/goexl/gox/field"
	"github.com/harluo/xorm/internal/config/internal"
	"github.com/harluo/xorm/internal/internal/constant"
	"xorm.io/xorm/names"
)

type Database struct {
	Server `json:",squash" validate:"required"`

	// 从库列表
	Slaves map[string]*Server `json:"slaves,omitempty"` // 使用映射方便配置文件编写

	// 表名规则
	Mapper string `default:"gonic" json:"mapper,omitempty" validate:"required,oneof=snake same gonic"`
	// 表名的前缀
	Suffix string `json:"suffix,omitempty"`
	// 表名后缀
	Prefix string `json:"prefix,omitempty"`
	// 连接的数据库名
	Schema string `json:"schema,omitempty" validate:"required"`

	// 额外参数
	Parameters internal.Parameters `json:"parameters,omitempty"`
	// 是否连接时测试数据库连接是否完好
	Ping bool `json:"ping,omitempty"`
	// 是否显示执行语句
	Verbose bool `json:"verbose,omitempty"`

	// 代理连接
	SSH *Ssh `json:"ssh,omitempty"`
	// 同步
	Sync Sync `json:"sync,omitempty"`
}

func (d *Database) TableMapper() (mapper names.Mapper) {
	switch d.Mapper {
	case constant.Gonic:
		mapper = new(names.GonicMapper)
	case constant.Snake:
		mapper = new(names.SnakeMapper)
	case constant.Same:
		mapper = new(names.SameMapper)
	default:
		mapper = new(names.GonicMapper)
	}

	return
}

func (d *Database) SSHEnabled() bool {
	return nil != d.SSH && d.SSH.Enable()
}

func (d *Database) SN(server *Server) (sn string, err error) {
	switch server.Type {
	case db.TypeMySQL:
		sn = fmt.Sprintf("%s:%s@%s(%s:%d)", server.Username, server.Password, server.Protocol, server.Host, server.Port)
		if "" != strings.TrimSpace(d.Schema) {
			sn = fmt.Sprintf("%s/%s", sn, strings.TrimSpace(d.Schema))
		}
	case db.TypePostgres:
		sn = fmt.Sprintf(
			"host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
			server.Host, server.Port,
			server.Username, server.Password,
			d.Schema,
		)
	case db.TypeSQLite:
		sn = d.Schema
	case db.TypeSQLite3:
		sn = d.Schema
		if "" != server.Username && "" != server.Password {
			d.Parameters[d.Sqlite.Name] = ""
			d.Parameters[d.Sqlite.User] = d.Username
			d.Parameters[d.Sqlite.Password] = d.Password
			d.Parameters[d.Sqlite.Crypt] = "sha512"
		}
	default:
		err = exception.New().Message("不支持的数据库类型").Field(field.New("type", d.Type)).Build()
	}

	err = d.sslParam(d.Parameters) // 增加安全连接参数
	if nil == err && 0 != len(d.Parameters) { // 增加参数
		sn = d.Parameters.Merge(sn, d.Type)
	}

	return
}
