package utils

import (
	"strings"
	"time"

	"github.com/pkg/errors"
	redis "gopkg.in/redis.v5"
)

type RedisV5Client struct {
	Options *redis.Options
	Client  *redis.Client
}

func (s *RedisV5Client) Conn() (err error) {
	s.Client = redis.NewClient(s.Options)
	if _, err = s.Client.Ping().Result(); err != nil {
		err = errors.Wrap(err, "RedisV5Client Conn")
	}
	return
}

func (s *RedisV5Client) Get(key string) (value string, err error) {
	if "" == key {
		err = errors.Wrap(errors.New("key empty"), "RedisV5Client Get")
		return
	}

	value, err = s.Client.Get(key).Result()
	//fmt.Printf("v:[%v], e:[%v]\n", value, err)
	if err != nil {
		//if err == redis.Nil {
		//	fmt.Println("nil nil nil nil")
		//}
		if err == redis.Nil {
			value = ""
			err = nil
		} else {
			err = errors.Wrap(err, "RedisV5Client Get")
		}
	}
	return
}

func (s *RedisV5Client) Set(key string, value interface{}, expiration time.Duration) (err error) {
	if err = s.Client.Set(key, value, expiration).Err(); err != nil {
		err = errors.Wrap(err, "RedisV5Client Set")
	}
	return
}

func (s *RedisV5Client) Del(key string) (err error) {
	if strings.TrimSpace(key) == "" {
		err = errors.Wrap(errors.New("param `key` empty"), "RedisV5Client Del")
	}
	if err = s.Client.Del(key).Err(); err != nil {
		err = errors.Wrap(err, "RedisV5Client Del")
	}
	return
}

func (s *RedisV5Client) SAdd(key string, members ...interface{}) (err error) {
	if strings.TrimSpace(key) == "" {
		err = errors.Wrap(errors.New("param `key` empty"), "RedisV5Client SAdd")
	}
	if err = s.Client.SAdd(key, members...).Err(); err != nil {
		err = errors.Wrap(err, "RedisV5Client SAdd")
	}
	return
}

func (s *RedisV5Client) SIsMember(key string, member interface{}) (boolCmd *redis.BoolCmd) {
	/*
		if strings.TrimSpace(key) == "" {
			err = errors.Wrap(errors.New("param `key` empty"), "RedisV5Client SIsMember")
		}
	*/
	return s.Client.SIsMember(key, member)
}
