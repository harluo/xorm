package xorm

import (
	"github.com/goexl/id"
	"github.com/harluo/di"
)

type Identity struct {
	Id uint64 `xorm:"pk BIGINT notnull default(0) comment('编号，用来唯一标识数据')" json:"id,string,omitempty"`
}

func (i *Identity) BeforeInsert() {
	if i.Id == 0 {
		di.New().Instance().Get(i.set).Build().Apply()
	}
}

func (i *Identity) set(generator id.Generator) (err error) {
	if generated, ne := generator.Next(); nil != ne {
		err = ne
	} else {
		i.Id = generated.Get()
	}

	return
}
