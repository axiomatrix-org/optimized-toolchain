# 速率限制中间件
为了防范恶意flood攻击，需要对同一IP对同一controller的访问频率做出限制。

## 使用方法
```go
// 创建速率限制配置
// 参数1:窗口期内可访问次数
// 参数2:窗口期时间，以s为计数单位
// 下述配置为1秒种内可访问同一controller 5次
var defaultLimitConfig = am_ratelimit.NewRateLimitConfig(5, 1)

gin.SetMode(gin.DebugMode)
r := gin.Default()
// 设置中间件
r.POST("/test", defaultLimitConfig.RateLimitMiddleware, func(context *gin.Context) {
    ...
})
```