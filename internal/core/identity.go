package core

import (
	"github.com/goexl/id"
	"github.com/goexl/model"
	"github.com/harluo/di"
)

type Identity struct {
	model.Base `xorm:"extends"`
}

func (i *Identity) BeforeInsert() {
	if i.Id == 0 {
		di.New().Instance().Get(i.setId).Build().Apply()
	}
}

func (i *Identity) setId(generator id.Generator) (err error) {
	if generated, ne := generator.Next(); nil != ne {
		err = ne
	} else {
		i.Id = generated.Get()
	}

	return
}
