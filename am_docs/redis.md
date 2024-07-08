# Redis 工具包
集成了Redis配置、set和get。
## 使用方法
> Attention⚠️：如果需要使用JWT套件，请务必配置Redis，因为该套件依赖Redis执行。
### 配置Redis
```go
// 使用默认配置 127.0.0.1:6379 password="" DB=0
am_redis.Setup(nil)

// 使用自定义配置
var redisConn = am_redis.RedisConn{
	Addr: "localhost",
	Port: "6379",
	Password: "",
	DB: 0,
}

am_redis.Setup(&redisConn)
```

### 存入Redis
```go
// 参数1: key
// 参数2: value
// 参数3: 过期时间，以s为计数单位。0表示永不过期
am_redis.SetValue("key", "value", 0)
```

### 从Redis取出
```go
// 参数1: key
am_redis.GetValue("key")
```

### 从Redis中删除
```go
// 参数1: key
am_redis.DelValue("key")
```