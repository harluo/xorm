package xorm

import (
	"time"
)

type Deleted struct {
	Deleted time.Time `xorm:"deleted" json:"deleted,omitempty"`
}
