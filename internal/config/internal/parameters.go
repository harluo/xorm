package internal

import (
	"fmt"
	"strings"

	"github.com/goexl/db"
	"github.com/goexl/gox"
	"github.com/harluo/xorm/internal/config/internal/constant"
)

type Parameters map[string]any

func (p Parameters) Merge(dsn string, typ db.Type) (final string) {
	values := make([]string, len(p))
	for key, value := range p {
		builder := gox.StringBuilder(key)
		if stringed := gox.ToString(value); "" != stringed {
			builder.Append(constant.Equal).Append(stringed)
		}
		values = append(values, builder.String())
	}

	switch typ {
	case db.TypePostgres:
		final = fmt.Sprintf("%s %s", dsn, strings.Join(values, constant.Space))
	case db.TypeMySQL:
		final = fmt.Sprintf("%s?%s", dsn, strings.Join(values, constant.And))
	default:
		final = fmt.Sprintf("%s?%s", dsn, strings.Join(values, constant.And))
	}

	return
}
