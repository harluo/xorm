package xorm

import (
	"time"
)

// CreatedAt 创建时间
type CreatedAt struct {
	Identity `xorm:"extends"`

	// 创建时间
	Created time.Time `xorm:"created notnull default(CURRENT_TIMESTAMP) comment('创建时间，创建时自动设置')" json:"created,omitempty"` // nolint:lll
}
