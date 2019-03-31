package cache

import (
	"errors"
	"github.com/astaxie/beego/cache"
	_ "github.com/astaxie/beego/cache/redis"
	"github.com/astaxie/beego/logs"
	"github.com/qiangjiyi/casbin_auth/common/util"
)

var (
	DefaultConn     = "127.0.0.1:6379" // redis conn addr info
	DefaultDbNum    = "0"              // default conn redis DB number
	DefaultPassword = ""               // conn redis password
)

// RedisCache create a redis cache adapter
func RedisCache(key string, s ...string) (cache.Cache, error) {
	if key == "" {
		return nil, errors.New("Cache key can't be empty")
	}

	var config string
	if len(s) == 0 {
		config = `{"key":"` + key + `", "conn":"` + DefaultConn + `", "dbNum":"` + DefaultDbNum + `", "password":"` + DefaultPassword + `"}`
	} else {
		conn := DefaultConn
		dbNum := DefaultDbNum
		password := DefaultPassword
		for i, val := range s {
			switch i {
			case 0:
				conn = val
			case 1:
				dbNum = val
			case 2:
				password = val
			default:
				return nil, errors.New("Cache params are not correct. Should be order by (key, conn, dbNum, password)")
			}
		}
		config = `{"key":"` + key + `", "conn":"` + conn + `", "dbNum":"` + dbNum + `", "password":"` + password + `"}`
	}

	return cache.NewCache("redis", config)
}

// RedisCacheForSession
func RedisCacheForSession() cache.Cache {
	redis, err := RedisCache(util.TokenCache)
	if err != nil {
		logs.Error("Init redis cache error. ---> %+v", err)
		return nil
	}
	return redis
}
