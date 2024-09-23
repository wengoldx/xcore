// Copyright (c) 2018-Now Dunyu All Rights Reserved.
//
// Author      : https://www.wengold.net
// Email       : support@wengold.net
//
// Prismy.No | Date       | Modified by. | Description
// -------------------------------------------------------------------
// 00001       2021/08/11   yangping       New version
// -------------------------------------------------------------------

package mvc

import (
	"github.com/gomodule/redigo/redis"
	"github.com/wengoldx/xcore/logger"
)

// HSet set the string value of a hash field.
func (c *WingRedisConn) HSet(key string, field, value any) bool {
	con := c.redisPool.Get()
	defer con.Close()

	set, err := redis.Bool(con.Do("HSET", c.serviceNamespace+key, field, value))
	if err != nil {
		logger.E("Redis:HSET [key"+key+"] err:", err)
		return false
	}
	return set
}

// HGet get the string value of a hash field.
func (c *WingRedisConn) HGet(key string, field any) (string, error) {
	con := c.redisPool.Get()
	defer con.Close()

	return redis.String(con.Do("HGET", c.serviceNamespace+key, field))
}

// HGetInt get the int value of a hash field.
func (c *WingRedisConn) HGetInt(key string, field any) (int, error) {
	con := c.redisPool.Get()
	defer con.Close()

	return redis.Int(con.Do("HGET", c.serviceNamespace+key, field))
}

// HGetInt64 get the int64 value of a hash field.
func (c *WingRedisConn) HGetInt64(key string, field any) (int64, error) {
	con := c.redisPool.Get()
	defer con.Close()

	return redis.Int64(con.Do("HGET", c.serviceNamespace+key, field))
}

// HGetUint64 get the uint64 value of a hash field.
func (c *WingRedisConn) HGetUint64(key string, field any) (uint64, error) {
	con := c.redisPool.Get()
	defer con.Close()

	return redis.Uint64(con.Do("HGET", c.serviceNamespace+key, field))
}

// HGetFloat64 get the float value of a hash field.
func (c *WingRedisConn) HGetFloat64(key string, field any) (float64, error) {
	con := c.redisPool.Get()
	defer con.Close()

	return redis.Float64(con.Do("HGET", c.serviceNamespace+key, field))
}

// HGetBytes get the bytes array of a hash field.
func (c *WingRedisConn) HGetBytes(key string, field any) ([]byte, error) {
	con := c.redisPool.Get()
	defer con.Close()

	return redis.Bytes(con.Do("HGET", c.serviceNamespace+key, field))
}

// HGetBool get the bool value of a hash field.
func (c *WingRedisConn) HGetBool(key string, field any) (bool, error) {
	con := c.redisPool.Get()
	defer con.Close()

	return redis.Bool(con.Do("HGET", c.serviceNamespace+key, field))
}
