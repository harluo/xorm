package core

import (
	"github.com/harluo/xorm/internal/core/internal/constant"
)

type Selector struct {
	engines map[string]*Engine
}

func newSelector(engines map[string]*Engine) *Selector {
	return &Selector{
		engines: engines,
	}
}

func (s *Selector) Default() *Engine {
	return s.engines[constant.NameDefault]
}

func (s *Selector) Exists(name string) (engine *Engine, exists bool) {
	engine, exists =s.engines[name]

	return
}

func (s *Selector) Select(name string) *Engine {
	return s.engines[name]
}
