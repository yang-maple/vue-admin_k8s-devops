package dao

import (
	"context"
	"errors"
	"github.com/go-redis/redis/v8"
	"kubeops/db"
	"time"
)

var RdbValue rdbValue

type rdbValue struct{}

// SetValue 设置缓存
func (rdb *rdbValue) SetValue(key, value string, expiration time.Duration) error {
	err := db.Rdb.Set(context.Background(), key, value, expiration*time.Second).Err()
	if err != nil {
		return err
	}
	return nil
}

// GetValue 读取缓存
func (rdb *rdbValue) GetValue(key string) (value string, err error) {
	value, err = db.Rdb.Get(context.Background(), key).Result()
	if errors.Is(err, redis.Nil) {
		return "", errors.New("key does not exist")
	} else if err != nil {
		return "", err
	}
	return value, nil
}

// DelValue 删除缓存
func (rdb *rdbValue) DelValue(key string) error {
	err := db.Rdb.Del(context.Background(), key).Err()
	if err != nil {
		return errors.New("redis del value error")
	}
	return nil

}
