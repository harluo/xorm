package core

import (
	"context"

	"github.com/goexl/gox"
)

// Session 事务
type Session struct {
	*shadowSession

	_ gox.Pointerized
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

func (s *Session) OrderBy(order any, args ...any) *Session {
	return &Session{
		shadowSession: s.shadowSession.OrderBy(order, args...),
	}
}
