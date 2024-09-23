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
	"errors"
	"fmt"
	"time"

	"github.com/astaxie/beego"
	"github.com/gomodule/redigo/redis"
	"github.com/wengoldx/xcore/invar"
	"github.com/wengoldx/xcore/logger"
)

// WingRedisConn content provider to support redis utils
type WingRedisConn struct {
	// redis connections pool, use redisPool.Get() to return one connect client
	// then call redisPool.Close() to close it after finished use.
	redisPool *redis.Pool

	// redis server host address, usually use ip.
	serverHost string

	// redis server auth password.
	serverAuthPwd string

	// serviceNamespace service namespace to distinguish others server, it will
	// append with key as prefix to get or set data to redis database, default is empty.
	//
	// `NOTICE` :
	//
	// you should config namespace in /conf/app.config file, or call WingRedis.SetNamespace(ns)
	// to set unique server namespace when multiple services connecting the one redis server.
	serviceNamespace string

	// deadlock max duration, default 20 seconds
	deadlockDuration int64
}

const (
	redisConfigHost = "%s::host"      // configs key of redis host and port
	redisConfigPwd  = "%s::pwd"       // configs key of redis password
	redisConfigNs   = "%s::namespace" // configs key of redis namespace
	redisConfigLock = "%s::deadlock"  // configs key of redis lock max duration
)

// The follow options may support by diffrent Redis version, get more info
// by link https://redis.io/commands webset.
const (
	OptEX      = "EX"      // seconds -- Set the specified expire time, in seconds
	OptPX      = "PX"      // milliseconds -- Set the specified expire time, in milliseconds.
	OptEXAT    = "EXAT"    // timestamp-seconds -- Set the specified Unix time at which the key will expire, in seconds.
	OptPXAT    = "PXAT"    // timestamp-milliseconds -- Set the specified Unix time at which the key will expire, in milliseconds.
	OptNX      = "NX"      // Only set the key if it does not already exist.
	OptXX      = "XX"      // Only set the key if it already exist.
	OptKEEPTTL = "KEEPTTL" // Retain the time to live associated with the key, use for SET commond.
	ExpNX      = "NX"      // Set expiry only when the key has no expiry.
	ExpXX      = "XX"      // Set expiry only when the key has an existing expiry.
	ExpGT      = "GT"      // Set expiry only when the new expiry is greater than current one.
	ExpLT      = "LT"      // Set expiry only when the new expiry is less than current one.

	CusOptDel = "DELETE" // The custom option to delete redis key after get commond execute.
)

// WingRedis a connecter to access redis database data,
// it will nil before mvc.OpenRedis() called
var WingRedis *WingRedisConn

// readRedisCofnigs read redis params from config file, than verify them if empty.
func readRedisCofnigs(session string) (string, string, string, int64, error) {
	host := beego.AppConfig.String(fmt.Sprintf(redisConfigHost, session))
	pwd := beego.AppConfig.String(fmt.Sprintf(redisConfigPwd, session))
	ns := beego.AppConfig.String(fmt.Sprintf(redisConfigNs, session)) // allow empty
	lock := beego.AppConfig.DefaultInt64(fmt.Sprintf(redisConfigLock, session), 20)

	if host == "" || pwd == "" {
		return "", "", "", 0, invar.ErrInvalidConfigs
	}

	if lock <= 0 {
		lock = 20 // default 20 seconds
	}
	return host, pwd, ns, lock, nil
}

// OpenRedis connect redis database server and auth password,
// the connections holded by mvc.WingRedis object.
//
// `NOTICE` :
//
// you must config redis params in /conf/app.config file as:
//
// ---
//
// #### Case 1 : For connect on prod mode.
//
//	[redis]
//	host = "127.0.0.1:6379"
//	pwd = "123456"
//	namespace = "project_namespace"
//	deadlock = 20
//
// #### Case 2 : For connect on dev mode.
//
//	[redis-dev]
//	host = "127.0.0.1:6379"
//	pwd = "123456"
//	namespace = "project_namespace"
//	deadlock = 20
//
// #### Case 3 : For both dev and prod mode, you can config all of up cases.
//
// #### Case 4 : For connect on prod mode without namespace and use default lock time.
//
//	[redis]
//	host = "127.0.0.1:6379"
//	pwd = "123456"
//
// ---
//
// The configs means as:
//
//	`host` - is the redis server host ip and port.
//	`pwd`  - is the redis server authenticate password.
//	`namespace` - is the prefix string or store key.
//	`deadlock`  - is the max time of deadlock, in seconds.
func OpenRedis() error {
	session := "redis"
	if beego.BConfig.RunMode == "dev" {
		session = session + "-dev"
	}

	host, pwd, ns, lock, err := readRedisCofnigs(session)
	if err != nil {
		return err
	}

	WingRedis = &WingRedisConn{nil, host, pwd, ns, lock}
	WingRedis.redisPool = &redis.Pool{
		MaxIdle: 16, MaxActive: 0, IdleTimeout: 300 * time.Second,
		Dial: func() (redis.Conn, error) {
			host, pwd := WingRedis.serverHost, WingRedis.serverAuthPwd
			c, err := redis.Dial("tcp", host) // dial TCP connection
			if err != nil {
				return nil, err
			}

			// authenticate connection password. see https://redis.io/commands/auth
			if _, err := c.Do("AUTH", pwd); err != nil {
				c.Close()
				panic(err)
			}
			return c, nil
		},
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			if _, err := c.Do("PING"); err != nil {
				return errors.New("Ping redis err: " + err.Error())
			}
			return nil
		},
	}
	return nil
}

