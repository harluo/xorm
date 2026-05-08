package xorm

import (
	"time"
)

// UpdatedAt 更新时间
type UpdatedAt struct {
	Identity `xorm:"extends"`

	// 最后更新时间
	Updated time.Time `xorm:"updated notnull default(CURRENT_TIMESTAMP) comment('最后更新时间，每次更新时自动修改')" json:"updated,omitempty"` // nolint:lll
}
