package core

import (
	"github.com/harluo/di"
)

func init() {
	di.New().Instance().Put(
		newEngine,
		newTransaction,
		newSynchronizer,
	).Build().Apply()
}
