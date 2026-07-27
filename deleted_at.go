package xorm

import (
	"time"
)

type DeletedAt struct {
	Model `xorm:"extends"`

	Deleted time.Time `xorm:"deleted comment('最后删除时间，用于实现软删除功能')" json:"deleted,omitempty"`
}
