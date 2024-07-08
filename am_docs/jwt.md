# JWT 工具套件
用于JWT（Json Web Token）的生成和验证。工具包附带一个全功能中间件。

## 使用方法
### Token的生成

```go
// 配置claims
var tokenClaims = am_jwt.TokenClaims{
	Email: "example@example.com", // 用户的邮件地址
	Role: "user", // Role有四个取值级别：root、admin、user和temp，其中temp仅用于注册和重设密码
	Exp: 600, // token过期时间，以s为计数单位
	Issuer: "James", // 签发人
	SECRET: "ROMANCETILLDEATH", // 签发密钥
}

var token = am_jwt.GenToken(&tokenClaims) // 生成token字串
```

### Token验证中间件的使用
按照普通中间件的使用方式加到需要的controller上即可。
```go
gin.SetMode(gin.DebugMode)
r := gin.Default()
// 参数是该controller访问所需要的权限。比该权限高的token均可访问，但temp除外，需要权限为temp的controller，user权限无法访问。
r.POST("/test", am_jwt.JWTAuthMiddleware("user"), func(context *gin.Context) {
	...
})
```