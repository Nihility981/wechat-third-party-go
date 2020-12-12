package cache

import (
	"qlong"
	"qlong/cache"
)

var Tools = &Cache{}

// Cache 缓存类
type Cache struct {
	rc *cache.RedisClass
}

// Init 实始化
func (me *Cache) Init(host, db int) error {
	// 连接redis
	// redis主机为基础数据，redis库为基础数据库
	rc, rcErr := qlong.GetRedis(host, db)
	if rcErr != nil {
		return rcErr
	}
	me.rc = rc
	return nil
}
