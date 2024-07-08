package am_smtp

import (
	"bytes"
	"gopkg.in/gomail.v2"
	"html/template"
)

// EmailConnection结构体
type EmailConnection struct {
	Server   string // smtp服务器地址
	Port     int    // smtp服务器端口
	Username string // smtp登录名
	Password string // smtp密码，可以是授权码
}

/*
* 发送平信（纯文字信）
* 参数：
* 1. from string：发信邮箱地址
* 2. to []string：寄送邮箱地址，可以多个
* 3. subject string：邮件主题
* 4. text string：邮件内容
* 5. conn EmailConnection：邮件连接信息
 */
func SendPlainMail(
	from string,
	to []string,
	subject string,
	text string,
	conn EmailConnection,
) error {
	message := gomail.NewMessage()
	message.SetHeader("From", from)
	message.SetHeader("To", to...)
	message.SetHeader("Subject", subject)
	message.SetBody("text/plain", text)

	d := gomail.NewDialer(conn.Server, conn.Port, conn.Username, conn.Password)
	if err := d.DialAndSend(message); err != nil {
		return err
	}

	return nil
}

/*
* 发送HTML信
* 参数：
* 1. from string：发信邮箱地址
* 2. to []string：寄送邮箱地址，可以多个
* 3. subject string：邮件主题
* 4. html string：邮件内容html模板所在的路径，必须能找到
* 5. data interface{}：需要加载到html模板中的数据集合
* 6. conn EmailConnection：邮件连接信息
 */
func SendHTMLMail(
	from string,
	to []string,
	subject string,
	html string,
	data interface{},
	conn EmailConnection,
) error {
	tmpl, err := template.ParseFiles(html)
	if err != nil {
		return err
	}

	var body bytes.Buffer
	if err := tmpl.Execute(&body, data); err != nil {
		return err
	}

	message := gomail.NewMessage()
	message.SetHeader("From", from)
	message.SetHeader("To", to...)
	message.SetHeader("Subject", subject)
	message.SetBody("text/html", body.String())

	d := gomail.NewDialer(conn.Server, conn.Port, conn.Username, conn.Password)
	if err := d.DialAndSend(message); err != nil {
		return err
	}
	return nil
}
