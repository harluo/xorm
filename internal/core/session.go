package core

import (
	"context"
	"fmt"

	"github.com/goexl/gox"
)

// Session 事务
type Session struct {
	*shadowSession

	_ gox.Pointerized
}

func (s *Session) ColumnName(field string) string {
	return fmt.Sprintf("`%s`", s.Engine().GetColumnMapper().Obj2Table(field))
}

func (s *Session) Context(ctx context.Context) *Session {
	return &Session{
		shadowSession: s.shadowSession.Context(ctx),
	}
}

func (s *Session) Cols(columns ...string) *Session {
	return &Session{
		shadowSession: s.shadowSession.Cols(columns...),
	}
}

func (s *Session) Where(query any, args ...any) *Session {
	return &Session{
		shadowSession: s.shadowSession.Where(query, args...),
	}
}

func (s *Session) Table(name any) *Session {
	return &Session{
		shadowSession: s.shadowSession.Table(name),
	}
}

func (s *Session) Join(name any, condition any, args ...any) *Session {
	return &Session{
		shadowSession: s.shadowSession.Join("INNER", name, condition, args...),
	}
}

func (s *Session) LeftJoin(name any, condition any, args ...any) *Session {
	return &Session{
		shadowSession: s.shadowSession.Join("LEFT", name, condition, args...),
	}
}

func (s *Session) RightJoin(name any, condition any, args ...any) *Session {
	return &Session{
		shadowSession: s.shadowSession.Join("RIGHT", name, condition, args...),
	}
}

func (s *Session) OrderBy(order any, args ...any) *Session {
	return &Session{
		shadowSession: s.shadowSession.OrderBy(order, args...),
	}
}

func (s *Session) Desc(field string) *Session {
	return &Session{
		shadowSession: s.shadowSession.Desc(s.ColumnName(field)),
	}
}

func (s *Session) Asc(field string) *Session {
	return &Session{
		shadowSession: s.shadowSession.Asc(s.ColumnName(field)),
	}
}

func (s *Session) Id(id any) *Session {
	return &Session{
		shadowSession: s.shadowSession.ID(id),
	}
}

func (s *Session) Persist(object any) (affected int64, err error) {
	if ifa, ife := s.shadowSession.Insert(object); ife != nil {
		// !因为在在分表操作，添加不成功时，需要检查是不是表不存在
		// !存在这么一种可能性，就是程序一直长期未重启过，表结构就不会被同步也不会创建表
		affected, err = s.retryInsert(object, ife)
	} else {
		affected = ifa
	}

	return
}

func (s *Session) retryInsert(quota any, original error) (affected int64, err error) {
	if exists, tee := s.shadowSession.IsTableExist(quota); nil != tee {
		err = tee
	} else if !exists {
		err = s.shadowSession.Sync(quota)
	} else {
		err = original
	}

	if nil == err {
		affected, err = s.shadowSession.Insert(quota)
	}

	return
}
