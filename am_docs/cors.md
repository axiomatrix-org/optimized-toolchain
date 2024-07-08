# Cors
解决跨域问题的Cors中间件，适用于gin框架。
## 使用方法

```go
gin.SetMode(gin.DebugMode)
r := gin.Default()
r.Use(cors.Cors())
```
