package core

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/elliotchance/sshtunnel"
	"github.com/goexl/gox"
	"github.com/goexl/gox/field"
	"github.com/goexl/log"
	"github.com/harluo/xorm/internal/config"
	"github.com/harluo/xorm/internal/internal"
	"golang.org/x/crypto/ssh"
	"xorm.io/core"
	"xorm.io/xorm"
)

type Engine struct {
	*shadowEngine

	_ gox.Pointerized
}

func (e *Engine) ColumnName(field string) string {
	return fmt.Sprintf("`%s`", e.GetColumnMapper().Table2Obj(field))
}

func (e *Engine) Context(ctx context.Context) *Session {
	return &Session{
		shadowSession: e.shadowEngine.Context(ctx),
	}
}

func (e *Engine) Cols(columns ...string) *Session {
	return &Session{
		shadowSession: e.shadowEngine.Cols(columns...),
	}
}

func (e *Engine) Where(query any, args ...any) *Session {
	return &Session{
		shadowSession: e.shadowEngine.Where(query, args...),
	}
}

func (e *Engine) Desc(field string) *Session {
	return &Session{
		shadowSession: e.shadowEngine.Desc(e.ColumnName(field)),
	}
}

func (e *Engine) Asc(field string) *Session {
	return &Session{
		shadowSession: e.shadowEngine.Asc(e.ColumnName(field)),
	}
}

func (e *Engine) Id(id any) *Session {
	return &Session{
		shadowSession: e.shadowEngine.ID(id),
	}
}

func newEngine(db *config.DB, logger log.Logger) (engine *Engine, err error) {
	engine = new(Engine)
	if ese := enableSSH(db, logger); nil != ese {
		err = ese
	} else if dsn, de := db.SN(); nil != de {
		err = de
	} else if engine.shadowEngine, err = xorm.NewEngine(db.Type.String(), dsn); nil == err {
		err = setupEngine(db, engine, logger)
	}

	return
}

func setupEngine(db *config.DB, engine *Engine, logger log.Logger) (err error) {
	// 替换成统一的日志框架
	engine.SetLogger(internal.NewXorm(logger))
	// 调试模式下打开各种可调试的选项
	if db.Verbose {
		engine.ShowSQL()
	}

	// 配置数据库连接池
	engine.SetMaxOpenConns(db.Connection.Open)
	engine.SetMaxIdleConns(db.Connection.Idle)
	engine.SetConnMaxLifetime(db.Connection.Lifetime)

	// 设置名称转换
	mapper := db.TableMapper()
	core.NewCacheMapper(core.GonicMapper{})
	if "" != strings.TrimSpace(db.Prefix) {
		mapper = core.NewPrefixMapper(mapper, db.Prefix)
	}
	if "" != strings.TrimSpace(db.Suffix) {
		mapper = core.NewSuffixMapper(mapper, db.Suffix)
	}
	engine.SetTableMapper(mapper) // !必须配置名称转移

	// 测试数据库连接成功
	if *db.Ping {
		logger.Info("开始测试数据库连接", field.New("db", db))
		err = engine.Ping()
	}

	return
}

func enableSSH(db *config.DB, logger log.Logger) (err error) {
	if !db.SSHEnabled() {
		return
	}

	password := db.SSH.Password
	keyfile := db.SSH.Keyfile
	auth := gox.Ift("" != password, ssh.Password(password), sshtunnel.PrivateKeyFile(keyfile))
	host := fmt.Sprintf("%s@%s", db.SSH.Username, db.SSH.Addr)
	if tunnel, ne := sshtunnel.NewSSHTunnel(host, auth, db.Host, "65512"); nil != ne {
		err = ne
	} else {
		tunnel.Log = internal.NewSsh(logger)
		go startTunnel(tunnel)
		time.Sleep(100 * time.Millisecond)
		db.Host = fmt.Sprintf("127.0.0.1:%d", tunnel.Local.Port)
	}

	return
}

func startTunnel(tunnel *sshtunnel.SSHTunnel) {
	_ = tunnel.Start()
}
