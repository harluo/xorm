package core

import (
	"runtime"

	"github.com/goexl/gox"
	"github.com/goexl/gox/field"
	"github.com/goexl/log"
	"xorm.io/xorm"
)

type Tx struct {
	engine *Engine
	logger log.Logger

	_ gox.Pointerized
}

func newTx(engine *Engine, logger log.Logger) *Tx {
	return &Tx{
		engine: engine,
		logger: logger,
	}
}

func (t *Tx) Do(fun Function, fields ...gox.Field[any]) (int64, error) {
	return t.do(func(tx *Session) (int64, error) {
		return fun(tx)
	}, fields...)
}

func (t *Tx) do(fun Function, fields ...gox.Field[any]) (affected int64, err error) {
	session := t.engine.NewSession()
	if err = t.begin(session, fields...); nil != err {
		return
	}
	defer t.close(session, fields...)

	if affected, err = fun(&Session{Session: session}); nil != err {
		t.rollback(session, fields...)
	} else {
		t.commit(session, fields...)
	}

	return
}

func (t *Tx) begin(tx *xorm.Session, fields ...gox.Field[any]) (err error) {
	if err = tx.Begin(); nil != err {
		t.error(err, "开始数据库事务出错", fields...)
	}

	return
}

func (t *Tx) commit(tx *xorm.Session, fields ...gox.Field[any]) {
	if err := tx.Commit(); nil != err {
		t.error(err, "提交数据库事务出错", fields...)
	}
}

func (t *Tx) close(tx *xorm.Session, fields ...gox.Field[any]) {
	if err := tx.Close(); nil != err {
		t.error(err, "关闭数据库事务出错", fields...)
	}
}

func (t *Tx) rollback(tx *xorm.Session, fields ...gox.Field[any]) {
	if err := tx.Rollback(); nil != err {
		t.error(err, "回退数据库事务出错", fields...)
	}
}

func (t *Tx) error(err error, msg string, fields ...gox.Field[any]) {
	fun, _, line, _ := runtime.Caller(1)

	logFields := make([]gox.Field[any], 0, len(fields)+3)
	logFields = append(logFields, field.New("line", line))
	logFields = append(logFields, fields...)
	logFields = append(logFields, field.Error(err))
	t.logger.Error(msg, field.New("fun", runtime.FuncForPC(fun).Name()), logFields...)
}
