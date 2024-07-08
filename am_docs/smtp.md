# SMTP 发信服务包
用于通过Email的SMTP服务发送电子邮件。

## 使用方法
### 配置服务连接
```go
var emailConnection = am_smtp.EmailConnection{
	Server: "smtp.example.com",
	Port: 587,
	Username: "example@example.com",
	Password: "123456"
}
```
### 发送纯文字邮件
```go
am_smtp.SendPlainEmail(
    "example@example.com", // from，发信邮箱
    []string{"to@to.com"}, // to，收信邮箱，可以多个
    "subject", // subject，邮件主题
    "content", // text，邮件内容
    emailConnection // conn，上一步设定的email connection信息
)
```

### 发送HTML邮件
```go
// 设定html路径
pwd, err := os.Getwd()
path := filepath.Join(pwd, "template", "signup", "template.html")

// 设定填充数据（有的话）
type EmailData struct {
	Code string
	Time string
}

var emailData = EmailData{
	...
}

// 发送
am_smtp.SendHTMLEmail(
    "example@example.com", // from，发信邮箱
    []string{"to@to.com"}, // to，收信邮箱，可以多个
    "subject", // subject，邮件主题
    path, // html，HTML模板的路径
    emailData, // 填充数据
    emailConnection, // conn，上一步设定的email connection信息
)
```