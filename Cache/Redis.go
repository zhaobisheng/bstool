package Cache

import (
	"bsTool/Config"

	"github.com/garyburd/redigo/redis"
)

var pool *redis.Pool


func init() {
	Config.InitConfig("app.conf")
	host := Config.ReadKey("redis", "host")
	port := Config.ReadKey("redis", "port")
	pool = &redis.Pool{
		MaxIdle:     16,
		MaxActive:   0,
		IdleTimeout: 300,
		Dial: func() (redis.Conn, error) {
			return redis.Dial("tcp", host+":"+port)
		},
	}
}

func Get(key string) (string, error) {
	c := pool.Get()
	defer c.Close()
	rs, err := redis.String(c.Do("GET", key))
	if err != nil {
		return "", err
	}
	return rs, nil
}

func Set(key, value string) error {
	c := pool.Get()
	defer c.Close()
	_, err := c.Do("set", key, value)
	if err != nil {
		return err
	}
	return nil
}

func Exists(key string) (int, error) {
	c := pool.Get()
	defer c.Close()
	r, err := redis.Int(c.Do("EXISTS", key))
	if err != nil {
		return r, err
	}
	return r, nil
}
