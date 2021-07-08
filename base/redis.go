package base

import (
	"github.com/garyburd/redigo/redis"
	"github.com/mszhangyi/infra"
	"github.com/sirupsen/logrus"
	"time"
)

var pool *redis.Pool

func RedisPool() *redis.Pool {
	Check(pool)
	return pool
}

//redis starter，并且设置为全局
type RedisPoolStarter struct {
	infra.BaseStarter
}

func (r *RedisPoolStarter) Setup() {
	pool = &redis.Pool{
		MaxIdle:     props.RedisMaxIdle,                                  //空闲连接数
		MaxActive:   props.RedisMaxActive,                                //活跃连接数
		IdleTimeout: time.Duration(props.RedisIdleTimeout) * time.Second, //连接空闲关闭时间
		Dial: func() (redis.Conn, error) {
			//redis.DialPassword(props.RedisPwd)
			//此处对应redis ip及端口号
			conn, err := redis.Dial("tcp", props.RedisAddr)
			if err != nil {
				return nil, err
			}
			//此处1234对应redis密码
			if props.RedisPwd != "" {
				if _, err := conn.Do("AUTH", props.RedisPwd); err != nil {
					conn.Close()
					return nil, err
				}
			}
			return conn, err
		},
		TestOnBorrow: func(conn redis.Conn, t time.Time) error {
			if time.Since(t) < time.Minute {
				return nil
			}
			_, err := conn.Do("PING")
			return err
		},
	}

	c := pool.Get()
	_, err := c.Do("ping")
	if err != nil {
		logrus.Panic("redis：", err)
		panic(err)
	}
}
