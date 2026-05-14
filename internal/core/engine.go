package core

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/elliotchance/sshtunnel"
	"github.com/goexl/exception"
	"github.com/goexl/gox"
	"github.com/goexl/gox/field"
	"github.com/goexl/log"
	"github.com/harluo/xorm/internal/config"
	"github.com/harluo/xorm/internal/core/internal/constant"
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
	return fmt.Sprintf("`%s`", e.GetColumnMapper().Obj2Table(field))
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

func (e *Engine) Table(name any) *Session {
	return &Session{
		shadowSession: e.shadowEngine.Table(name),
	}
}

func (e *Engine) Join(name any, condition any, args ...any) *Session {
	return &Session{
		shadowSession: e.shadowEngine.Join("INNER", name, condition, args...),
	}
}

func (e *Engine) Left(name any, condition any, args ...any) *Session {
	return &Session{
		shadowSession: e.shadowEngine.Join("LEFT", name, condition, args...),
	}
}

func (e *Engine) Right(name any, condition any, args ...any) *Session {
	return &Session{
		shadowSession: e.shadowEngine.Join("RIGHT", name, condition, args...),
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

func (e *Engine) Limit(limit int, starts ...int) *Session {
	return &Session{
		shadowSession: e.shadowEngine.Limit(limit, starts...),
	}
}

func (e *Engine) Id(id any) *Session {
	return &Session{
		shadowSession: e.shadowEngine.ID(id),
	}
}

func newEngines(db *config.DB, logger log.Logger) (selector *Selector, err error) {
	engines := make(map[string]*Engine, len(db.Databases)+1)
	if err = initDefault(&db.Database, engines, logger); err != nil {
		return
	}

	for name, database := range db.Databases {
		if engine, nee := createEngine(database, logger); nee != nil {
			err = nee
		} else {
			engines[name] = engine
		}
		if err != nil {
			break
		}
	}
	if err == nil {
		selector = newSelector(engines)
	}

	return
}

func newEngine(selector *Selector) (engine *Engine, err error) {
	if saved, exists := selector.Exists(constant.NameDefault); !exists {
		err = exception.New().Build()
	} else {
		engine = saved
	}

	return
}

func initDefault(database *config.Database, engines map[string]*Engine, logger log.Logger) (err error) {
	if database.Host == "" {
		return
	}

	if engine, cee := createEngine(database, logger); cee != nil {
		err = cee
	} else {
		engines[constant.NameDefault] = engine
	}

	return
}

func createEngine(database *config.Database, logger log.Logger) (engine *Engine, err error) {
	engine = new(Engine)
	if ese := enableSSH(database, logger); nil != ese {
		err = ese
	} else if master, cme := createMaster(database); nil != cme {
		err = cme
	} else if slaves, cse := createSlaves(database); nil != cse {
		err = cse
	} else if engine.shadowEngine, err = xorm.NewEngineGroup(master, slaves); nil == err {
		err = setupEngine(database, engine, logger)
	}

	return
}

func createMaster(database *config.Database) (engine *xorm.Engine, err error) {
	if dsn, de := database.SN(&database.Server); nil != de {
		err = de
	} else {
		engine, err = xorm.NewEngine(database.Type.String(), dsn)
	}

	return
}

func createSlaves(database *config.Database) (slaves []*xorm.Engine, err error) {
	slaves = make([]*xorm.Engine, len(database.Slaves))
	for _, slave := range database.Slaves {
		if dsn, de := database.SN(slave); nil != de {
			err = de
		} else if engine, nee := xorm.NewEngine(database.Type.String(), dsn); nil != nee {
			err = nee
		} else {
			slaves = append(slaves, engine)
		}
	}

	return
}

func setupEngine(db *config.Database, engine *Engine, logger log.Logger) (err error) {
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
	if db.Ping {
		logger.Info("开始测试数据库连接", field.New("db", db))
		err = engine.Ping()
	}

	return
}

func enableSSH(database *config.Database, logger log.Logger) (err error) {
	if !database.SSHEnabled() {
		return
	}

	password := database.SSH.Password
	keyfile := database.SSH.Keyfile
	auth := gox.Ift("" != password, ssh.Password(password), sshtunnel.PrivateKeyFile(keyfile))
	host := fmt.Sprintf("%s@%s", database.SSH.Username, database.SSH.Addr)
	if tunnel, ne := sshtunnel.NewSSHTunnel(host, auth, database.Host, "65512"); nil != ne {
		err = ne
	} else {
		tunnel.Log = internal.NewSsh(logger)
		go startTunnel(tunnel)
		time.Sleep(100 * time.Millisecond)
		database.Host = fmt.Sprintf("127.0.0.1:%d", tunnel.Local.Port)
	}

	return
}

func startTunnel(tunnel *sshtunnel.SSHTunnel) {
	_ = tunnel.Start()
}
