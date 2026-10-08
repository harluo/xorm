package config

import (
	"runtime"
	"time"
)

const (
	// 每个处理器默认持有的最大连接数
	connectionsPerCPU = 4
	// 连接默认的最大存活时间
	defaultLifetime = 15 * time.Minute
)

type Connection struct {
	// 最大打开连接数
	Open int `yaml:"open" json:"open" xml:"open" toml:"open"`
	// 最大休眠连接数
	Idle int `yaml:"idle" json:"idle" xml:"idle" toml:"idle"`
	// 每个连接最大存活时间
	Lifetime time.Duration `yaml:"lifetime" json:"lifetime" xml:"lifetime" toml:"lifetime"`
}

// Default 设置连接池的默认参数，未显式配置的参数在启动时根据处理器数量动态计算
func (c *Connection) Default() {
	if c.Open <= 0 { // 未配置最大打开连接数时，根据处理器数量动态计算
		c.Open = connectionsPerCPU * runtime.GOMAXPROCS(0)
	}
	if c.Idle <= 0 { // 未配置最大休眠连接数时，与最大打开连接数保持一致，避免频繁重建连接
		c.Idle = c.Open
	}
	if c.Idle > c.Open { // 休眠连接数不能超过最大打开连接数
		c.Idle = c.Open
	}
	if c.Lifetime <= 0 { // 未配置连接最大存活时间时，使用默认值
		c.Lifetime = defaultLifetime
	}
}
