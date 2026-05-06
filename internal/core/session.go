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

func (s *Session) Left(name any, condition any, args ...any) *Session {
	return &Session{
		shadowSession: s.shadowSession.Join("LEFT", name, condition, args...),
	}
}

func (s *Session) Right(name any, condition any, args ...any) *Session {
	return &Session{
		shadowSession: s.shadowSession.Join("RIGHT", name, condition, args...),
	}
}

func (s *Session) OrderBy(order any, args ...any) *Session {
	return &Session{
		shadowSession: s.shadowSession.OrderBy(order, args...),
	}
}

func (s *Session) Desc(field gox.Column) *Session {
	return &Session{
		shadowSession: s.shadowSession.Desc(s.ColumnName(field.String())),
	}
}

func (s *Session) Asc(field gox.Column) *Session {
	return &Session{
		shadowSession: s.shadowSession.Asc(s.ColumnName(field.String())),
	}
}

func (s *Session) Limit(limit int, starts ...int) *Session {
	return &Session{
		shadowSession: s.shadowSession.Limit(limit, starts...),
	}
}

func (s *Session) Id(id any) *Session {
	return &Session{
		shadowSession: s.shadowSession.ID(id),
	}
}

func (s *Session) Reget(object any) (exists bool, err error) {
	if ca, ce := s.shadowSession.Get(object); ce != nil {
		exists, err = s.retryExists(object, ce, func(i *bool, e *error) {
			*i, *e = s.shadowSession.Get(object)
		})
	} else {
		exists = ca
	}

	return
}

func (s *Session) Recount(table any, object ...any) (affected int64, err error) {
	if ca, ce := s.shadowSession.Count(object...); ce != nil {
		affected, err = s.retryAffected(table, ce, func(i *int64, e *error) {
			*i, *e = s.shadowSession.Count(object...)
		})
	} else {
		affected = ca
	}

	return
}

func (s *Session) Refind(table any, objects any, conditions ...any) (err error) {
	if ce := s.shadowSession.Find(objects, conditions...); ce != nil {
		_, err = s.retryAffected(table, ce, func(i *int64, e *error) {
			*e = s.shadowSession.Find(objects, conditions...)
		})
	}

	return
}

func (s *Session) Persist(object any) (affected int64, err error) {
	if ifa, ife := s.shadowSession.Insert(object); ife != nil {
		affected, err = s.retryAffected(object, ife, func(i *int64, e *error) {
			*i, *e = s.shadowSession.Insert(object)
		})
	} else {
		affected = ifa
	}

	return
}

func (s *Session) retryExists(object any, original error, callback func(*bool, *error)) (exists bool, err error) {
	if exists, tee := s.shadowSession.IsTableExist(object); nil != tee {
		err = tee
	} else if !exists {
		err = s.shadowSession.Sync(object)
	} else {
		err = original
	}

	if nil == err {
		callback(&exists, &err)
	}

	return
}

func (s *Session) retryAffected(object any, original error, callback func(*int64, *error)) (affected int64, err error) {
	if exists, tee := s.shadowSession.IsTableExist(object); nil != tee {
		err = tee
	} else if !exists {
		err = s.shadowSession.Sync(object)
	} else {
		err = original
	}

	if nil == err {
		callback(&affected, &err)
	}

	return
}