// SetNamespace set server uniqu namespace
func (c *WingRedisConn) SetNamespace(ns string) {
	if ns != "" {
		c.serviceNamespace = ns
	}
}

// SetDeadlock set the max deadlock duration
func (c *WingRedisConn) SetDeadlock(dur int64) {
	if dur > 0 {
		c.deadlockDuration = dur
	}
}

// GetDeadlock get the max deadlock duration
func (c *WingRedisConn) GetDeadlock() int64 {
	return c.deadlockDuration
}

// NsKey transform origin key to namespaced key
func (c *WingRedisConn) NsKey(key string) string {
	return c.serviceNamespace + key
}

// NsKeys transform origin keys to namespaced keys
func (c *WingRedisConn) NsKeys(keys ...string) []string {
	return c.NsArrKeys(keys)
}

// NsArrKeys transform origin keys to namespaced keys
func (c *WingRedisConn) NsArrKeys(keys []string) []string {
	nskeys := []string{}
	for _, key := range keys {
		nskeys = append(nskeys, c.NsKey(key))
	}
	return nskeys
}

// ServerTime get redis server unix time is seconds and microsecodns.
//
// see command [time](https://redis.io/commands/time)
func (c *WingRedisConn) ServerTime() (int64, int64) {
	con := c.redisPool.Get()
	defer con.Close()

	st, err := redis.Int64s(con.Do("TIME"))
	if err != nil || len(st) != 2 {
		logger.E("Redis:TIME Failed set redis server time")
		return 0, 0
	}
	return st[0], st[1]
}

// setWithExpire set a value and expiration of a key.
//
// see commands [setex](https://redis.io/commands/setex),
// [psetex](https://redis.io/commands/psetex).
func (c *WingRedisConn) setWithExpire(key, commond string, value any, expire int64) error {
	con := c.redisPool.Get()
	defer con.Close()

	_, err := con.Do(commond, c.serviceNamespace+key, expire, value)
	return err
}

// parseGetOptions parse the GETEX, GETDEL commonds options
func (c *WingRedisConn) parseGetOptions(options ...any) (string, int64) {
	if len(options) == 0 {
		logger.E("Redis: invalid options, parse failed")
		return "", 0
	}

	switch option := options[0].(type) {
	case string:
		if len(options) > 1 /* parse expire */ {
			switch expire := options[1].(type) {
			case int64:
				return option, expire
			default:
				logger.E("Redis: expire value type is not int64")
			}
		} else {
			return option, 0 // just for CusOptDel
		}
	default:
		logger.E("Redis: option value type is not string")
	}
	return "", 0
}

// Get get a value of key, than set value expire time or delete by given options.
//
// see commands [get](https://redis.io/commands/get),
// [getex](https://redis.io/commands/getex),
// [getdel](https://redis.io/commands/getdel)
func (c *WingRedisConn) getWithOptions(key string, options ...any) (any, error) {
	con := c.redisPool.Get()
	defer con.Close()

	if options != nil {
		var reply any
		err := invar.ErrInvalidRedisOptions.Err()
		if option, expire := c.parseGetOptions(options); option != "" {
			switch option {
			case CusOptDel:
				reply, err = con.Do("GETDEL", c.serviceNamespace+key)
			case OptEX, OptPX, OptEXAT, OptPXAT:
				reply, err = con.Do("GETEX", c.serviceNamespace+key, option, expire)
			}
		}
		return redis.String(reply, err)
	}
	return con.Do("GET", c.serviceNamespace+key)
}

// getKeyExpire get the time to live for a key by given commond, it may return unexist
// key error if the key unexist or set expiration and now expire, or return keeplive
// error if the exist key has no associated expire.
//
// see commands [ttl](https://redis.io/commands/ttl),
// [pttl](https://redis.io/commands/pttl)
func (c *WingRedisConn) getKeyExpire(key, commond string) (int64, error) {
	con := c.redisPool.Get()
	defer con.Close()

	expire, err := redis.Int64(con.Do(commond, c.serviceNamespace+key))
	if expire == -2 {
		return 0, invar.ErrUnexistRedisKey
	} else if expire == -1 {
		return 0, invar.ErrNoAssociatedExpire
	}
	return expire, err
}

// setKeyExpire set the expiration for a key by given commond, the optional values
// can be set as ExpNX, ExpXX, ExpGT, ExpLT since Redis 7.0 support.
//
// see [expire](https://redis.io/commands/expire),
// [expireat](https://redis.io/commands/expireat)
func (c *WingRedisConn) setKeyExpire(key, commond string, expire int64, option ...string) bool {
	con := c.redisPool.Get()
	defer con.Close()

	set, err := false, invar.ErrInvalidRedisOptions.Err()
	if len(option) > 0 {
		switch option[0] {
		case ExpNX, ExpXX, ExpGT, ExpLT:
			set, err = redis.Bool(con.Do(commond, c.serviceNamespace+key, expire, option[0]))
		}
	} else {
		set, err = redis.Bool(con.Do(commond, c.serviceNamespace+key, expire))
	}

	if err != nil {
		logger.E("Redis:EXPIRE [key"+key+":"+commond+"] err:", err)
		return false
	}
	return set
}

// GetRedisPool get redis pool
func (c *WingRedisConn) GetRedisPool() redis.Conn {
	return c.redisPool.Get()
}
