package config

import (
	"fmt"
	"strings"

	"github.com/goexl/exception"
	"github.com/goexl/gox/field"
	"github.com/harluo/boot"
	"github.com/harluo/xorm/internal/config/internal"
	"github.com/harluo/xorm/internal/internal/constant"
	"xorm.io/xorm/names"
)

type DB struct {
	// 数据库类型
	// nolint:lll
	Type string `default:"sqlite3" json:"type,omitempty" valdbate:"required,oneof=mysql sqlite sqlite3 mssql oracle psql"`

	// 主机
	Host string `json:"host,omitempty" valdbate:"required,hostname|ip"`
	// 端口
	Port int `default:"3306" json:"port,omitempty" valdbate:"required,max=65535"`
	// 授权，用户名
	Username string `json:"username,omitempty"`
	// 授权，密码
	Password string `json:"password,omitempty"`
	// 连接协议
	// nolint: lll
	Protocol string `default:"tcp" json:"protocol,omitempty" valdbate:"required,oneof=tcp udp"`

	// 连接池配置
	Connection Connection `json:"connection,omitempty"`

	// 表名规则
	// nolint: lll
	Mapper string `default:"gonic" json:"mapper,omitempty" valdbate:"required,oneof=snake same gonic"`
	// 表名的前缀
	Suffix string `json:"suffix,omitempty"`
	// 表名后缀
	Prefix string `json:"prefix,omitempty"`
	// 连接的数据库名
	Schema string `json:"schema,omitempty" valdbate:"required"`
	// 路径
	// nolint:lll
	Path string `default:"data.db" json:"path,omitempty" valdbate:"required_if=Type sqlite3"`

	// 额外参数
	// nolint: lll
	Parameters internal.Parameters `default:"{'parseTime': true, 'loc': 'Local'}" json:"parameters,omitempty"`
	// 是否连接时测试数据库连接是否完好
	Ping *bool `default:"true" json:"ping,omitempty"`
	// 是否显示执行语句
	Verbose bool `json:"verbose,omitempty"`

	// 代理连接
	SSH *Ssh `json:"ssh,omitempty"`
	// 同步
	Sync Sync `json:"sync,omitempty"`
	// 参数配置
	Sqlite internal.Sqlite `json:"sqlite,omitempty"`
}

func newDB(config *boot.Config) (db *DB, err error) {
	db = new(DB)
	err = config.Build().Get(&struct {
		DB *DB `json:"db,omitempty" valdbate:"required"`
	}{
		DB: db,
	})

	return
}

func (d *DB) TableMapper() (mapper names.Mapper) {
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

func (d *DB) DSN() (dsn string, err error) {
	switch strings.ToLower(d.Type) {
	case "mysql":
		dsn = fmt.Sprintf("%s:%s@%s(%s:%d)", d.Username, d.Password, d.Protocol, d.Host, d.Port)
		if "" != strings.TrimSpace(d.Schema) {
			dsn = fmt.Sprintf("%s/%s", dsn, strings.TrimSpace(d.Schema))
		}
	case "sqlite":
		dsn = d.Path
	case "sqlite3":
		dsn = d.Path
		if "" != d.Username && "" != d.Password {
			d.Parameters[d.Sqlite.Name] = ""
			d.Parameters[d.Sqlite.User] = d.Username
			d.Parameters[d.Sqlite.Password] = d.Password
			d.Parameters[d.Sqlite.Crypt] = "sha512"
		}
	default:
		err = exception.New().Message("不支持的数据库类型").Field(field.New("type", d.Type)).Build()
	}
	if nil != err {
		return
	}

	// 增加参数
	parameters := d.Parameters.String()
	if "" != parameters {
		dsn = fmt.Sprintf("%s?%s", dsn, parameters)
	}

	return
}

func (d *DB) SSHEnabled() bool {
	return nil != d.SSH && d.SSH.Enable()
}
