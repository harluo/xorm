package core

import (
	"github.com/harluo/di"
)

func init() {
	di.New().Instance().Put(
		newEngine,
		newTx,
		newSynchronizer,
		newDB,
	).Build().Apply()
}
