package am_redis

import (
	"errors"
	"github.com/go-redis/redis"
	"net"
	"strconv"
	"time"
)

// 工具错误类型
var (
	RedisGetNilError = errors.New("redis get nil")        // 未找到值
	TxFailedError    = errors.New("transaction failed")   // 握手失败
	TimeoutError     = errors.New("timeout")              // 连接超时
	NilPointError    = errors.New("no redis connections") // 没有设定redis连接
)

// redis连接结构体
type RedisConn struct {
	Addr     string `json:"addr"`
	Port     int    `json:"port"`
	Password string `json:"password"`
	DB       int    `json:"db"`
}

// redis连接参数默认值
var redisConn *RedisConn = &RedisConn{
	Addr:     "127.0.0.1",
	Port:     6379,
	Password: "",
	DB:       0,
}

// redis连接client
var client *redis.Client

// 初始化redis连接
func Setup(conn *RedisConn) error {
	connectionToRedis := redisConn // redis连接参数默认值
	if conn != nil {               // 修改redis连接参数
		connectionToRedis = conn
	}
	client = redis.NewClient(&redis.Options{ // 获取redis连接
		Addr:     connectionToRedis.Addr + ":" + strconv.Itoa(connectionToRedis.Port),
		Password: connectionToRedis.Password,
		DB:       connectionToRedis.DB,
	})
	_, err := client.Ping().Result() // 测试连通情况
	if err != nil {                  // 如果不通，抛出err
		return err
	}
	return nil // 通，返回client连接
}

/*
* 向redis中添加资料
* 参数：
* 1. key string 键
* 2. value string 值
* 3. exp int 过期时间，秒为单位，0永不过期
 */
func SetValue(key string, value string, exp int) error {
	if client == nil { // 如果没有redis连接，则返回空指针异常
		return NilPointError
	}
	client.Set(key, value, time.Duration(exp)*time.Second)
	return nil
}

/*
* 从redis中获取资料
* 1. key string 键
 */
func GetValue(key string) (string, error) {
	if client == nil { // 如果没有redis连接，则返回空指针异常
		return "", NilPointError
	}
	result, err := client.Get(key).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) { // 未找到值
			return "", RedisGetNilError
		} else if errors.Is(err, redis.TxFailedErr) { // 握手失败
			return "", TxFailedError
		} else if netErr, ok := err.(net.Error); ok && netErr.Timeout() { // 连接超时
			return "", TimeoutError
		} else { // 其他错误
			return "", err
		}
	}
	return result, nil
}

/*
* 删除资料
* 参数：
* 1. key：键
 */
func DelValue(key string) error {
	if client == nil { // 如果没有redis连接，则返回空指针异常
		return NilPointError
	}

	client.Del(key)
	return nil
}
