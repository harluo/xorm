package core

import (
	"github.com/harluo/di"
)

func init() {
	di.New().Instance().Put(
		newEngines,
		newEngine,
		newTx,
		newDB,
	).Build().Apply()
}
